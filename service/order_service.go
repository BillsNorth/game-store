package service

import (
	"errors"
	"game-store/entity"
	"game-store/repository"
)

type OrderService interface {
	Checkout(userID int) (int, error)
	GetUserOrders(userID int) ([]entity.Order, error)
	GetOrderDetailsByUserID(userID int) ([]entity.OrderDetail, error)
}

type orderService struct {
	orderRepo repository.OrderRepository
	cartRepo  repository.CartRepository
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
	}
}

func (s *orderService) Checkout(userID int) (int, error) {
	// Take the cart items for the user
	cartItems, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return 0, err
	}

	if len(cartItems) == 0 {
		return 0, errors.New("keranjang belanja kamu masih kosong")
	}

	// Calculate the grand total of the cart items
	var grandTotal float64
	for _, item := range cartItems {
		grandTotal += item.Price * float64(item.Quantity)
	}

	// Run the checkout process with transaction
	return s.orderRepo.CreateOrderWithTransaction(userID, cartItems, grandTotal)
}

func (s *orderService) GetUserOrders(userID int) ([]entity.Order, error) {
	return s.orderRepo.GetOrdersByUserID(userID)
}

func (s *orderService) GetOrderDetailsByUserID(userID int) ([]entity.OrderDetail, error) {
	return s.orderRepo.GetOrderDetailsByUserID(userID)
}
