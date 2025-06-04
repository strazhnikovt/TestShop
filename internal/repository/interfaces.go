package repository

import "github.com/strazhnikovt/TestShop/internal/domain"

type UserRepository interface {
	Create(user *domain.User) error
	GetByLogin(login string) (*domain.User, error)
}

type ProductRepository interface {
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(id int) error
	GetAll() ([]domain.Product, error)
	GetByID(id int) (*domain.Product, error)
}

type OrderRepository interface {
	Create(order *domain.Order) error
}

type Repositories struct {
	User    UserRepository
	Product ProductRepository
	Order   OrderRepository
}
