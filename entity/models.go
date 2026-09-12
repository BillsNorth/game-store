package entity

import (
	"time"
)

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
	ID           int
	CategoryID   int
	CategoryName string
	Title        string
	Price        float64
	CreatedAt    string
}

type GameKey struct {
	ID         int
	GameID     int
	LicenseKey string
	Status     string
}

type Cart struct {
	ID        int
	UserID    int
	GameID    int
	GameTitle string
	Price     float64
	Quantity  int
	CreatedAt string
}

type Order struct {
	ID         int
	UserID     int
	TotalPrice float64
	CreatedAt  string
}

type OrderDetail struct {
	ID         int
	OrderID    int
	GameID     int
	GameTitle  string
	GameKeyID  int
	LicenseKey string
	Price      float64
}

type AllOrderDetail struct {
	OrderID         int
	DateBuy         string
	Buyer           string
	Title           string
	GameKeyID       int
	LicenseKey      string
	PriceAtPurchase float64
}
