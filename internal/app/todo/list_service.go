package todo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
	"github.com/TTekmii/todo-list-app/internal/lib/logger/sl"
)

type TodoListService struct {
	listRepo repo.TodoList
	logger   *slog.Logger
}

func NewTodoListService(listRepo repo.TodoList, logger *slog.Logger) *TodoListService {
	return &TodoListService{
		listRepo: listRepo,
		logger:   logger.With("component", "todo_list_service"),
	}
}

func (s *TodoListService) Create(ctx context.Context, userID int, list model.TodoList) (int, error) {
	if list.Title == "" {
		s.logger.Warn("validation failed: empty title",
			slog.Int("user_id", userID),
		)
		return 0, fmt.Errorf("title is required: %w", ErrInvalidInput)
	}

	s.logger.Debug("creating list", slog.Int("user_id", userID), slog.String("title", list.Title))

	listID, err := s.listRepo.Create(ctx, userID, list)
	if err != nil {
		s.logger.Error("failed to create list",
			slog.Int("user_id", userID),
			sl.Err(err),
		)
		return 0, fmt.Errorf("failed to create list: %w", err)
	}

	s.logger.Info("list created",
		slog.Int("list_id", listID),
	)
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
	s.logger.Info("deleting list", slog.Int("list_id", listID))

	err := s.listRepo.Delete(ctx, userID, listID)
	if err != nil {
		s.logger.Error("failed to delete list",
			slog.Int("list_id", listID),
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("failed to delete list: %d: %w", listID, err)
	}

	s.logger.Info("list deleted",
		slog.Int("list_id", listID),
	)
	return nil
}

func (s *TodoListService) Update(ctx context.Context, userID, listID int, input model.UpdateListInput) error {
	if !input.HasChanges() {
		return ErrNoFieldsToUpdate
	}

	if input.Title != nil && *input.Title == "" {
		return fmt.Errorf("title cannot be empty: %w", ErrInvalidInput)
	}

	s.logger.Debug("updating list",
		slog.Int("user_id", userID),
		slog.Int("list_id", listID),
	)

	err := s.listRepo.Update(ctx, userID, listID, input)
	if err != nil {
		return fmt.Errorf("failed to update list: %d: %w", listID, err)
	}

	return nil
}
