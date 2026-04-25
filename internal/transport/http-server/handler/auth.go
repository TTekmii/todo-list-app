package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type signUnInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

// signUp godoc
// @Summary      Register a new user
// @Description  Create a new user account with username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      object{username=string,password=string,name=string}  true  "User credentials"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input signUnInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Auth.Register(c.Request.Context(), input.Name, input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

// signIn godoc
// @Summary      Login user
// @Description  Get JWT token by username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      object{username=string,password=string}  true  "User credentials"
// @Success      200    {object}  map[string]string  "JWT Token"
// @Failure      400    {object}  map[string]string
// @Failure      401    {object}  map[string]string
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.ShouldBindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Auth.Login(c.Request.Context(), input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenResponse{
		Token: token,
	})
}
