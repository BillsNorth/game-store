package repository

import (
	"database/sql"
	"game-store/entity"
)

type UserRepository interface {
	Register(user entity.User, profile entity.UserProfile) error
	FindByEmail(email string) (entity.User, error)
	TopUpBalance(userID int, amount float64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Register(user entity.User, profile entity.UserProfile) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	queryUser := "INSERT INTO users (email, password, role) VALUES (?, ?, ?)"
	res, err := tx.Exec(queryUser, user.Email, user.Password, user.Role)
	if err != nil {
		tx.Rollback()
		return err
	}

	userID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	queryProfile := "INSERT INTO user_profiles (user_id, full_name, wallet_balance) VALUES (?, ?, ?)"
	_, err = tx.Exec(queryProfile, userID, profile.FullName, profile.WalletBalance)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *userRepository) FindByEmail(email string) (entity.User, error) {
	var user entity.User
	query := "SELECT id, email, password, role, created_at FROM users WHERE email = ?"
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) TopUpBalance(userID int, amount float64) error {
	query := `UPDATE user_profiles SET wallet_balance = wallet_balance + ? WHERE user_id = ?`
	_, err := r.db.Exec(query, amount, userID)
	return err
}
