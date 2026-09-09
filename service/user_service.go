package service

import (
	"errors"
	"game-store/entity"
	"game-store/repository"
)

type UserService interface {
	Register(email, password, fullName string) error
	Login(email, password string) (entity.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(email, password, fullName string) error {
	// Input Validate
	if email == "" || password == "" || fullName == "" {
		return errors.New("semua field (email, password, nama lengkap) harus diisi")
	}

	// Email Checker
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser.ID != 0 {
		return errors.New("email sudah terdaftar, gunakan email lain")
	}

	// Make User and user profile
	user := entity.User{
		Email:    email,
		Password: password,
		Role:     "customer",
	}

	profile := entity.UserProfile{
		FullName:      fullName,
		WalletBalance: 0.0,
	}

	// Save to database
	return s.userRepo.Register(user, profile)
}

func (s *userService) Login(email, password string) (entity.User, error) {
	// Validate Input
	if email == "" || password == "" {
		return entity.User{}, errors.New("email dan password tidak boleh kosong")
	}

	// Search user from email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return entity.User{}, errors.New("email atau password salah")
	}

	// Password Match Checker
	if user.Password != password {
		return entity.User{}, errors.New("email atau password salah")
	}

	return user, nil
}
