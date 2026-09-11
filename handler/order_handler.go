package handler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"game-store/service"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

// Checkout
func (h *OrderHandler) Checkout(userID int) {
	fmt.Println("\n========== CHECKOUT ==========")

	orderID, err := h.orderService.Checkout(userID)

	if err != nil {
		fmt.Println("Checkout gagal:", err)
		return
	}

	fmt.Println("Checkout berhasil!")
	fmt.Println("Order ID:", orderID)
}

// Menampilkan order milik user
func (h *OrderHandler) ShowOrders(userID int) {
	fmt.Println("\n========== RIWAYAT ORDER ==========")

	orders, err := h.orderService.GetUserOrders(userID)

	if err != nil {
		fmt.Println("Gagal mengambil order:", err)
		return
	}

	if len(orders) == 0 {
		fmt.Println("Belum ada order.")
		return
	}

	for _, order := range orders {
		fmt.Printf(
			"Order ID: %d | Total: Rp%.0f | Created: %s\n",
			order.ID,
			order.TotalPrice,
			order.CreatedAt,
		)
	}
}

// Detail order
func (h *OrderHandler) ShowOrderDetail() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n========== DETAIL ORDER ==========")

	fmt.Print("Order ID: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	orderID, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Order ID harus berupa angka.")
		return
	}

	details, err := h.orderService.GetOrderDetails(orderID)

	if err != nil {
		fmt.Println("Gagal mengambil detail order:", err)
		return
	}

	if len(details) == 0 {
		fmt.Println("Detail order tidak ditemukan.")
		return
	}

	for _, detail := range details {
		fmt.Printf(
			"Game: %s | Key: %s | Harga: Rp%.0f\n",
			detail.GameTitle,
			detail.LicenseKey,
			detail.Price,
		)
	}
}
