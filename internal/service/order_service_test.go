package service

import (
	"errors"
	"testing"

	"github.com/strazhnikovt/TestShop/internal/domain"
	"github.com/strazhnikovt/TestShop/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_CreateOrder(t *testing.T) {
	t.Run("ValidOrder", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		orderService := NewOrderService(mockOrderRepo, mockProductRepo)

		order := &domain.Order{
			UserID: 1,
			OrderItems: []domain.OrderItem{
				{ProductID: 1, Quantity: 2},
				{ProductID: 2, Quantity: 1},
			},
		}

		mockProductRepo.
			On("GetByID", 1).
			Return(&domain.Product{ID: 1, Quantity: 10}, nil)
		mockProductRepo.
			On("GetByID", 2).
			Return(&domain.Product{ID: 2, Quantity: 5}, nil)

		mockOrderRepo.
			On("Create", mock.MatchedBy(func(o *domain.Order) bool {
				return !o.CreatedAt.IsZero() &&
					o.UserID == 1 &&
					len(o.OrderItems) == 2
			})).
			Return(nil)

		err := orderService.CreateOrder(order)
		assert.NoError(t, err)

		mockProductRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("InsufficientQuantity", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		orderService := NewOrderService(mockOrderRepo, mockProductRepo)

		order := &domain.Order{
			UserID: 1,
			OrderItems: []domain.OrderItem{
				{ProductID: 1, Quantity: 20}, // request more than available
			},
		}

		mockProductRepo.
			On("GetByID", 1).
			Return(&domain.Product{ID: 1, Quantity: 10}, nil)

		err := orderService.CreateOrder(order)
		assert.Error(t, err)
		assert.Equal(t, "insufficient product quantity", err.Error())

		mockOrderRepo.AssertNotCalled(t, "Create")
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		mockOrderRepo := new(mocks.MockOrderRepository)
		mockProductRepo := new(mocks.MockProductRepository)
		orderService := NewOrderService(mockOrderRepo, mockProductRepo)

		order := &domain.Order{
			UserID: 1,
			OrderItems: []domain.OrderItem{
				{ProductID: 999, Quantity: 1}, // non-existent product
			},
		}

		mockProductRepo.
			On("GetByID", 999).
			Return((*domain.Product)(nil), errors.New("not found"))

		err := orderService.CreateOrder(order)
		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())

		mockOrderRepo.AssertNotCalled(t, "Create")
		mockProductRepo.AssertExpectations(t)
	})
}
