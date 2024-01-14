package service

import (
	"github.com/matiaspub/todo-api/pkg/entity"
	"github.com/matiaspub/todo-api/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, todoList entity.TodoList) (int, error)
	GetAll(userId int) ([]entity.TodoList, error)
	GetOne(userId int, listId int) (entity.TodoList, error)
	Delete(userId int, listId int) error
	Update(userId int, listId int, list entity.UpdateListInput) error
}

type TodoItem interface {
	Create(userId int, listId int, item entity.TodoItem) (int, error)
	GetAll(userId int, listId int) ([]entity.TodoItem, error)
	GetOne(userId int, itemId int) (entity.TodoItem, error)
	Update(userId int, itemId int, input entity.UpdateTodoItemInput) error
	Delete(userId int, itemId int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
