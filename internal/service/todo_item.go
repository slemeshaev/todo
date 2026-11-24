package service

import (
	"errors"

	todo "github.com/slemeshaev/todo/internal/models"
	"github.com/slemeshaev/todo/internal/repository"
)

var (
	// ErrItemNotFound is returned when the item does not exist or does not belong to the user.
	ErrItemNotFound = errors.New("item is not found")
)

type TodoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		// list does not exist or does not belong to user
		if errors.Is(err, repository.ErrNotFound) {
			return 0, ErrListNotFound
		}

		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	item, err := s.repo.GetById(userId, itemId)
	if errors.Is(err, repository.ErrNotFound) {
		// item does not exist or does not belong to user
		return todo.TodoItem{}, ErrItemNotFound
	}

	return item, err
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	err := s.repo.Delete(userId, itemId)
	if errors.Is(err, repository.ErrNotFound) {
		// item does not exist or does not belong to user
		return ErrItemNotFound
	}

	return err
}

func (s *TodoItemService) Update(userId, itemId int, input todo.UpdateItemInput) error {
	err := s.repo.Update(userId, itemId, input)
	if errors.Is(err, repository.ErrNotFound) {
		// item does not exist or does not belong to user
		return ErrItemNotFound
	}

	return err
}
