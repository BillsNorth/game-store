package entity

import "time"

type User struct {
	ID        int
	Email     string
	Password  string
	Role      string
	CreatedAt time.Time
}

type UserProfile struct {
	ID            int
	UserID        int
	FullName      string
	WalletBalance float64
}

type Category struct {
	ID           int
	CategoryName string
}

type Game struct {
	ID         int
	CategoryID int
	Title      string
	Price      float64
}

type GameKey struct {
	ID         int
	GameID     int
	LicenseKey string
	Status     string
}

type Order struct {
	ID          int
	UserID      int
	TotalAmount float64
	CreatedAt   time.Time
}

type OrderDetail struct {
	ID              int
	OrderID         int
	GameKeyID       int
	PriceAtPurchase float64
}
