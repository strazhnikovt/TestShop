package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/strazhnikovt/TestShop/internal/domain"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	query := `
        INSERT INTO products (description, tags, quantity, price) 
        VALUES ($1, $2, $3, $4) RETURNING id`
	return r.db.QueryRow(
		query,
		product.Description,
		pq.Array(product.Tags),
		product.Quantity,
		product.Price,
	).Scan(&product.ID)
}

func (r *ProductRepository) Update(product *domain.Product) error {
	query := `
        UPDATE products 
           SET description = $1, tags = $2, quantity = $3, price = $4 
         WHERE id = $5`
	_, err := r.db.Exec(
		query,
		product.Description,
		pq.Array(product.Tags),
		product.Quantity,
		product.Price,
		product.ID,
	)
	return err
}

func (r *ProductRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *ProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	query := `SELECT id, description, tags, quantity, price FROM products`
	if err := r.db.Select(&products, query); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*domain.Product, error) {
	var product domain.Product
	query := `SELECT id, description, tags, quantity, price FROM products WHERE id = $1`
	err := r.db.Get(&product, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}
