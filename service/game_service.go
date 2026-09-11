package service

import (
	"errors"
	"game-store/entity"
	"game-store/repository"
	"regexp"
)

var licenseKeyPattern = regexp.MustCompile(`^[A-Z0-9]{4}-[A-Z0-9]{4}-[A-Z0-9]{4}$`)

// GameWithStock holds the game details along with its available stock
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

// GetAllGames take all games from the database and return them as a slice of GameWithStock
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

// GetGameByID take the game ID and return the game details along with its category name and available stock
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
	if !licenseKeyPattern.MatchString(licenseKey) {
		return errors.New("format license key tidak valid, gunakan format XXXX-YYYY-ZZZZ")
	}
	return s.gameRepo.AddGameKey(gameID, licenseKey)
}
