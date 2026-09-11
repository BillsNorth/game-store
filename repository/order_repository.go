package repository

import (
	"database/sql"
	"fmt"
	"game-store/entity"
)

type OrderRepository interface {
	CreateOrderWithTransaction(userID int, cartItems []entity.Cart, grandTotal float64) (int, error)
	GetOrdersByUserID(userID int) ([]entity.Order, error)
	GetOrderDetailsByUserID(userID int) ([]entity.OrderDetail, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderWithTransaction(userID int, cartItems []entity.Cart, grandTotal float64) (int, error) {
	// Database Transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	// Checking balance
	var currentBalance float64
	err = tx.QueryRow(`SELECT wallet_balance FROM user_profiles WHERE user_id = ?`, userID).Scan(&currentBalance)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil data saldo: %v", err)
	}

	if currentBalance < grandTotal {
		return 0, fmt.Errorf("saldo tidak cukup. Saldo kamu: Rp%.2f, Total belanja: Rp%.2f", currentBalance, grandTotal)
	}

	// Current balance is sufficient, proceed with the transaction
	_, err = tx.Exec(`UPDATE user_profiles SET wallet_balance = wallet_balance - ? WHERE user_id = ?`, grandTotal, userID)
	if err != nil {
		return 0, fmt.Errorf("gagal memotong saldo: %v", err)
	}

	// Record the order in the orders table
	res, err := tx.Exec(`INSERT INTO orders (user_id, total_price) VALUES (?, ?)`, userID, grandTotal)
	if err != nil {
		return 0, fmt.Errorf("gagal membuat order: %v", err)
	}

	orderID64, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	orderID := int(orderID64)

	// Take each item in the cart and assign a game key to it
	for _, item := range cartItems {
		for i := 0; i < item.Quantity; i++ {
			var keyID int
			var licenseKey string

			keyQuery := `SELECT id, license_key FROM game_keys WHERE game_id = ? AND status = 'available' LIMIT 1 FOR UPDATE`
			err := tx.QueryRow(keyQuery, item.GameID).Scan(&keyID, &licenseKey)
			if err == sql.ErrNoRows {
				return 0, fmt.Errorf("stok game '%s' mendadak habis", item.GameTitle)
			} else if err != nil {
				return 0, err
			}

			_, err = tx.Exec(`UPDATE game_keys SET status = 'sold' WHERE id = ?`, keyID)
			if err != nil {
				return 0, err
			}

			// Catat ke order_details (Ubah price jadi price_at_purchase)
			_, err = tx.Exec(
				`INSERT INTO order_details (order_id, game_key_id, price_at_purchase) VALUES (?, ?, ?)`,
				orderID, keyID, item.Price,
			)
			if err != nil {
				return 0, err
			}
		}
	}

	// Empty the user's cart after successfully creating the order
	_, err = tx.Exec(`DELETE FROM carts WHERE user_id = ?`, userID)
	if err != nil {
		return 0, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return orderID, nil
}

func (r *orderRepository) GetOrdersByUserID(userID int) ([]entity.Order, error) {
	query := `SELECT id, user_id, total_price, created_at FROM orders WHERE user_id = ? ORDER BY id DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var o entity.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.TotalPrice, &o.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, rows.Err()
}

func (r *orderRepository) GetOrderDetailsByUserID(userID int) ([]entity.OrderDetail, error) {
	query := `
		SELECT od.id, od.order_id, g.id, g.title, od.game_key_id, gk.license_key, od.price_at_purchase
		FROM order_details od
		JOIN orders o ON od.order_id = o.id
		JOIN game_keys gk ON od.game_key_id = gk.id
		JOIN games g ON gk.game_id = g.id
		WHERE o.user_id = ?
		ORDER BY od.order_id DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []entity.OrderDetail
	for rows.Next() {
		var d entity.OrderDetail
		if err := rows.Scan(&d.ID, &d.OrderID, &d.GameID, &d.GameTitle, &d.GameKeyID, &d.LicenseKey, &d.Price); err != nil {
			return nil, err
		}
		details = append(details, d)
	}
	return details, rows.Err()
}
