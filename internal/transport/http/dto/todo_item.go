package dto

import "github.com/TTekmii/todo-list-app/internal/domain/model"

type CreateItemRequest struct {
	Title       string `json:"title" binding:"required,min=2"`
	Description string `json:"description" binding:"omitempty"`
}

func (r CreateItemRequest) ToDomain() model.TodoItem {
	return model.TodoItem{
		Title:       r.Title,
		Description: r.Description,
		Done:        false,
	}
}

type UpdateItemRequest struct {
	Title       *string `json:"title" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
	Done        *bool   `json:"done" binding:"omitempty"`
}

func (r UpdateItemRequest) ToDomain() model.UpdateItemInput {
	return model.UpdateItemInput{
		Title:       r.Title,
		Description: r.Description,
		Done:        r.Done,
	}
}

type ItemResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func ItemFromDomain(i model.TodoItem) ItemResponse {
	return ItemResponse{
		ID:          i.ID,
		Title:       i.Title,
		Description: i.Description,
		Done:        i.Done,
	}
}
