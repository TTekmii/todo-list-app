package handler

import (
	"log/slog"
	"net/http"

	"github.com/TTekmii/todo-list-app/internal/app/auth"
	"github.com/TTekmii/todo-list-app/internal/app/todo"
	"github.com/gin-gonic/gin"
)

type Service struct {
	Auth     *auth.Service
	TodoList *todo.TodoListService
	TodoItem *todo.TodoItemService
}

type Handler struct {
	services *Service
	logger   *slog.Logger
}

func NewHandler(services *Service) *Handler {
	return &Handler{services: services}
}

func getUserId(c *gin.Context) (int, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return 0, nil
	}

	id, ok := userID.(int)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return 0, nil
	}

	return id, nil
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, map[string]string{
		"message": message,
	})
}

type statusResponse struct {
	Status string `json:"status"`
}
