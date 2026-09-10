package repository

import (
	"database/sql"
	"game-store/entity"
)

type CartRepository interface {
	AddToCart(userID, gameID, quantity int) error
	GetCartByUserID(userID int) ([]entity.Cart, error)
	DeleteCartItem(cartID, userID int) error
	ClearCartByUserID(userID int) error
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

// AddToCart adds an item to the cart or updates the quantity if the game is already in the cart.
func (r *cartRepository) AddToCart(userID, gameID, quantity int) error {
	// Cek apakah item sudah ada di keranjang
	var existingID, existingQty int
	checkQuery := `SELECT id, quantity FROM carts WHERE user_id = ? AND game_id = ?`
	err := r.db.QueryRow(checkQuery, userID, gameID).Scan(&existingID, &existingQty)

	if err == sql.ErrNoRows {
		insertQuery := `INSERT INTO carts (user_id, game_id, quantity) VALUES (?, ?, ?)`
		_, err := r.db.Exec(insertQuery, userID, gameID, quantity)
		return err
	} else if err != nil {
		return err
	}

	updateQuery := `UPDATE carts SET quantity = quantity + ? WHERE id = ?`
	_, err = r.db.Exec(updateQuery, quantity, existingID)
	return err
}

// GetCartByUserID retrieves the list of carts belonging to a specific user, along with the game details.
func (r *cartRepository) GetCartByUserID(userID int) ([]entity.Cart, error) {
	query := `
		SELECT c.id, c.user_id, c.game_id, g.title, g.price, c.quantity, c.created_at
		FROM carts c
		JOIN games g ON c.game_id = g.id
		WHERE c.user_id = ?
		ORDER BY c.id ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []entity.Cart
	for rows.Next() {
		var cart entity.Cart
		err := rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.GameID,
			&cart.GameTitle,
			&cart.Price,
			&cart.Quantity,
			&cart.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return carts, nil
}

// DeleteCartItem removes one item from the cart based on cart_id and user_id.
func (r *cartRepository) DeleteCartItem(cartID, userID int) error {
	query := `DELETE FROM carts WHERE id = ? AND user_id = ?`
	_, err := r.db.Exec(query, cartID, userID)
	return err
}

// ClearCartByUserID removes all items from the user's cart after checkout.
func (r *cartRepository) ClearCartByUserID(userID int) error {
	query := `DELETE FROM carts WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}
