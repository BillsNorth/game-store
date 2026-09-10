package repository

import (
	"database/sql"
	"game-store/entity"
)

type GameRepository interface {
	GetAll() ([]entity.Game, error)
	GetByID(id int) (entity.Game, error)
	GetAvailableKeyCount(gameID int) (int, error)
	AddGameKey(gameID int, licenseKey string) error
}

type gameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{db: db}
}

// GetAll take all games from the database and return them as a slice of entity.Game
func (r *gameRepository) GetAll() ([]entity.Game, error) {
	query := `
		SELECT g.id, g.category_id, c.category_name, g.title, g.price
		FROM games g
		JOIN categories c ON g.category_id = c.id
		ORDER BY g.id ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []entity.Game
	for rows.Next() {
		var g entity.Game
		err := rows.Scan(&g.ID, &g.CategoryID, &g.CategoryName, &g.Title, &g.Price)
		if err != nil {
			return nil, err
		}
		games = append(games, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}

// GetByID take the game ID and return the game details along with its category name
func (r *gameRepository) GetByID(id int) (entity.Game, error) {
	query := `
		SELECT g.id, g.category_id, c.category_name, g.title, g.price
		FROM games g
		JOIN categories c ON g.category_id = c.id
		WHERE g.id = ?
	`

	var g entity.Game
	err := r.db.QueryRow(query, id).Scan(&g.ID, &g.CategoryID, &g.CategoryName, &g.Title, &g.Price)
	if err != nil {
		return g, err
	}

	return g, nil
}

// GetAvailableKeyCount take the game ID and return the count of available keys for that game
func (r *gameRepository) GetAvailableKeyCount(gameID int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM game_keys 
		WHERE game_id = ? AND status = 'available'
	`

	var count int
	err := r.db.QueryRow(query, gameID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *gameRepository) AddGameKey(gameID int, licenseKey string) error {
	query := `INSERT INTO game_keys (game_id, license_key, status) VALUES (?, ?, 'available')`
	_, err := r.db.Exec(query, gameID, licenseKey)
	return err
}
