package handler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"game-store/service"
)

type GameHandler struct {
	gameService service.GameService
}

func NewGameHandler(gameService service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// Browse semua game
func (h *GameHandler) BrowseGames() {
	fmt.Println("\n========== DAFTAR GAME ==========")

	games, err := h.gameService.GetAllGames()
	if err != nil {
		fmt.Println("Gagal mengambil data game:", err)
		return
	}

	if len(games) == 0 {
		fmt.Println("Belum ada game tersedia.")
		return
	}

	for _, item := range games {
		fmt.Printf(
			"ID: %d | %s | Rp%.0f | Stok: %d\n",
			item.Game.ID,
			item.Game.Title,
			item.Game.Price,
			item.Stock,
		)
	}
}

// Detail game berdasarkan ID
func (h *GameHandler) ShowGameDetail() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nMasukkan ID game: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	id, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("ID game harus berupa angka.")
		return
	}

	game, err := h.gameService.GetGameByID(id)
	if err != nil {
		fmt.Println("Gagal:", err)
		return
	}

	fmt.Println("\n========== DETAIL GAME ==========")
	fmt.Println("ID          :", game.Game.ID)
	fmt.Println("Nama        :", game.Game.Title)
	fmt.Println("Harga       :", game.Game.Price)
	fmt.Println("Category   :", game.Game.CategoryName)
	fmt.Println("Stok        :", game.Stock)
}
