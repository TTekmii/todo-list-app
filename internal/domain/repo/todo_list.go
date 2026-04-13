package repo

import (
	"context"

	"github.com/TTekmii/todo-list-app/internal/domain/model"
)

type TodoList interface {
	Create(ctx context.Context, userId int, list model.TodoList) (int, error)
	GetAll(ctx context.Context, userId int) ([]model.TodoList, error)
	GetById(ctx context.Context, userId, listId int) (model.TodoList, error)
	Delete(ctx context.Context, userId, listId int) error
	Update(ctx context.Context, userId, listId int, input model.UpdateListInput) error
}
