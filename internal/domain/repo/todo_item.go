package repo

import (
	"context"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
)

type TodoItem interface {
	Create(ctx context.Context, listID int, item model.TodoItem) (int, error)
	GetAll(ctx context.Context, userID, listID int) ([]model.TodoItem, error)
	GetById(ctx context.Context, userID, itemID int) (model.TodoItem, error)
	Delete(ctx context.Context, userID, itemID int) error
	Update(ctx context.Context, userID, itemID int, input model.UpdateItemInput) error
}
