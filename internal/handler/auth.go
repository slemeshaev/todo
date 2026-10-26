package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	todo "github.com/slemeshaev/todo/internal/models"
	"github.com/slemeshaev/todo/internal/service"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body todo.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]

func (h *Handler) signUp(c *gin.Context) {
	var input todo.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		if errors.Is(err, service.ErrUsernameTaken) {
			newErrorResponse(c, http.StatusConflict, err.Error())
			return
		}

		newInternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"id": id,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			newErrorResponse(c, http.StatusUnauthorized, err.Error())
			return
		}

		newInternalErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"token": token,
	})
}
