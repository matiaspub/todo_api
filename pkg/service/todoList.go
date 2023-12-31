package service

import (
	"github.com/matiaspub/todo-api/pkg/entity"
	"github.com/matiaspub/todo-api/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{
		repo: repo,
	}
}

func (s *TodoListService) Create(userId int, todo entity.TodoList) (int, error) {
	return s.repo.Create(userId, todo)
}

func (s *TodoListService) GetAll(userId int) ([]entity.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetOne(userId int, listId int) (entity.TodoList, error) {
	return s.repo.GetOne(userId, listId)
}

func (s *TodoListService) Delete(userId int, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *TodoListService) Update(userId int, listId int, list entity.UpdateListInput) error {
	if err := list.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, list)
}
