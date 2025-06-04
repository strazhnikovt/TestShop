package domain

import "time"

type Order struct {
	ID         int         `db:"id" json:"id,omitempty"`
	UserID     int         `db:"user_id" json:"user_id,omitempty"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	OrderItems []OrderItem `db:"-" json:"order_items,omitempty"`
}

type OrderItem struct {
	ID           int     `db:"id" json:"id,omitempty"`
	OrderID      int     `db:"order_id" json:"order_id,omitempty"`
	ProductID    int     `db:"product_id" json:"product_id,omitempty"`
	Quantity     int     `db:"quantity" json:"quantity,omitempty"`
	PriceAtOrder float64 `db:"price_at_order" json:"price_at_order,omitempty"`
}

type OrderCreateRequest struct {
	UserID int                `json:"user_id" validate:"required"`
	Items  []OrderItemRequest `json:"items" validate:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID int `json:"product_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"min=1"`
}
