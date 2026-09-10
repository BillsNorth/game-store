package service

import (
	"errors"
	"fmt"
	"game-store/entity"
	"game-store/repository"
)

// CartSummary holds a list of items and the overall price total in the cart
type CartSummary struct {
	Items      []entity.Cart
	GrandTotal float64
}

type CartService interface {
	AddToCart(userID, gameID, quantity int) error
	GetCart(userID int) (CartSummary, error)
	RemoveFromCart(cartID, userID int) error
}

type cartService struct {
	cartRepo repository.CartRepository
	gameRepo repository.GameRepository
}

func NewCartService(cartRepo repository.CartRepository, gameRepo repository.GameRepository) CartService {
	return &cartService{
		cartRepo: cartRepo,
		gameRepo: gameRepo,
	}
}

// AddToCart validates game stock before adding it to the cart.
func (s *cartService) AddToCart(userID, gameID, quantity int) error {
	if quantity <= 0 {
		return errors.New("jumlah pembelian minimal 1")
	}

	game, err := s.gameRepo.GetByID(gameID)
	if err != nil {
		return errors.New("game tidak ditemukan")
	}

	availableStock, err := s.gameRepo.GetAvailableKeyCount(game.ID)
	if err != nil {
		return err
	}

	if availableStock == 0 {
		return errors.New("stok game ini sedang habis")
	}

	userCart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return err
	}

	currentQtyInCart := 0
	for _, item := range userCart {
		if item.GameID == gameID {
			currentQtyInCart = item.Quantity
			break
		}
	}

	if (currentQtyInCart + quantity) > availableStock {
		return fmt.Errorf("stok tidak mencukupi. Sisa stok: %d, di keranjangmu sudah ada: %d", availableStock, currentQtyInCart)
	}

	return s.cartRepo.AddToCart(userID, gameID, quantity)
}

// GetCart retrieves the list of items in the cart along with the total price.
func (s *cartService) GetCart(userID int) (CartSummary, error) {
	items, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return CartSummary{}, err
	}

	var grandTotal float64
	for _, item := range items {
		grandTotal += item.Price * float64(item.Quantity)
	}

	return CartSummary{
		Items:      items,
		GrandTotal: grandTotal,
	}, nil
}

// RemoveFromCart removes an item from the cart.
func (s *cartService) RemoveFromCart(cartID, userID int) error {
	return s.cartRepo.DeleteCartItem(cartID, userID)
}
