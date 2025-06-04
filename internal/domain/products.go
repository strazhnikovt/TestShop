package domain

import (
	"github.com/lib/pq"
)

type Product struct {
	ID          int            `db:"id" json:"id,omitempty"`
	Description string         `db:"description" json:"description,omitempty"`
	Tags        pq.StringArray `db:"tags"        json:"tags,omitempty"`
	Quantity    int            `db:"quantity"    json:"quantity,omitempty"`
	Price       float64        `db:"price"       json:"price,omitempty"`
}

type ProductCreateRequest struct {
	Description string   `json:"description" validate:"required"`
	Tags        []string `json:"tags"`
	Quantity    int      `json:"quantity"    validate:"min=0"`
	Price       float64  `json:"price"       validate:"min=0"`
}
