package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	todoapp "todo-app"
	"todo-app/pkg/service"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todoapp.User
	err := c.BindJSON(&input)
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		newResponceError(c, h.log, http.StatusUnprocessableEntity, "Invalid email address.")
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
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

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	err := c.BindJSON(&input)
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, err.Error())
		return
	}
	exist, err := h.services.Authorization.ExistsUser(input.Username)
	if err != nil {
		newResponceError(c, h.log, http.StatusUnauthorized, err.Error())
		return
	}
	if !exist {
		newResponceError(c, h.log, http.StatusUnauthorized, "Invalid credentials. The specified username doesn't exist.")
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			newResponceError(c, h.log, http.StatusUnauthorized, "Invalid user password.")
			return
		} else {
			newResponceError(c, h.log, http.StatusInternalServerError, "Error creating signed token.")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
