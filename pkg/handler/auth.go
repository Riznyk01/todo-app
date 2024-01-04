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
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

//"required,email,max=64"`
//"required,min=8,max=64"`

// TODO: Implement authentication for both username and email
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	err := c.BindJSON(&input)
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, err.Error())
		return
	}
	exist, err := h.services.Authorization.ExistsUser(input.Email)
	if err != nil {
		newResponceError(c, h.log, http.StatusUnauthorized, err.Error())
		return
	}
	if !exist {
		newResponceError(c, h.log, http.StatusUnauthorized, "Invalid credentials. The specified username doesn't exist.")
		return
	}
	accessToken, refreshToken, err := h.services.Authorization.GenerateTokenPair(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			newResponceError(c, h.log, http.StatusUnauthorized, "Invalid user password.")
			return
		} else {
			newResponceError(c, h.log, http.StatusInternalServerError, "Error creating signed tokens.")
			return
		}
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accesstoken":  accessToken,
		"refreshtoken": refreshToken,
	})
}
