package handler

import (
	"net/http"
	"strconv"

	"github.com/TTekmii/todo-list-app/internal/transport/http/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	var req dto.CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	id, err := h.services.TodoList.Create(c.Request.Context(), userID, req.ToDomain())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) getAllLists(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(c.Request.Context(), userID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var response []dto.ListResponse
	for _, list := range lists {
		response = append(response, dto.ListFromDomain(list))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *Handler) getListById(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoList.GetById(c.Request.Context(), userID, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.ListFromDomain(list))
}

func (h *Handler) updateList(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var req dto.UpdateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.services.TodoList.Update(c.Request.Context(), userID, id, req.ToDomain()); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

func (h *Handler) deleteList(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoList.Delete(c.Request.Context(), userID, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
