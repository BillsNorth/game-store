package handler

import (
	"fmt"

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
func (h *OrderHandler) ShowOrderDetail(userID int) {
	fmt.Println("\n========== DETAIL ORDER ==========")

	details, err := h.orderService.GetOrderDetailsByUserID(userID)

	if err != nil {
		fmt.Println("Gagal mengambil detail order:", err)
		return
	}

	if len(details) == 0 {
		fmt.Println("Kamu belum pernah melakukan pemesanan.")
		return
	}

	currentOrderID := -1
	for _, detail := range details {
		if detail.OrderID != currentOrderID {
			currentOrderID = detail.OrderID
			fmt.Printf("\nOrder ID: %d\n", currentOrderID)
		}
		fmt.Printf(
			"  Game: %s | Key: %s | Harga: Rp%.0f\n",
			detail.GameTitle,
			detail.LicenseKey,
			detail.Price,
		)
	}
}

func (h *OrderHandler) ShowAllOrderDetail() {
	fmt.Println("\n========== ALL ORDER ==========")

	details, err := h.orderService.GetAllOrderDetails()

	if err != nil {
		fmt.Println("Gaga mengambil detail order:", err)
		return
	}

	if len(details) == 0 {
		fmt.Println("Belum ada yang membeli lisensi game")
	}

	currentOrderID := -1
	for _, detail := range details {
		if detail.OrderID != currentOrderID {
			currentOrderID = detail.OrderID
			fmt.Printf("\nOrder ID: %v | DATE BUY: %s | Buyer: %s \n", detail.OrderID, detail.DateBuy, detail.Buyer)
		}

		fmt.Printf(
			"             Game: %s | Key: %s | Harga: Rp%.0f\n",
			detail.Title,
			detail.LicenseKey,
			detail.PriceAtPurchase,
		)
	}

}
