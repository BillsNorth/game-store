package main

import (
	"log"

	"game-store/cli"
	"game-store/config"
	"game-store/handler"
	"game-store/repository"
	"game-store/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Intialisasi Database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer db.Close()

	// Intialisasi Repositories
	userRepo := repository.NewUserRepository(db)
	gameRepo := repository.NewGameRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Intialisasi Services
	userService := service.NewUserService(userRepo)
	gameService := service.NewGameService(gameRepo)
	cartService := service.NewCartService(cartRepo, gameRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo)

	// Intialisasi Handlers
	gameHandler := handler.NewGameHandler(gameService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	userHandler := handler.NewUserHandler(userService)

	// Intialisasi CLI dan jalankan
	gameStoreCLI := cli.NewCLI(gameHandler, cartHandler, orderHandler, userHandler, gameService)
	gameStoreCLI.Start()
}
