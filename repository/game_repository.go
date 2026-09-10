package repository

import (
	"database/sql"
	"game-store/entity"
)

type GameRepository interface {
	GetAll() ([]entity.Game, error)
	GetByID(id int) (entity.Game, error)
	GetAvailableKeyCount(gameID int) (int, error)
}

type gameRepository struct {
	db *sql.DB
}

func NewGameRepository(db *sql.DB) GameRepository {
	return &gameRepository{db: db}
}

// GetAll mengambil semua daftar game beserta nama kategorinya
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

// GetByID mengambil detail 1 game berdasarkan ID
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

// GetAvailableKeyCount menghitung sisa stok lisensi game yang masih 'available'
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
