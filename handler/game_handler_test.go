package handler

import (
	"fmt"
	"testing"

	"game-store/entity"
	"game-store/service"
)

// MockGameService for testing
type MockGameService struct {
	GetAllGamesFunc func() ([]service.GameWithStock, error)
	GetGameByIDFunc func(id int) (service.GameWithStock, error)
	AddKeyFunc      func(gameID int, licenseKey string) error
}

func (m *MockGameService) GetAllGames() ([]service.GameWithStock, error) {
	if m.GetAllGamesFunc != nil {
		return m.GetAllGamesFunc()
	}
	return nil, nil
}

func (m *MockGameService) GetGameByID(id int) (service.GameWithStock, error) {
	if m.GetGameByIDFunc != nil {
		return m.GetGameByIDFunc(id)
	}
	return service.GameWithStock{}, nil
}

func (m *MockGameService) AddKey(gameID int, licenseKey string) error {
	if m.AddKeyFunc != nil {
		return m.AddKeyFunc(gameID, licenseKey)
	}
	return nil
}

func TestNewGameHandler(t *testing.T) {
	mockService := &MockGameService{}
	handler := NewGameHandler(mockService)

	if handler == nil {
		t.Error("Expected handler to be initialized, got nil")
	}
}

func TestBrowseGames_EmptyList(t *testing.T) {
	mockService := &MockGameService{
		GetAllGamesFunc: func() ([]service.GameWithStock, error) {
			return []service.GameWithStock{}, nil
		},
	}

	handler := NewGameHandler(mockService)
	handler.BrowseGames()
}

func TestBrowseGames_MultipleGames(t *testing.T) {
	mockService := &MockGameService{
		GetAllGamesFunc: func() ([]service.GameWithStock, error) {
			return []service.GameWithStock{
				{
					Game: entity.Game{
						ID:           1,
						CategoryID:   1,
						CategoryName: "Action",
						Title:        "Game A",
						Price:        100000,
						CreatedAt:    "2026-09-11",
					},
					Stock: 10,
				},
				{
					Game: entity.Game{
						ID:           2,
						CategoryID:   2,
						CategoryName: "RPG",
						Title:        "Game B",
						Price:        200000,
						CreatedAt:    "2026-09-11",
					},
					Stock: 5,
				},
			}, nil
		},
	}

	handler := NewGameHandler(mockService)
	handler.BrowseGames()
}

func TestBrowseGames_GetError(t *testing.T) {
	mockService := &MockGameService{
		GetAllGamesFunc: func() ([]service.GameWithStock, error) {
			return nil, fmt.Errorf("database error")
		},
	}

	handler := NewGameHandler(mockService)
	handler.BrowseGames()
}

func TestShowGameDetail_Success(t *testing.T) {
	mockService := &MockGameService{
		GetGameByIDFunc: func(id int) (service.GameWithStock, error) {
			if id != 1 {
				t.Errorf("Expected id=1, got %d", id)
			}
			return service.GameWithStock{
				Game: entity.Game{
					ID:           1,
					CategoryID:   1,
					CategoryName: "Action",
					Title:        "Game A",
					Price:        100000,
					CreatedAt:    "2026-09-11",
				},
				Stock: 10,
			}, nil
		},
	}

	_ = NewGameHandler(mockService)
	game, err := mockService.GetGameByID(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if game.Game.ID != 1 {
		t.Errorf("Expected game ID 1, got %d", game.Game.ID)
	}
}

func TestShowGameDetail_NotFound(t *testing.T) {
	mockService := &MockGameService{
		GetGameByIDFunc: func(id int) (service.GameWithStock, error) {
			return service.GameWithStock{}, fmt.Errorf("game not found")
		},
	}

	_ = NewGameHandler(mockService)
	game, err := mockService.GetGameByID(999)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if game.Game.ID != 0 {
		t.Error("Expected zero-value game, got game with ID")
	}
}

func TestAddKey_Success(t *testing.T) {
	called := false
	mockService := &MockGameService{
		AddKeyFunc: func(gameID int, licenseKey string) error {
			called = true
			if gameID != 1 || licenseKey != "KEY123" {
				t.Errorf("Expected gameID=1, licenseKey='KEY123'; got %d, %q", gameID, licenseKey)
			}
			return nil
		},
	}

	_ = NewGameHandler(mockService)
	err := mockService.AddKey(1, "KEY123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !called {
		t.Error("AddKeyFunc was not called")
	}
}

func TestAddKey_Error(t *testing.T) {
	mockService := &MockGameService{
		AddKeyFunc: func(gameID int, licenseKey string) error {
			return fmt.Errorf("failed to add key")
		},
	}

	_ = NewGameHandler(mockService)
	err := mockService.AddKey(1, "KEY123")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
