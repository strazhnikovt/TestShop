package service

import (
	"errors"
	"testing"

	"github.com/strazhnikovt/TestShop/internal/domain"
	"github.com/strazhnikovt/TestShop/internal/service/mocks"

	"github.com/stretchr/testify/assert"
)

func TestProductService_CreateProduct(t *testing.T) {
	t.Run("ValidProduct", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		productService := NewProductService(mockRepo)

		product := &domain.Product{
			Description: "Laptop",
			Tags:        []string{"electronics", "tech"},
			Quantity:    10,
			Price:       999.99,
		}

		mockRepo.On("Create", product).Return(nil)

		err := productService.CreateProduct(product)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		productService := NewProductService(mockRepo)

		product := &domain.Product{
			Description: "Laptop",
			Tags:        []string{"electronics"},
			Quantity:    10,
			Price:       999.99,
		}

		mockRepo.On("Create", product).Return(errors.New("db error"))

		err := productService.CreateProduct(product)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_UpdateProduct(t *testing.T) {
	t.Run("ValidUpdate", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		productService := NewProductService(mockRepo)

		product := &domain.Product{
			ID:          1,
			Description: "Updated Laptop",
			Tags:        []string{"updated"},
			Quantity:    5,
			Price:       1099.99,
		}

		mockRepo.On("Update", product).Return(nil)

		err := productService.UpdateProduct(product)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateError", func(t *testing.T) {
		mockRepo := new(mocks.MockProductRepository)
		productService := NewProductService(mockRepo)

		product := &domain.Product{ID: 1}

		mockRepo.On("Update", product).Return(errors.New("update error"))

		err := productService.UpdateProduct(product)
		assert.Error(t, err)
		assert.Equal(t, "update error", err.Error())

		mockRepo.AssertExpectations(t)
	})
}
