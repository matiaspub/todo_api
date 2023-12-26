package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/matiaspub/todo-api/pkg/entity"
)

type AuthPg struct {
	db *sqlx.DB
}

func NewAuthPg(db *sqlx.DB) *AuthPg {
	return &AuthPg{db: db}
}

func (a *AuthPg) CreateUser(user entity.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO users (name, username, password_hash) values ($1, $2, $3) RETURNING id")
	row := a.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}