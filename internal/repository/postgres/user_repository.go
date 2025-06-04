package postgres

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/strazhnikovt/TestShop/internal/domain"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
        INSERT INTO users 
            (first_name, last_name, login, full_name, age, is_married, password, role)
        VALUES 
            ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id`
	return r.db.QueryRow(
		query,
		user.FirstName,
		user.LastName,
		user.Login,
		user.FullName,
		user.Age,
		user.IsMarried,
		user.Password,
		user.Role,
	).Scan(&user.ID)
}

func (r *UserRepository) GetByLogin(login string) (*domain.User, error) {
	query := `
        SELECT 
            id, first_name, last_name, login, full_name, age, is_married, password, role 
        FROM users 
        WHERE login = $1`

	var user domain.User
	err := r.db.Get(&user, query, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
