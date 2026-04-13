package todo

import (
	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
)

type TodoItemService struct {
	itemRepo repo.TodoItem
	listRepo repo.TodoList
}

func NewTodoItemService(itemRepo repo.TodoItem, listRepo repo.TodoList) *TodoItemService {
	return &TodoItemService{itemRepo: itemRepo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId, listId int, item model.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.itemRepo.Create(listId, item)
}

func (s *TodoItemService) GetAll(userId, listId int) ([]model.TodoItem, error) {
	return s.itemRepo.GetAll(userId, listId)
}

func (s *TodoItemService) GetById(userId, itemId int) (model.TodoItem, error) {
	return s.itemRepo.GetById(userId, itemId)
}

func (s *TodoItemService) Delete(userId, itemId int) error {
	return s.itemRepo.Delete(userId, itemId)
}

func (s *TodoItemService) Update(userId, itemId int, input model.UpdateItemInput) error {
	return s.itemRepo.Update(userId, itemId, input)
}
