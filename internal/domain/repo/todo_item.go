package repo

import (
	"context"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
)

type TodoItem interface {
	Create(ctx context.Context, listId int, item model.TodoItem) (int, error)
	GetAll(ctx context.Context, userId, listId int) ([]model.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (model.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	Update(ctx context.Context, userId, itemId int, input model.UpdateItemInput) error
}
