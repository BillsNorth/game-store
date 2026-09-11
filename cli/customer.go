package cli

import (
	"bufio"
	"fmt"
	"strings"

	"game-store/handler"
)

type CustomerCLI struct {
	gameHandler  *handler.GameHandler
	cartHandler  *handler.CartHandler
	orderHandler *handler.OrderHandler
	userHandler  *handler.UserHandler
	scanner      *bufio.Scanner
	userID       int
	userEmail    string
}

func NewCustomerCLI(
	gameHandler *handler.GameHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	userHandler *handler.UserHandler,
	scanner *bufio.Scanner,
	userID int,
	userEmail string,
) *CustomerCLI {
	return &CustomerCLI{
		gameHandler:  gameHandler,
		cartHandler:  cartHandler,
		orderHandler: orderHandler,
		userHandler:  userHandler,
		scanner:      scanner,
		userID:       userID,
		userEmail:    userEmail,
	}
}

// ShowMenu displays the customer menu
func (c *CustomerCLI) ShowMenu() bool {
	for {
		fmt.Printf("\n=================================\n")
		fmt.Printf(" DASHBOARD CUSTOMER (%s)\n", c.userEmail)
		fmt.Println("=================================")
		fmt.Println("1. Lihat Katalog Game")
		fmt.Println("2. Tambah Game ke Keranjang")
		fmt.Println("3. Lihat Keranjang Belanja")
		fmt.Println("4. Checkout (Beli)")
		fmt.Println("5. Lihat Riwayat Pesanan")
		fmt.Println("6. Lihat Detail Pesanan")
		fmt.Println("7. Top-up Saldo Wallet")
		fmt.Println("8. Logout")
		fmt.Print("Pilih menu: ")

		c.scanner.Scan()
		choice := strings.TrimSpace(c.scanner.Text())

		switch choice {
		case "1":
			c.handleBrowseGames()

		case "2":
			c.handleAddToCart()

		case "3":
			c.handleShowCart()

		case "4":
			c.handleCheckout()

		case "5":
			c.handleShowOrders()

		case "6":
			c.handleShowOrderDetail()

		case "7":
			c.handleTopUp()

		case "8":
			fmt.Println("\n[Sukses] Berhasil logout.")
			return false

		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

// handleBrowseGames displays all games
func (c *CustomerCLI) handleBrowseGames() {
	c.gameHandler.BrowseGames()
}

// handleAddToCart adds game to cart
func (c *CustomerCLI) handleAddToCart() {
	c.cartHandler.AddToCart(c.userID)
}

// handleShowCart displays the user's cart
func (c *CustomerCLI) handleShowCart() {
	c.cartHandler.ShowCart(c.userID)
}

// handleCheckout processes the checkout
func (c *CustomerCLI) handleCheckout() {
	c.orderHandler.Checkout(c.userID)
}

// handleShowOrders displays user's order history
func (c *CustomerCLI) handleShowOrders() {
	c.orderHandler.ShowOrders(c.userID)
}

// handleShowOrderDetail displays details of a specific order
func (c *CustomerCLI) handleShowOrderDetail() {
	c.orderHandler.ShowOrderDetail()
}

// handleTopUp handles wallet top-up
func (c *CustomerCLI) handleTopUp() {
	fmt.Println("\n========== TOP-UP WALLET ==========")
	amount := c.userHandler.PromptTopUp()
	c.userHandler.HandleTopUp(c.userID, amount)
}
