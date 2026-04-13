package todo

import (
	"github.com/TTekmii/todo-list-app/internal/domain/model"
	"github.com/TTekmii/todo-list-app/internal/domain/repo"
)

type TodoListService struct {
	listRepo repo.TodoList
}

func NewTodoListService(listRepo repo.TodoList) *TodoListService {
	return &TodoListService{listRepo: listRepo}
}

func (s *TodoListService) Create(userId int, list model.TodoList) (int, error) {
	return s.listRepo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]model.TodoList, error) {
	return s.listRepo.GetAll(userId)
}

func (s *TodoListService) GetById(userId, listId int) (model.TodoList, error) {
	return s.listRepo.GetById(userId, listId)
}

func (s *TodoListService) Delete(userId, listId int) error {
	return s.listRepo.Delete(userId, listId)
}

func (s *TodoListService) Update(userId, listId int, input model.UpdateListInput) error {
	if !input.HasChanges() {
		return model.ErrNoFieldsToUpdate
	}

	return s.listRepo.Update(userId, listId, input)
}
