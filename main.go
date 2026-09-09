package main

import (
	"log"

	"game-store/config"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer db.Close()
}
