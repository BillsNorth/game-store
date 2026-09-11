package cli

import (
	"bufio"
	"fmt"
	"strings"

	"game-store/entity"
	"game-store/handler"
)

type AuthCLI struct {
	userHandler *handler.UserHandler
	scanner     *bufio.Scanner
}

func NewAuthCLI(userHandler *handler.UserHandler, scanner *bufio.Scanner) *AuthCLI {
	return &AuthCLI{
		userHandler: userHandler,
		scanner:     scanner,
	}
}

// ShowAuthMenu displays login and register options
func (a *AuthCLI) ShowAuthMenu() (*entity.User, bool) {
	for {
		fmt.Println("\n=================================")
		fmt.Println("     WELCOME TO GAME STORE CLI   ")
		fmt.Println("=================================")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Keluar")
		fmt.Print("Pilih menu: ")

		a.scanner.Scan()
		choice := strings.TrimSpace(a.scanner.Text())

		switch choice {
		case "1":
			user, loggedIn := a.showLoginMenu()
			if loggedIn {
				return user, true
			}

		case "2":
			a.showRegisterMenu()

		case "3":
			fmt.Println("Terima kasih telah menggunakan Game Store CLI!")
			return nil, false

		default:
			fmt.Println("Pilihan tidak valid.")
		}
	}
}

// showLoginMenu handles login UI and delegates to handler
func (a *AuthCLI) showLoginMenu() (*entity.User, bool) {
	fmt.Println("\n========== LOGIN ==========")
	email, password := a.userHandler.PromptLogin()

	user, err := a.userHandler.HandleLogin(email, password)
	if err != nil {
		return nil, false
	}

	return &user, true
}

// showRegisterMenu handles registration UI and delegates to handler
func (a *AuthCLI) showRegisterMenu() {
	fmt.Println("\n========== REGISTER ==========")
	email, password, fullName := a.userHandler.PromptRegister()

	a.userHandler.HandleRegister(email, password, fullName)
}
