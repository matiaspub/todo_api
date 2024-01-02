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
	Delete(userId int, listId int) error
	Update(userId int, listId int, list entity.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item entity.TodoItem) (int, error)
	GetAll(userId int, listId int) ([]entity.TodoItem, error)
	GetOne(userId int, itemId int) (entity.TodoItem, error)
	Update(userId int, itemId int, input entity.UpdateTodoItemInput) error
	Delete(userId int, itemId int) error
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
		TodoItem:      NewTodoItemPg(db),
	}
}
