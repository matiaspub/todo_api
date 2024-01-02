package service

import (
	"github.com/matiaspub/todo-api/pkg/entity"
	"github.com/matiaspub/todo-api/pkg/repository"
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func (s *TodoItemService) Delete(userId int, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId int, itemId int, input entity.UpdateTodoItemInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, itemId, input)
}

func (s *TodoItemService) GetOne(userId int, itemId int) (entity.TodoItem, error) {
	return s.repo.GetOne(userId, itemId)
}

func (s *TodoItemService) GetAll(userId int, listId int) ([]entity.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) Create(userId int, listId int, item entity.TodoItem) (int, error) {
	_, err := s.listRepo.GetOne(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}
