package cli

import (
	"bufio"
	"os"

	"game-store/handler"
	"game-store/service"
)

type CLI struct {
	authCLI      *AuthCLI
	gameHandler  *handler.GameHandler
	cartHandler  *handler.CartHandler
	orderHandler *handler.OrderHandler
	userHandler  *handler.UserHandler
	gameService  service.GameService
	scanner      *bufio.Scanner
}

func NewCLI(
	gameHandler *handler.GameHandler,
	cartHandler *handler.CartHandler,
	orderHandler *handler.OrderHandler,
	userHandler *handler.UserHandler,
	gameService service.GameService,
) *CLI {
	scanner := bufio.NewScanner(os.Stdin)
	authCLI := NewAuthCLI(userHandler, scanner)

	return &CLI{
		authCLI:      authCLI,
		gameHandler:  gameHandler,
		cartHandler:  cartHandler,
		orderHandler: orderHandler,
		userHandler:  userHandler,
		gameService:  gameService,
		scanner:      scanner,
	}
}

// Start initiates the CLI application
func (cli *CLI) Start() {
	for {
		// Show authentication menu
		user, shouldContinue := cli.authCLI.ShowAuthMenu()
		if !shouldContinue {
			break
		}

		// Handle menu based on user role
		if user.Role == "admin" {
			adminCLI := NewAdminCLI(cli.gameHandler, cli.gameService, cli.scanner, user.Email)
			if !adminCLI.ShowMenu() {
				// Logout, go back to auth menu
				continue
			}
		} else {
			customerCLI := NewCustomerCLI(
				cli.gameHandler,
				cli.cartHandler,
				cli.orderHandler,
				cli.userHandler,
				cli.scanner,
				user.ID,
				user.Email,
			)
			if !customerCLI.ShowMenu() {
				// Logout, go back to auth menu
				continue
			}
		}
	}
}
