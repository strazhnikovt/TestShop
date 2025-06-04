package service

import (
	"errors"
	"time"

	"github.com/strazhnikovt/TestShop/internal/domain"
	"github.com/strazhnikovt/TestShop/internal/repository"
)

type OrderService struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewOrderService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(order *domain.Order) error {
	for _, item := range order.OrderItems {
		product, err := s.productRepo.GetByID(item.ProductID)
		if err != nil {
			return errors.New("product not found")
		}

		if product.Quantity < item.Quantity {
			return errors.New("insufficient product quantity")
		}
	}

	order.CreatedAt = time.Now()
	return s.orderRepo.Create(order)
}
