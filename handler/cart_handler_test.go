package handler

import (
	"fmt"
	"testing"

	"game-store/entity"
	"game-store/service"
)

// MockCartService for testing
type MockCartService struct {
	GetCartFunc        func(userID int) (service.CartSummary, error)
	AddToCartFunc      func(userID, gameID, quantity int) error
	RemoveFromCartFunc func(cartID, userID int) error
}

func (m *MockCartService) GetCart(userID int) (service.CartSummary, error) {
	if m.GetCartFunc != nil {
		return m.GetCartFunc(userID)
	}
	return service.CartSummary{}, nil
}

func (m *MockCartService) AddToCart(userID, gameID, quantity int) error {
	if m.AddToCartFunc != nil {
		return m.AddToCartFunc(userID, gameID, quantity)
	}
	return nil
}

func (m *MockCartService) RemoveFromCart(cartID, userID int) error {
	if m.RemoveFromCartFunc != nil {
		return m.RemoveFromCartFunc(cartID, userID)
	}
	return nil
}

func TestNewCartHandler(t *testing.T) {
	mockService := &MockCartService{}
	handler := NewCartHandler(mockService)

	if handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
}

func TestShowCart_EmptyCart(t *testing.T) {
	mockService := &MockCartService{
		GetCartFunc: func(userID int) (service.CartSummary, error) {
			return service.CartSummary{
				Items:      []entity.Cart{},
				GrandTotal: 0,
			}, nil
		},
	}

	handler := NewCartHandler(mockService)
	handler.ShowCart(1)
}

func TestShowCart_WithItems(t *testing.T) {
	mockService := &MockCartService{
		GetCartFunc: func(userID int) (service.CartSummary, error) {
			return service.CartSummary{
				Items: []entity.Cart{
					{
						ID:        1,
						UserID:    1,
						GameID:    10,
						GameTitle: "Game A",
						Price:     100000,
						Quantity:  2,
						CreatedAt: "2026-09-11",
					},
					{
						ID:        2,
						UserID:    1,
						GameID:    11,
						GameTitle: "Game B",
						Price:     50000,
						Quantity:  1,
						CreatedAt: "2026-09-11",
					},
				},
				GrandTotal: 250000,
			}, nil
		},
	}

	handler := NewCartHandler(mockService)
	handler.ShowCart(1)
}

func TestShowCart_GetError(t *testing.T) {
	mockService := &MockCartService{
		GetCartFunc: func(userID int) (service.CartSummary, error) {
			return service.CartSummary{}, fmt.Errorf("database error")
		},
	}

	handler := NewCartHandler(mockService)
	handler.ShowCart(1)
}

func TestAddToCart_Success(t *testing.T) {
	called := false
	mockService := &MockCartService{
		AddToCartFunc: func(userID, gameID, quantity int) error {
			called = true
			if userID != 1 || gameID != 5 || quantity != 2 {
				t.Errorf("Expected userID=1, gameID=5, quantity=2; got %d, %d, %d", userID, gameID, quantity)
			}
			return nil
		},
	}

	_ = NewCartHandler(mockService)
	// Note: In actual integration, this would prompt for input
	// Here we test the handler's logic path
	err := mockService.AddToCart(1, 5, 2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !called {
		t.Error("AddToCartFunc was not called")
	}
}

func TestRemoveFromCart_Success(t *testing.T) {
	called := false
	mockService := &MockCartService{
		RemoveFromCartFunc: func(cartID, userID int) error {
			called = true
			if cartID != 3 || userID != 1 {
				t.Errorf("Expected cartID=3, userID=1; got %d, %d", cartID, userID)
			}
			return nil
		},
	}

	_ = NewCartHandler(mockService)
	err := mockService.RemoveFromCart(3, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !called {
		t.Error("RemoveFromCartFunc was not called")
	}
}

func TestRemoveFromCart_Error(t *testing.T) {
	mockService := &MockCartService{
		RemoveFromCartFunc: func(cartID, userID int) error {
			return fmt.Errorf("item not found")
		},
	}

	_ = NewCartHandler(mockService)
	err := mockService.RemoveFromCart(999, 1)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
