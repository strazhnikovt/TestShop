package service

import (
	"github.com/strazhnikovt/TestShop/internal/domain"
	"github.com/strazhnikovt/TestShop/internal/repository"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *domain.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProduct(product *domain.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetByID(id int) (*domain.Product, error) {
	return s.repo.GetByID(id)
}
