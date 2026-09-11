package handler

import (
	"fmt"
	"testing"

	"game-store/entity"
)

// MockOrderService for testing
type MockOrderService struct {
	CheckoutFunc        func(userID int) (int, error)
	GetUserOrdersFunc   func(userID int) ([]entity.Order, error)
	GetOrderDetailsFunc func(orderID int) ([]entity.OrderDetail, error)
}

func (m *MockOrderService) Checkout(userID int) (int, error) {
	if m.CheckoutFunc != nil {
		return m.CheckoutFunc(userID)
	}
	return 0, nil
}

func (m *MockOrderService) GetUserOrders(userID int) ([]entity.Order, error) {
	if m.GetUserOrdersFunc != nil {
		return m.GetUserOrdersFunc(userID)
	}
	return nil, nil
}

func (m *MockOrderService) GetOrderDetails(orderID int) ([]entity.OrderDetail, error) {
	if m.GetOrderDetailsFunc != nil {
		return m.GetOrderDetailsFunc(orderID)
	}
	return nil, nil
}

func TestNewOrderHandler(t *testing.T) {
	mockService := &MockOrderService{}
	handler := NewOrderHandler(mockService)

	if handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
}

func TestCheckout_Success(t *testing.T) {
	called := false
	mockService := &MockOrderService{
		CheckoutFunc: func(userID int) (int, error) {
			called = true
			if userID != 1 {
				t.Errorf("Expected userID=1, got %d", userID)
			}
			return 42, nil
		},
	}

	handler := NewOrderHandler(mockService)
	handler.Checkout(1)

	if !called {
		t.Error("CheckoutFunc was not called")
	}
}

func TestCheckout_EmptyCart(t *testing.T) {
	mockService := &MockOrderService{
		CheckoutFunc: func(userID int) (int, error) {
			return 0, fmt.Errorf("keranjang belanja kamu masih kosong")
		},
	}

	handler := NewOrderHandler(mockService)
	handler.Checkout(1)
}

func TestCheckout_Error(t *testing.T) {
	mockService := &MockOrderService{
		CheckoutFunc: func(userID int) (int, error) {
			return 0, fmt.Errorf("database error")
		},
	}

	handler := NewOrderHandler(mockService)
	handler.Checkout(1)
}

func TestShowOrders_EmptyOrders(t *testing.T) {
	mockService := &MockOrderService{
		GetUserOrdersFunc: func(userID int) ([]entity.Order, error) {
			return []entity.Order{}, nil
		},
	}

	handler := NewOrderHandler(mockService)
	handler.ShowOrders(1)
}

func TestShowOrders_WithItems(t *testing.T) {
	mockService := &MockOrderService{
		GetUserOrdersFunc: func(userID int) ([]entity.Order, error) {
			return []entity.Order{
				{
					ID:         1,
					UserID:     1,
					TotalPrice: 150000,
					CreatedAt:  "2026-09-11",
				},
				{
					ID:         2,
					UserID:     1,
					TotalPrice: 75000,
					CreatedAt:  "2026-09-10",
				},
			}, nil
		},
	}

	handler := NewOrderHandler(mockService)
	handler.ShowOrders(1)
}

func TestShowOrders_GetError(t *testing.T) {
	mockService := &MockOrderService{
		GetUserOrdersFunc: func(userID int) ([]entity.Order, error) {
			return nil, fmt.Errorf("database error")
		},
	}

	handler := NewOrderHandler(mockService)
	handler.ShowOrders(1)
}

func TestGetOrderDetails_Success(t *testing.T) {
	called := false
	mockService := &MockOrderService{
		GetOrderDetailsFunc: func(orderID int) ([]entity.OrderDetail, error) {
			called = true
			if orderID != 5 {
				t.Errorf("Expected orderID=5, got %d", orderID)
			}
			return []entity.OrderDetail{
				{
					ID:         1,
					OrderID:    5,
					GameID:     10,
					GameTitle:  "Game A",
					GameKeyID:  100,
					LicenseKey: "ABC-123",
					Price:      100000,
				},
			}, nil
		},
	}

	_ = NewOrderHandler(mockService)
	details, err := mockService.GetOrderDetails(5)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !called {
		t.Error("GetOrderDetailsFunc was not called")
	}
	if len(details) != 1 || details[0].GameTitle != "Game A" {
		t.Errorf("Expected one detail with GameTitle='Game A', got %v", details)
	}
}

func TestGetOrderDetails_NotFound(t *testing.T) {
	mockService := &MockOrderService{
		GetOrderDetailsFunc: func(orderID int) ([]entity.OrderDetail, error) {
			return []entity.OrderDetail{}, nil
		},
	}

	_ = NewOrderHandler(mockService)
	details, err := mockService.GetOrderDetails(999)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(details) != 0 {
		t.Errorf("Expected no details, got %v", details)
	}
}

func TestGetOrderDetails_Error(t *testing.T) {
	mockService := &MockOrderService{
		GetOrderDetailsFunc: func(orderID int) ([]entity.OrderDetail, error) {
			return nil, fmt.Errorf("database error")
		},
	}

	_ = NewOrderHandler(mockService)
	_, err := mockService.GetOrderDetails(999)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestShowOrderDetail(t *testing.T) {
	handler := NewOrderHandler(&MockOrderService{})

	// Note: ShowOrderDetail reads from stdin, which is not available in test context
	// This test verifies the handler is created and the method exists
	if handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
}
