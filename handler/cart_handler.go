package handler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"game-store/service"
)

type CartHandler struct {
	cartService service.CartService
}

func NewCartHandler(cartService service.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

// Menampilkan cart
func (h *CartHandler) ShowCart(userID int) {
	fmt.Println("\n========== CART ==========")

	cart, err := h.cartService.GetCart(userID)
	if err != nil {
		fmt.Println("Gagal mengambil cart:", err)
		return
	}

	if len(cart.Items) == 0 {
		fmt.Println("Keranjang masih kosong.")
		return
	}

	for _, item := range cart.Items {
		fmt.Printf(
			"Cart ID: %d | Game ID: %d | %s | Qty: %d | Rp%.0f\n",
			item.ID,
			item.GameID,
			item.GameTitle,
			item.Quantity,
			item.Price,
		)
	}

	fmt.Println("----------------------------")
	fmt.Printf("Grand Total: Rp%.0f\n", cart.GrandTotal)
}

// Tambah game ke cart
func (h *CartHandler) AddToCart(userID int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n========== TAMBAH KE CART ==========")

	fmt.Print("Game ID: ")
	gameInput, _ := reader.ReadString('\n')
	gameInput = strings.TrimSpace(gameInput)

	gameID, err := strconv.Atoi(gameInput)
	if err != nil {
		fmt.Println("Game ID harus berupa angka.")
		return
	}

	fmt.Print("Quantity: ")
	qtyInput, _ := reader.ReadString('\n')
	qtyInput = strings.TrimSpace(qtyInput)

	quantity, err := strconv.Atoi(qtyInput)
	if err != nil {
		fmt.Println("Quantity harus berupa angka.")
		return
	}

	err = h.cartService.AddToCart(
		userID,
		gameID,
		quantity,
	)

	if err != nil {
		fmt.Println("Gagal:", err)
		return
	}

	fmt.Println("Game berhasil ditambahkan ke cart.")
}

// Hapus item dari cart
func (h *CartHandler) RemoveFromCart(userID int) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n========== HAPUS CART ==========")

	fmt.Print("Cart ID yang ingin dihapus: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	cartID, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Cart ID harus berupa angka.")
		return
	}

	err = h.cartService.RemoveFromCart(
		cartID,
		userID,
	)

	if err != nil {
		fmt.Println("Gagal menghapus item:", err)
		return
	}

	fmt.Println("Item berhasil dihapus dari cart.")
}
