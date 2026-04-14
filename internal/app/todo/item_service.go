package todo

import (
	"context"
	"fmt"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
)

type TodoItemService struct {
	itemRepo repo.TodoItem
	listRepo repo.TodoList
}

func NewTodoItemService(itemRepo repo.TodoItem, listRepo repo.TodoList) *TodoItemService {
	return &TodoItemService{
		itemRepo: itemRepo,
		listRepo: listRepo,
	}
}

func (s *TodoItemService) Create(ctx context.Context, userID, listID int, item model.TodoItem) (int, error) {
	if item.Title == "" {
		return 0, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}
	_, err := s.listRepo.GetById(ctx, userID, listID)
	if err != nil {
		return 0, fmt.Errorf("access denied to list %d: %w", listID, err)
	}

	itemID, err := s.itemRepo.Create(ctx, listID, item)
	if err != nil {
		return 0, fmt.Errorf("failed to create item: %w", err)
	}

	return itemID, nil
}

func (s *TodoItemService) GetAll(ctx context.Context, userID, listID int) ([]model.TodoItem, error) {
	_, err := s.listRepo.GetById(ctx, userID, listID)
	if err != nil {
		return nil, fmt.Errorf("access denied to list %d: %w", listID, err)
	}

	items, err := s.itemRepo.GetAll(ctx, userID, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	return items, nil
}

func (s *TodoItemService) GetById(ctx context.Context, userID, itemID int) (model.TodoItem, error) {
	item, err := s.itemRepo.GetById(ctx, userID, itemID)
	if err != nil {
		return model.TodoItem{}, fmt.Errorf("failed to fetch item %d: %w", itemID, err)
	}

	return item, nil
}

func (s *TodoItemService) Delete(ctx context.Context, userID, itemID int) error {
	err := s.itemRepo.Delete(ctx, userID, itemID)
	if err != nil {
		return fmt.Errorf("failed to delete item %d: %w", itemID, err)
	}
	return nil
}

func (s *TodoItemService) Update(ctx context.Context, userID, itemID int, input model.UpdateItemInput) error {
	if !input.HasChanges() {
		return nil
	}

	if input.Title != nil && *input.Title == "" {
		return fmt.Errorf("title cannot be empty: %w", ErrInvalidInput)
	}

	err := s.itemRepo.Update(ctx, userID, itemID, input)
	if err != nil {
		return fmt.Errorf("failed to update item %d: %w", itemID, err)
	}
	return nil
}
