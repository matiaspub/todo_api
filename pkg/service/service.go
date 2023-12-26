package service

import (
	"github.com/matiaspub/todo-api/pkg/entity"
	"github.com/matiaspub/todo-api/pkg/repository"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
