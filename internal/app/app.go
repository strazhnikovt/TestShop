package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/strazhnikovt/TestShop/internal/config"
	httpdelivery "github.com/strazhnikovt/TestShop/internal/delivery/http"
	"github.com/strazhnikovt/TestShop/internal/domain"
	"github.com/strazhnikovt/TestShop/internal/repository"
	"github.com/strazhnikovt/TestShop/internal/repository/postgres"
	"github.com/strazhnikovt/TestShop/internal/service"
	"github.com/strazhnikovt/TestShop/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Config   *config.Config
	Logger   *logging.Logger
	Server   *http.Server
	Services *service.Services
}

func New() (*App, error) {
	cfg := config.Load()
	logger := logging.NewLogger()

	db, err := postgres.Connect(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := postgres.RunMigrations(cfg.DatabaseURL); err != nil {
		return nil, fmt.Errorf("migrations failed: %w", err)
	}

	repos := postgres.NewRepositories(db)

	// ensure a default admin user exists
	if err := createDefaultAdmin(repos.User, cfg.AdminLogin, cfg.AdminPass); err != nil {
		return nil, fmt.Errorf("failed to create default admin: %w", err)
	}

	// build services with repository collection and JWT secret
	services := service.NewServices(repos, cfg.JWTSecret)
	handlers := httpdelivery.NewHandlers(
		services.User,
		services.Product,
		services.Order,
		services.JWTManager,
	)

	router := chi.NewRouter()
	router.Use(httpdelivery.LoggingMiddleware(logger))
	httpdelivery.RegisterRoutes(router, handlers, services.JWTManager)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	return &App{
		Config:   cfg,
		Logger:   logger,
		Server:   server,
		Services: services,
	}, nil
}

func (a *App) Run() error {
	a.Logger.Printf("Server started on port %d", a.Config.Port)
	return a.Server.ListenAndServe()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return a.Server.Shutdown(ctx)
}

// createDefaultAdmin ensures the admin user exists using interface type
func createDefaultAdmin(repo repository.UserRepository, login, password string) error {
	admin, err := repo.GetByLogin(login)
	if err != nil {
		return err
	}
	if admin == nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		adminUser := &domain.User{
			FirstName: "Admin",
			LastName:  "System",
			Login:     login,
			FullName:  "Admin System",
			Age:       30,
			Password:  string(hashedPassword),
			Role:      "admin",
		}
		if err := repo.Create(adminUser); err != nil {
			return err
		}
	}
	return nil
}
