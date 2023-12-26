package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/matiaspub/todo-api/pkg/entity"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPg(db),
	}
}
