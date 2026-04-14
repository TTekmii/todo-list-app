package repo

import (
	"context"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
)

type TodoList interface {
	Create(ctx context.Context, userID int, list model.TodoList) (int, error)
	GetAll(ctx context.Context, userID int) ([]model.TodoList, error)
	GetById(ctx context.Context, userID, listID int) (model.TodoList, error)
	Delete(ctx context.Context, userID, listID int) error
	Update(ctx context.Context, userID, listID int, input model.UpdateListInput) error
}
