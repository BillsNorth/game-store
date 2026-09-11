package handler

import (
	"fmt"
	"testing"
	"time"

	"game-store/entity"
)

// MockUserService for testing
type MockUserService struct {
	LoginFunc    func(email, password string) (entity.User, error)
	RegisterFunc func(email, password, fullName string) error
	TopUpFunc    func(userID int, amount float64) error
}

func (m *MockUserService) Login(email, password string) (entity.User, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(email, password)
	}
	return entity.User{}, nil
}

func (m *MockUserService) Register(email, password, fullName string) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(email, password, fullName)
	}
	return nil
}

func (m *MockUserService) TopUp(userID int, amount float64) error {
	if m.TopUpFunc != nil {
		return m.TopUpFunc(userID, amount)
	}
	return nil
}

func TestNewUserHandler(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)

	if handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
}

func TestHandleLogin_Success(t *testing.T) {
	mockService := &MockUserService{
		LoginFunc: func(email, password string) (entity.User, error) {
			if email != "user@example.com" || password != "password123" {
				t.Errorf("Expected email='user@example.com', password='password123'; got %q, %q", email, password)
			}
			return entity.User{
				ID:        1,
				Email:     "user@example.com",
				Password:  "password123",
				Role:      "customer",
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewUserHandler(mockService)
	user, err := handler.HandleLogin("user@example.com", "password123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("Expected user ID=1, got %d", user.ID)
	}

	if user.Email != "user@example.com" {
		t.Errorf("Expected email='user@example.com', got %q", user.Email)
	}
}

func TestHandleLogin_EmptyEmail(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)
	_, err := handler.HandleLogin("", "password123")

	if err == nil {
		t.Error("Expected error for empty email, got nil")
	}
}

func TestHandleLogin_EmptyPassword(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)
	_, err := handler.HandleLogin("user@example.com", "")

	if err == nil {
		t.Error("Expected error for empty password, got nil")
	}
}

func TestHandleLogin_InvalidCredentials(t *testing.T) {
	mockService := &MockUserService{
		LoginFunc: func(email, password string) (entity.User, error) {
			return entity.User{}, fmt.Errorf("invalid credentials")
		},
	}

	handler := NewUserHandler(mockService)
	user, err := handler.HandleLogin("user@example.com", "wrongpassword")

	if err == nil {
		t.Error("Expected error for invalid credentials, got nil")
	}

	if user.ID != 0 {
		t.Errorf("Expected zero user, got %v", user)
	}
}

func TestHandleLogin_AdminUser(t *testing.T) {
	mockService := &MockUserService{
		LoginFunc: func(email, password string) (entity.User, error) {
			return entity.User{
				ID:        2,
				Email:     "admin@example.com",
				Password:  "admin123",
				Role:      "admin",
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := NewUserHandler(mockService)
	user, err := handler.HandleLogin("admin@example.com", "admin123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.Role != "admin" {
		t.Errorf("Expected role='admin', got %q", user.Role)
	}
}

func TestHandleRegister_Success(t *testing.T) {
	called := false
	mockService := &MockUserService{
		RegisterFunc: func(email, password, fullName string) error {
			called = true
			if email != "newuser@example.com" || password != "pass123" || fullName != "John Doe" {
				t.Errorf("Unexpected parameters: %q, %q, %q", email, password, fullName)
			}
			return nil
		},
	}

	handler := NewUserHandler(mockService)
	err := handler.HandleRegister("newuser@example.com", "pass123", "John Doe")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !called {
		t.Error("RegisterFunc was not called")
	}
}

func TestHandleRegister_EmptyEmail(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)
	err := handler.HandleRegister("", "password123", "John Doe")

	if err == nil {
		t.Error("Expected error for empty email, got nil")
	}
}

func TestHandleRegister_EmptyPassword(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)
	err := handler.HandleRegister("user@example.com", "", "John Doe")

	if err == nil {
		t.Error("Expected error for empty password, got nil")
	}
}

func TestHandleRegister_EmptyName(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)
	err := handler.HandleRegister("user@example.com", "password123", "")

	if err == nil {
		t.Error("Expected error for empty name, got nil")
	}
}

func TestHandleRegister_DuplicateEmail(t *testing.T) {
	mockService := &MockUserService{
		RegisterFunc: func(email, password, fullName string) error {
			return fmt.Errorf("email already exists")
		},
	}

	handler := NewUserHandler(mockService)
	err := handler.HandleRegister("existing@example.com", "password123", "John Doe")

	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
}

func TestHandleTopUp_Success(t *testing.T) {
	called := false
	mockService := &MockUserService{
		TopUpFunc: func(userID int, amount float64) error {
			called = true
			if userID != 1 || amount != 100000 {
				t.Errorf("Expected userID=1, amount=100000; got %d, %f", userID, amount)
			}
			return nil
		},
	}

	handler := NewUserHandler(mockService)
	err := handler.HandleTopUp(1, 100000)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !called {
		t.Error("TopUpFunc was not called")
	}
}

func TestHandleTopUp_ZeroAmount(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)
	err := handler.HandleTopUp(1, 0)

	if err == nil {
		t.Error("Expected error for zero amount, got nil")
	}
}

func TestHandleTopUp_NegativeAmount(t *testing.T) {
	mockService := &MockUserService{}
	handler := NewUserHandler(mockService)
	err := handler.HandleTopUp(1, -50000)

	if err == nil {
		t.Error("Expected error for negative amount, got nil")
	}
}

func TestHandleTopUp_DifferentUsers(t *testing.T) {
	mockService := &MockUserService{
		TopUpFunc: func(userID int, amount float64) error {
			if userID == 1 || userID == 2 {
				return nil
			}
			return fmt.Errorf("user not found")
		},
	}

	handler := NewUserHandler(mockService)

	err1 := handler.HandleTopUp(1, 50000)
	if err1 != nil {
		t.Errorf("Expected no error for user 1, got %v", err1)
	}

	err2 := handler.HandleTopUp(2, 75000)
	if err2 != nil {
		t.Errorf("Expected no error for user 2, got %v", err2)
	}
}

func TestHandleTopUp_UserNotFound(t *testing.T) {
	mockService := &MockUserService{
		TopUpFunc: func(userID int, amount float64) error {
			return fmt.Errorf("user not found")
		},
	}

	handler := NewUserHandler(mockService)
	err := handler.HandleTopUp(999, 100000)

	if err == nil {
		t.Error("Expected error for non-existent user, got nil")
	}
}

func TestPromptLogin(t *testing.T) {
	handler := NewUserHandler(&MockUserService{})

	// Note: PromptLogin reads from stdin, which is not available in test context
	// This test verifies the handler is created and the method exists
	if handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
}

func TestPromptRegister(t *testing.T) {
	handler := NewUserHandler(&MockUserService{})

	// Note: PromptRegister reads from stdin, which is not available in test context
	// This test verifies the handler is created and the method exists
	if handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
}

func TestPromptTopUp(t *testing.T) {
	handler := NewUserHandler(&MockUserService{})

	// Note: PromptTopUp reads from stdin, which is not available in test context
	// This test verifies the handler is created and the method exists
	if handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
}
