package dto

import "github.com/TTekmii/todo-list-app/internal/domain/model"

type CreateListRequest struct {
	Title       string `json:"title" binding:"required,min=3"`
	Description string `json:"description" binding:"omitempty"`
}

func (r CreateListRequest) ToDomain() model.TodoList {
	return model.TodoList{
		Title:       r.Title,
		Description: r.Description,
	}
}

type UpdateListRequest struct {
	Title       *string `json:"title" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
}

func (r UpdateListRequest) ToDomain() model.UpdateListInput {
	return model.UpdateListInput{
		Title:       r.Title,
		Description: r.Description,
	}
}

type ListResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func ListFromDomain(l model.TodoList) ListResponse {
	return ListResponse{
		ID:          l.ID,
		Title:       l.Title,
		Description: l.Description,
	}
}
