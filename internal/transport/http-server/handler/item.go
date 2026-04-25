package handler

import (
	"net/http"
	"strconv"

	"github.com/TTekmii/todo-list-app/internal/transport/http-server/dto"
	"github.com/gin-gonic/gin"
)

// createItem godoc
// @Summary      Create a new todo item
// @Description  Create a new item in a specific todo list for the authenticated user
// @Tags         items
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id     path      integer                                      true  "List ID"
// @Param        input  body      object{title=string,description=string}      true  "Item details"
// @Success      201    {object}  map[string]int                               "Item ID"
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists/{id}/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var req dto.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	item := req.ToDomain()

	id, err := h.services.TodoItem.Create(c.Request.Context(), userID, listID, item)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

// getAllItems godoc
// @Summary      Get all items in a list
// @Description  Retrieve all todo items for a specific list for the authenticated user
// @Tags         items
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      integer  true  "List ID"
// @Success      200    {array}  dto.ItemResponse
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists/{id}/items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	listID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	items, err := h.services.TodoItem.GetAll(c.Request.Context(), userID, listID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var response []dto.ItemResponse
	for _, item := range items {
		response = append(response, dto.ItemFromDomain(item))
	}

	c.JSON(http.StatusOK, response)
}

// getItemById godoc
// @Summary      Get item by ID
// @Description  Retrieve a specific todo item by its ID within a list
// @Tags         items
// @Produce      json
// @Security     BearerAuth
// @Param        listId  path      integer  true  "List ID"
// @Param        itemId  path      integer  true  "Item ID"
// @Success      200    {object}  dto.ItemResponse
// @Failure      401    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists/{listId}/items/{itemId} [get]
func (h *Handler) getItemById(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	itemID, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid item id param")
		return
	}

	item, err := h.services.TodoItem.GetById(c.Request.Context(), userID, itemID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, dto.ItemFromDomain(item))
}

// updateItem godoc
// @Summary      Update a todo item
// @Description  Update a specific todo item by its ID within a list
// @Tags         items
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        listId  path      integer                                                      true  "List ID"
// @Param        itemId  path      integer                                                      true  "Item ID"
// @Param        input   body      object{title=string,description=string,done=boolean}         true  "Updated item details"
// @Success      200    {object}  map[string]string                                              "Status ok"
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists/{listId}/items/{itemId} [put]
func (h *Handler) updateItem(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var req dto.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.services.TodoItem.Update(c.Request.Context(), userID, id, req.ToDomain()); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// deleteItem godoc
// @Summary      Delete a todo item
// @Description  Delete a specific todo item by its ID within a list
// @Tags         items
// @Security     BearerAuth
// @Param        listId  path      integer  true  "List ID"
// @Param        itemId  path      integer  true  "Item ID"
// @Success      200    {object}  map[string]string  "Status ok"
// @Failure      401    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /api/lists/{listId}/items/{itemId} [delete]
func (h *Handler) deleteItem(c *gin.Context) {
	userID, err := getUserId(c)
	if err != nil {
		return
	}

	itemID, err := strconv.Atoi(c.Param("itemID"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.services.TodoItem.Delete(c.Request.Context(), userID, itemID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
