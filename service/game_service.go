package service

import (
	"errors"
	"game-store/entity"
	"game-store/repository"
)

// GameWithStock memuat detail game beserta jumlah stok lisensinya
type GameWithStock struct {
	Game  entity.Game
	Stock int
}

type GameService interface {
	GetAllGames() ([]GameWithStock, error)
	GetGameByID(id int) (GameWithStock, error)
	AddKey(gameID int, licenseKey string) error
}

type gameService struct {
	gameRepo repository.GameRepository
}

func NewGameService(gameRepo repository.GameRepository) GameService {
	return &gameService{gameRepo: gameRepo}
}

// GetAllGames mengambil semua game beserta jumlah stoknya masing-masing
func (s *gameService) GetAllGames() ([]GameWithStock, error) {
	games, err := s.gameRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var result []GameWithStock
	for _, g := range games {
		stock, err := s.gameRepo.GetAvailableKeyCount(g.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, GameWithStock{
			Game:  g,
			Stock: stock,
		})
	}

	return result, nil
}

// GetGameByID mengambil detail 1 game beserta stoknya
func (s *gameService) GetGameByID(id int) (GameWithStock, error) {
	if id <= 0 {
		return GameWithStock{}, errors.New("ID game tidak valid")
	}

	game, err := s.gameRepo.GetByID(id)
	if err != nil {
		return GameWithStock{}, errors.New("game tidak ditemukan")
	}

	stock, err := s.gameRepo.GetAvailableKeyCount(game.ID)
	if err != nil {
		return GameWithStock{}, err
	}

	return GameWithStock{
		Game:  game,
		Stock: stock,
	}, nil
}

func (s *gameService) AddKey(gameID int, licenseKey string) error {
	if licenseKey == "" {
		return errors.New("license key tidak boleh kosong")
	}
	return s.gameRepo.AddGameKey(gameID, licenseKey)
}
