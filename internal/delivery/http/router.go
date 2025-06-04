package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/strazhnikovt/TestShop/internal/service"
	"github.com/strazhnikovt/TestShop/pkg/auth"
)

type Handlers struct {
	User    *UserHandler
	Product *ProductHandler
	Order   *OrderHandler
}

func NewHandlers(
	userService *service.UserService,
	productService *service.ProductService,
	orderService *service.OrderService,
	jwtManager *auth.JWTManager,
) *Handlers {
	return &Handlers{
		User:    NewUserHandler(userService, jwtManager),
		Product: NewProductHandler(productService),
		Order:   NewOrderHandler(orderService),
	}
}

func RegisterRoutes(r *chi.Mux, h *Handlers, jwtManager *auth.JWTManager) {
	// Public routes
	r.Post("/register", h.User.Register)
	r.Post("/login", h.User.Login)
	r.Get("/products", h.Product.GetAllProducts)

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(AuthMiddleware(jwtManager))

		r.Route("/orders", func(r chi.Router) {
			r.Post("/", h.Order.CreateOrder)
		})

		// Admin routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(AdminOnlyMiddleware)

			r.Post("/products", h.Product.CreateProduct)
			r.Put("/products/{id}", h.Product.UpdateProduct)
			r.Delete("/products/{id}", h.Product.DeleteProduct)
		})
	})
}
