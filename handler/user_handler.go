package handler

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"game-store/entity"
	"game-store/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// HandleLogin handles user login input and validation
func (h *UserHandler) HandleLogin(email, password string) (entity.User, error) {
	// Validate input
	if email == "" || password == "" {
		return entity.User{}, fmt.Errorf("email dan password tidak boleh kosong")
	}

	// Attempt login
	user, err := h.userService.Login(email, password)
	if err != nil {
		fmt.Printf("[Error] Login Gagal: %v\n", err)
		return entity.User{}, err
	}

	fmt.Printf("[Sukses] Selamat datang, %s!\n", user.Email)
	return user, nil
}

// HandleRegister handles user registration input and validation
func (h *UserHandler) HandleRegister(email, password, fullName string) error {
	// Validate input
	if email == "" || password == "" || fullName == "" {
		return fmt.Errorf("semua field (email, password, nama lengkap) harus diisi")
	}

	// Attempt registration
	err := h.userService.Register(email, password, fullName)
	if err != nil {
		fmt.Printf("[Error] Register Gagal: %v\n", err)
		return err
	}

	fmt.Println("[Sukses] Pendaftaran berhasil! Silakan login.")
	return nil
}

// HandleTopUp handles wallet top-up input and validation
func (h *UserHandler) HandleTopUp(userID int, amount float64) error {
	// Validate amount
	if amount <= 0 {
		return fmt.Errorf("nominal top-up harus lebih dari 0")
	}

	// Attempt top-up
	err := h.userService.TopUp(userID, amount)
	if err != nil {
		fmt.Printf("[Error] Top-up gagal: %v\n", err)
		return err
	}

	fmt.Printf("[Sukses] Berhasil top-up saldo sebesar Rp%.2f!\n", amount)
	return nil
}

// PromptLogin displays login prompt and reads input
func (h *UserHandler) PromptLogin() (email, password string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	emailInput, _ := reader.ReadString('\n')
	email = strings.TrimSpace(emailInput)

	fmt.Print("Password: ")
	passwordInput, _ := reader.ReadString('\n')
	password = strings.TrimSpace(passwordInput)

	return email, password
}

// PromptRegister displays registration prompt and reads input
func (h *UserHandler) PromptRegister() (email, password, fullName string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Email: ")
	emailInput, _ := reader.ReadString('\n')
	email = strings.TrimSpace(emailInput)

	fmt.Print("Password: ")
	passwordInput, _ := reader.ReadString('\n')
	password = strings.TrimSpace(passwordInput)

	fmt.Print("Nama Lengkap: ")
	fullNameInput, _ := reader.ReadString('\n')
	fullName = strings.TrimSpace(fullNameInput)

	return email, password, fullName
}

// PromptTopUp displays top-up prompt and reads input
func (h *UserHandler) PromptTopUp() float64 {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Masukkan Nominal Top-up: Rp")
	amountInput, _ := reader.ReadString('\n')
	amountInput = strings.TrimSpace(amountInput)

	var amount float64
	fmt.Sscanf(amountInput, "%f", &amount)

	return amount
}
