package cli

import (
	"bufio"
	"fmt"
	"strings"

	"game-store/handler"
	"game-store/service"
)

type AdminCLI struct {
	gameHandler *handler.GameHandler
	gameService service.GameService
	scanner     *bufio.Scanner
	userEmail   string
}

func NewAdminCLI(gameHandler *handler.GameHandler, gameService service.GameService, scanner *bufio.Scanner, userEmail string) *AdminCLI {
	return &AdminCLI{
		gameHandler: gameHandler,
		gameService: gameService,
		scanner:     scanner,
		userEmail:   userEmail,
	}
}

// ShowMenu displays the admin menu
func (a *AdminCLI) ShowMenu() bool {
	for {
		fmt.Printf("\n=================================\n")
		fmt.Printf(" MENU ADMIN (%s)\n", a.userEmail)
		fmt.Println("=================================")
		fmt.Println("1. Lihat Katalog Game")
		fmt.Println("2. Tambah Stok Game Key Baru")
		fmt.Println("3. Logout")
		fmt.Print("Pilih menu: ")

		a.scanner.Scan()
		choice := strings.TrimSpace(a.scanner.Text())

		switch choice {
		case "1":
			a.handleBrowseGames()

		case "2":
			a.handleAddGameKey()

		case "3":
			fmt.Println("\n[Sukses] Berhasil logout.")
			return false

		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

// handleBrowseGames displays all games
func (a *AdminCLI) handleBrowseGames() {
	a.gameHandler.BrowseGames()
}

// handleAddGameKey adds a new game key
func (a *AdminCLI) handleAddGameKey() {
	fmt.Print("\nMasukkan ID Game: ")
	a.scanner.Scan()
	var gameID int
	fmt.Sscanf(a.scanner.Text(), "%d", &gameID)

	fmt.Print("Masukkan Key Baru (contoh: XXXX-YYYY-ZZZZ): ")
	a.scanner.Scan()
	key := strings.TrimSpace(a.scanner.Text())

	err := a.gameService.AddKey(gameID, key)
	if err != nil {
		fmt.Printf("\n[Error] Gagal tambah key: %v\n", err)
	} else {
		fmt.Println("\n[Sukses] Key lisensi baru berhasil ditambahkan!")
	}
}
