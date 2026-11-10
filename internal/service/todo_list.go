package service

import (
	"errors"
	"fmt"

	todo "github.com/slemeshaev/todo/internal/models"
	"github.com/slemeshaev/todo/internal/repository"
)

var (
	// ErrListNotFound is returned when the list does not exist or belongs to another user.
	ErrListNotFound = errors.New("list not found")

	// ErrNothingToUpdate is returned when the update request sets no fields.
	ErrNothingToUpdate = errors.New("nothing to update")
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
	list, err := s.repo.GetById(userId, listId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return todo.TodoList{}, ErrListNotFound
		}
		return todo.TodoList{}, fmt.Errorf("get list: %w", err)
	}

	return list, nil
}

func (s *TodoListService) Delete(userId, listId int) error {
	if err := s.repo.Delete(userId, listId); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrListNotFound
		}
		return fmt.Errorf("delete list: %w", err)
	}

	return nil
}

func (s *TodoListService) Update(userId, listId int, input todo.UpdateListInput) error {
	if input.Title == nil && input.Description == nil {
		return ErrNothingToUpdate
	}

	if err := s.repo.Update(userId, listId, input); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrListNotFound
		}
		return fmt.Errorf("update list: %w", err)
	}

	return nil
}
