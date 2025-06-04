package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/strazhnikovt/TestShop/internal/domain"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *domain.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			err = fmt.Errorf("rollback error: %v, original error: %w", rbErr, err)
		}
	}()

	insertOrderQuery := `
        INSERT INTO orders (user_id, created_at) 
        VALUES ($1, $2) 
        RETURNING id`
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	if err := tx.QueryRow(insertOrderQuery, order.UserID, order.CreatedAt).Scan(&order.ID); err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	insertItemQuery := `
        INSERT INTO order_items 
            (order_id, product_id, quantity, price_at_order)
        VALUES ($1, $2, $3, (SELECT price FROM products WHERE id = $2))`
	for _, item := range order.OrderItems {
		if _, err := tx.Exec(insertItemQuery, order.ID, item.ProductID, item.Quantity); err != nil {
			return fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	updateProductQuery := `UPDATE products SET quantity = quantity - $1 WHERE id = $2`
	for _, item := range order.OrderItems {
		result, err := tx.Exec(updateProductQuery, item.Quantity, item.ProductID)
		if err != nil {
			return fmt.Errorf("failed to update product quantity: %w", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rowsAffected == 0 {
			return fmt.Errorf("no product found with ID %d", item.ProductID)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
