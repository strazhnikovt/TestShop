package domain

type User struct {
	ID        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Login     string `db:"login" json:"login"`
	FullName  string `db:"full_name" json:"full_name"`
	Age       int    `db:"age" json:"age"`
	IsMarried bool   `db:"is_married" json:"is_married"`
	Password  string `db:"password" json:"-"`
	Role      string `db:"role" json:"role"`
}

type UserRegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Login     string `json:"login" validate:"required"`
	Age       int    `json:"age" validate:"required,min=18"`
	IsMarried bool   `json:"is_married"`
	Password  string `json:"password" validate:"required,min=8"`
}

type UserLoginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
