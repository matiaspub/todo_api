package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/matiaspub/todo-api/pkg/entity"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username string, password string) (entity.User, error)
}

type TodoList interface {
	Create(userId int, todo entity.TodoList) (int, error)
	GetAll(userId int) ([]entity.TodoList, error)
	GetOne(userId int, listId int) (entity.TodoList, error)
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
		TodoList:      NewTodoListPg(db),
	}
}
