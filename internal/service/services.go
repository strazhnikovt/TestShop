package service

import (
	"github.com/strazhnikovt/TestShop/internal/repository/postgres"
	"github.com/strazhnikovt/TestShop/pkg/auth"
)

type Services struct {
	User       *UserService
	Product    *ProductService
	Order      *OrderService
	JWTManager *auth.JWTManager
}

func NewServices(
	repos *postgres.Repositories,
	jwtSecret string,
) *Services {
	jwtManager := auth.NewJWTManager(jwtSecret)
	return &Services{
		User:       NewUserService(repos.User),
		Product:    NewProductService(repos.Product),
		Order:      NewOrderService(repos.Order, repos.Product),
		JWTManager: jwtManager,
	}
}
