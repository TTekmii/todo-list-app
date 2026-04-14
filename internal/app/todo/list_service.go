package todo

import (
	"context"
	"fmt"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
)

type TodoListService struct {
	listRepo repo.TodoList
}

func NewTodoListService(listRepo repo.TodoList) *TodoListService {
	return &TodoListService{
		listRepo: listRepo,
	}
}

func (s *TodoListService) Create(ctx context.Context, userID int, list model.TodoList) (int, error) {
	if list.Title == "" {
		return 0, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}

	listID, err := s.listRepo.Create(ctx, userID, list)
	if err != nil {
		return 0, fmt.Errorf("failed to create list: %w", err)
	}

	return listID, nil
}

func (s *TodoListService) GetAll(ctx context.Context, userID int) ([]model.TodoList, error) {
	lists, err := s.listRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lists: %w", err)
	}

	return lists, nil
}

func (s *TodoListService) GetById(ctx context.Context, userID, listID int) (model.TodoList, error) {
	list, err := s.listRepo.GetById(ctx, userID, listID)
	if err != nil {
		return model.TodoList{}, fmt.Errorf("failed to fetch list: %w", err)
	}

	return list, nil
}

func (s *TodoListService) Delete(ctx context.Context, userID, listID int) error {
	err := s.listRepo.Delete(ctx, userID, listID)
	if err != nil {
		return fmt.Errorf("failed to delete list: %d: %w", listID, err)
	}

	return nil
}

func (s *TodoListService) Update(ctx context.Context, userID, listID int, input model.UpdateListInput) error {
	if !input.HasChanges() {
		return ErrNoFieldsToUpdate
	}

	if input.Title != nil && *input.Title == "" {
		return fmt.Errorf("title cannot be empty: %w", ErrInvalidInput)
	}

	err := s.listRepo.Update(ctx, userID, listID, input)
	if err != nil {
		return fmt.Errorf("failed to update list: %d: %w", listID, err)
	}

	return nil
}
