package handler

import (
	"net/http"
	"strconv"

	"github.com/TTekmii/todo-list-app/internal/transport/http-server/dto"
	"github.com/gin-gonic/gin"
)

// createList godoc
// @Summary      Create a new todo list
// @Description  Create a new list for the authenticated user
// @Tags         lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input  body      object{title=string,description=string}  true  "List details"
// @Success      200    {object}  map[string]int  "List ID"
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists [post]
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

// getAllLists godoc
// @Summary      Get all user lists
// @Description  Retrieve all todo lists for the authenticated user
// @Tags         lists
// @Produce      json
// @Security     BearerAuth
// @Success      200    {object}  map[string][]dto.ListResponse  "Wrapped list of items"
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists [get]
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

// getListById godoc
// @Summary      Get list by ID
// @Description  Retrieve a specific todo list by its ID for the authenticated user
// @Tags         lists
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      integer  true  "List ID"
// @Success      200    {object}  dto.ListResponse
// @Failure      401    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists/{id} [get]
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

// updateList godoc
// @Summary      Update a todo list
// @Description  Update a specific todo list by its ID for the authenticated user
// @Tags         lists
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      integer  true  "List ID"
// @Param        input  body      object{title=string,description=string}  true  "Updated list details"
// @Success      200    {object}  map[string]string  "Status ok"
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists/{id} [put]
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

// deleteList godoc
// @Summary      Delete a todo list
// @Description  Delete a specific todo list by its ID for the authenticated user
// @Tags         lists
// @Security     BearerAuth
// @Param        id  path      integer  true  "List ID"
// @Success      200    {object}  map[string]string  "Status ok"
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists/{id} [delete]
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
