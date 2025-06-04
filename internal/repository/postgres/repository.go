package postgres

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/strazhnikovt/TestShop/internal/repository"
)

// Connect opens a PostgreSQL connection using sqlx.
func Connect(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}

// RunMigrations executes the SQL statements from the migration file.
func RunMigrations(databaseURL string) error {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	migrationSQL, err := os.ReadFile("migrations/00001_init_schema.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		return err
	}
	return nil
}

// Repositories holds implementations of all repository interfaces.
type Repositories struct {
	User    repository.UserRepository
	Product repository.ProductRepository
	Order   repository.OrderRepository
}

// NewRepositories constructs a Repositories struct from an open *sqlx.DB.
func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		User:    NewUserRepository(db),
		Product: NewProductRepository(db),
		Order:   NewOrderRepository(db),
	}
}
