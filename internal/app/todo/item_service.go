package todo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
	"github.com/TTekmii/todo-list-app/internal/lib/logger/sl"
)

type TodoItemService struct {
	itemRepo repo.TodoItem
	listRepo repo.TodoList
	logger   *slog.Logger
}

func NewTodoItemService(itemRepo repo.TodoItem, listRepo repo.TodoList, logger *slog.Logger) *TodoItemService {
	return &TodoItemService{
		itemRepo: itemRepo,
		listRepo: listRepo,
		logger:   logger.With("component", "todo_item_service"),
	}
}

func (s *TodoItemService) Create(ctx context.Context, userID, listID int, item model.TodoItem) (int, error) {
	if item.Title == "" {
		s.logger.Warn("validation failed: empty title",
			slog.Int("user_id", userID),
			slog.Int("list_id", listID),
		)
		return 0, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}

	s.logger.Debug("creating item",
		slog.Int("user_id", userID),
		slog.Int("list_id", listID),
		slog.String("title", item.Title),
	)

	_, err := s.listRepo.GetById(ctx, userID, listID)
	if err != nil {
		s.logger.Warn("access denied to list",
			slog.Int("user_id", userID),
			slog.Int("list_id", listID),
			sl.Err(err),
		)
		return 0, fmt.Errorf("access denied to list %d: %w", listID, err)
	}

	itemID, err := s.itemRepo.Create(ctx, listID, item)
	if err != nil {
		s.logger.Error("failed to create item",
			slog.Int("user_id", userID),
			slog.Int("list_id", listID),
			slog.String("title", item.Title),
			sl.Err(err),
		)
		return 0, fmt.Errorf("failed to create item: %w", err)
	}

	s.logger.Info("item created successfully",
		slog.Int("item_id", itemID),
		slog.Int("user_id", userID),
		slog.Int("list_id", listID),
	)

	return itemID, nil
}

func (s *TodoItemService) GetAll(ctx context.Context, userID, listID int) ([]model.TodoItem, error) {
	s.logger.Debug("fetching items for list",
		slog.Int("user_id", userID),
		slog.Int("list_id", listID),
	)

	_, err := s.listRepo.GetById(ctx, userID, listID)
	if err != nil {
		s.logger.Warn("access denied to list",
			slog.Int("user_id", userID),
			slog.Int("list_id", listID),
			sl.Err(err),
		)
		return nil, fmt.Errorf("access denied to list %d: %w", listID, err)
	}

	items, err := s.itemRepo.GetAll(ctx, userID, listID)
	if err != nil {
		s.logger.Error("failed to fetch items",
			slog.Int("user_id", userID),
			slog.Int("list_id", listID),
			sl.Err(err),
		)
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}

	s.logger.Debug("items fetched successfully",
		slog.Int("user_id", userID),
		slog.Int("list_id", listID),
		slog.Int("count", len(items)),
	)

	return items, nil
}

func (s *TodoItemService) GetById(ctx context.Context, userID, itemID int) (model.TodoItem, error) {
	s.logger.Debug("fetching item",
		slog.Int("user_id", userID),
		slog.Int("item_id", itemID),
	)

	item, err := s.itemRepo.GetById(ctx, userID, itemID)
	if err != nil {
		s.logger.Error("failed to fetch item",
			slog.Int("user_id", userID),
			slog.Int("item_id", itemID),
			sl.Err(err),
		)
		return model.TodoItem{}, fmt.Errorf("failed to fetch item %d: %w", itemID, err)
	}

	return item, nil
}

func (s *TodoItemService) Delete(ctx context.Context, userID, itemID int) error {
	s.logger.Info("deleting item",
		slog.Int("user_id", userID),
		slog.Int("item_id", itemID),
	)

	err := s.itemRepo.Delete(ctx, userID, itemID)
	if err != nil {
		s.logger.Error("failed to delete item",
			slog.Int("user_id", userID),
			slog.Int("item_id", itemID),
			sl.Err(err),
		)
		return fmt.Errorf("failed to delete item %d: %w", itemID, err)
	}

	s.logger.Info("item deleted successfully",
		slog.Int("user_id", userID),
		slog.Int("item_id", itemID),
	)

	return nil
}

func (s *TodoItemService) Update(ctx context.Context, userID, itemID int, input model.UpdateItemInput) error {
	if !input.HasChanges() {
		s.logger.Debug("update skipped: no changes",
			slog.Int("user_id", userID),
			slog.Int("item_id", itemID),
		)
		return ErrNoFieldsToUpdate
	}

	if input.Title != nil && *input.Title == "" {
		s.logger.Warn("validation failed: empty title in update",
			slog.Int("user_id", userID),
			slog.Int("item_id", itemID),
		)
		return fmt.Errorf("title cannot be empty: %w", ErrInvalidInput)
	}

	s.logger.Debug("updating item",
		slog.Int("user_id", userID),
		slog.Int("item_id", itemID),
		slog.Bool("title_changed", input.Title != nil),
		slog.Bool("description_changed", input.Description != nil),
		slog.Bool("done_changed", input.Done != nil),
	)

	err := s.itemRepo.Update(ctx, userID, itemID, input)
	if err != nil {
		s.logger.Error("failed to update item",
			slog.Int("user_id", userID),
			slog.Int("item_id", itemID),
			sl.Err(err),
		)
		return fmt.Errorf("failed to update item %d: %w", itemID, err)
	}

	s.logger.Info("item updated successfully",
		slog.Int("user_id", userID),
		slog.Int("item_id", itemID),
	)

	return nil
}
