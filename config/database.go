package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDB() (*sql.DB, error) {
	_ = godotenv.Load()

	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DB_URL tidak ditemukan di .env")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka driver: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal koneksi ke database: %w", err)
	}

	fmt.Println("Berhasil terhubung ke database MySQL (game_store)!")
	return db, nil
}
