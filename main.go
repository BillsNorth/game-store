package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"game-store/config"
	"game-store/entity"
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

	scanner := bufio.NewScanner(os.Stdin)

	var activeUser *entity.User

	for {
		if activeUser == nil {
			// LOGIN / REGISTER MENU
			fmt.Println("\n=================================")
			fmt.Println("     WELCOME TO GAME STORE CLI   ")
			fmt.Println("=================================")
			fmt.Println("1. Login")
			fmt.Println("2. Register")
			fmt.Println("3. Keluar")
			fmt.Print("Pilih menu: ")

			scanner.Scan()
			choice := strings.TrimSpace(scanner.Text())

			switch choice {
			case "1":
				fmt.Print("Email: ")
				scanner.Scan()
				email := strings.TrimSpace(scanner.Text())

				fmt.Print("Password: ")
				scanner.Scan()
				password := strings.TrimSpace(scanner.Text())

				user, err := userService.Login(email, password)
				if err != nil {
					fmt.Printf("\n[Error] Login Gagal: %v\n", err)
				} else {
					activeUser = &user
					fmt.Printf("\n[Sukses] Selamat datang, %s!\n", activeUser.Email)
				}

			case "2":
				fmt.Print("Email: ")
				scanner.Scan()
				email := strings.TrimSpace(scanner.Text())

				fmt.Print("Password: ")
				scanner.Scan()
				password := strings.TrimSpace(scanner.Text())

				fmt.Print("Nama Lengkap: ")
				scanner.Scan()
				fullName := strings.TrimSpace(scanner.Text())

				err := userService.Register(email, password, fullName)
				if err != nil {
					fmt.Printf("\n[Error] Register Gagal: %v\n", err)
				} else {
					fmt.Println("\n[Sukses] Pendaftaran berhasil! Silakan login.")
				}

			case "3":
				fmt.Println("Terima kasih telah menggunakan Game Store CLI!")
				return

			default:
				fmt.Println("Pilihan tidak valid.")
			}

		} else {
			// Check user role and show appropriate menu
			if activeUser.Role == "admin" {
				// Admin Menu
				fmt.Printf("\n=================================\n")
				fmt.Printf(" MENU ADMIN (%s)\n", activeUser.Email)
				fmt.Println("=================================")
				fmt.Println("1. Lihat Katalog Game")
				fmt.Println("2. Tambah Stok Game Key Baru")
				fmt.Println("3. Logout")
				fmt.Print("Pilih menu: ")

				scanner.Scan()
				choice := strings.TrimSpace(scanner.Text())

				switch choice {
				case "1":
					games, _ := gameService.GetAllGames()
					fmt.Println("\n--- KATALOG GAME ---")
					for _, g := range games {
						fmt.Printf("ID: %d | Judul: %s | Stok: %d key\n", g.Game.ID, g.Game.Title, g.Stock)
					}
				case "2":
					fmt.Print("Masukkan ID Game: ")
					scanner.Scan()
					var gameID int
					fmt.Sscanf(scanner.Text(), "%d", &gameID)

					fmt.Print("Masukkan Key Baru (contoh: XXXX-YYYY-ZZZZ): ")
					scanner.Scan()
					key := strings.TrimSpace(scanner.Text())

					err := gameService.AddKey(gameID, key)
					if err != nil {
						fmt.Printf("\n[Error] Gagal tambah key: %v\n", err)
					} else {
						fmt.Println("\n[Sukses] Key lisensi baru berhasil ditambahkan!")
					}
				case "3":
					activeUser = nil
					fmt.Println("\n[Sukses] Berhasil logout.")
				}

			} else {
				// Customer Menu
				fmt.Printf("\n=================================\n")
				fmt.Printf(" DASHBOARD CUSTOMER (%s)\n", activeUser.Email)
				fmt.Println("=================================")
				fmt.Println("1. Lihat Katalog Game")
				fmt.Println("2. Tambah Game ke Keranjang")
				fmt.Println("3. Lihat Keranjang Belanja")
				fmt.Println("4. Checkout (Beli)")
				fmt.Println("5. Lihat Riwayat Pesanan & Key")
				fmt.Println("6. Top-up Saldo Wallet")
				fmt.Println("7. Logout")
				fmt.Print("Pilih menu: ")

				scanner.Scan()
				choice := strings.TrimSpace(scanner.Text())

				switch choice {
				case "1":
					games, _ := gameService.GetAllGames()
					fmt.Println("\n--- KATALOG GAME ---")
					fmt.Printf("%-4s | %-30s | %-12s | %-12s | %-6s\n", "ID", "Judul Game", "Kategori", "Harga", "Stok")
					fmt.Println(strings.Repeat("-", 75))
					for _, g := range games {
						fmt.Printf("%-4d | %-30s | %-12s | Rp%-10.2f | %d key\n", g.Game.ID, g.Game.Title, g.Game.CategoryName, g.Game.Price, g.Stock)
					}
				case "2":
					fmt.Print("Masukkan ID Game: ")
					scanner.Scan()
					var gameID int
					fmt.Sscanf(scanner.Text(), "%d", &gameID)

					fmt.Print("Masukkan Jumlah: ")
					scanner.Scan()
					var qty int
					fmt.Sscanf(scanner.Text(), "%d", &qty)

					err := cartService.AddToCart(activeUser.ID, gameID, qty)
					if err != nil {
						fmt.Printf("\n[Error] Gagal masuk keranjang: %v\n", err)
					} else {
						fmt.Println("\n[Sukses] Game berhasil ditambahkan ke keranjang!")
					}
				case "3":
					summary, err := cartService.GetCart(activeUser.ID)
					if err != nil {
						fmt.Printf("\n[Error] Gagal mengambil keranjang: %v\n", err)
						continue
					}
					if len(summary.Items) == 0 {
						fmt.Println("\nKeranjang belanja kamu masih kosong.")
						continue
					}
					fmt.Println("\n--- ISI KERANJANG BELANJA ---")
					fmt.Printf("%-4s | %-30s | %-12s | %-6s | %-12s\n", "ID", "Judul Game", "Harga Satuan", "Qty", "Total")
					fmt.Println(strings.Repeat("-", 75))
					for _, item := range summary.Items {
						subtotal := item.Price * float64(item.Quantity)
						fmt.Printf("%-4d | %-30s | Rp%-10.2f | %-6d | Rp%-10.2f\n", item.ID, item.GameTitle, item.Price, item.Quantity, subtotal)
					}
					fmt.Println(strings.Repeat("-", 75))
					fmt.Printf("GRAND TOTAL: Rp%.2f\n", summary.GrandTotal)
				case "4":
					orderID, err := orderService.Checkout(activeUser.ID)
					if err != nil {
						fmt.Printf("\n[Error] Checkout gagal: %v\n", err)
					} else {
						fmt.Printf("\n[Sukses] Checkout Berhasil! Nomor Order ID: #%d\n", orderID)
						fmt.Println("Lisensi game dapat dilihat di menu Riwayat Pesanan.")
					}
				case "5":
					orders, err := orderService.GetUserOrders(activeUser.ID)
					if err != nil {
						fmt.Printf("\n[Error] Gagal mengambil riwayat: %v\n", err)
						continue
					}
					if len(orders) == 0 {
						fmt.Println("\nKamu belum pernah melakukan transaksi.")
						continue
					}
					fmt.Println("\n--- RIWAYAT TRANSAKSI ---")
					for _, o := range orders {
						fmt.Printf("\nOrder ID: #%d | Tanggal: %s | Total: Rp%.2f\n", o.ID, o.CreatedAt, o.TotalPrice)
						details, err := orderService.GetOrderDetails(o.ID)
						if err == nil {
							for _, d := range details {
								fmt.Printf("  -> Game: %-25s | Key: %s\n", d.GameTitle, d.LicenseKey)
							}
						}
					}
				case "6":
					fmt.Print("Masukkan Nominal Top-up: Rp")
					scanner.Scan()
					var amount float64
					fmt.Sscanf(scanner.Text(), "%f", &amount)

					err := userService.TopUp(activeUser.ID, amount)
					if err != nil {
						fmt.Printf("\n[Error] Top-up gagal: %v\n", err)
					} else {
						fmt.Printf("\n[Sukses] Berhasil top-up saldo sebesar Rp%.2f!\n", amount)
					}
				case "7":
					activeUser = nil
					fmt.Println("\n[Sukses] Berhasil logout.")
				default:
					fmt.Println("Pilihan tidak valid.")
				}
			}
		}
	}
}
