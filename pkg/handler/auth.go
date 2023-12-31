package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	"strings"
	todoapp "todo-app"
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
		newResponceError(c, h.log, http.StatusUnauthorized, "Invalid credentials. The specified email doesn't exist.")
		return
	}
	if err := h.services.Authorization.CheckUserPassword(input.Email, input.Password); err != nil {
		newResponceError(c, h.log, http.StatusUnauthorized, "Invalid credentials.")
		return
	}
	accessToken, refreshToken, err := h.services.Authorization.GenerateTokenPair(input.Email)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, "Error creating signed tokens.")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"accesstoken":  accessToken,
		"refreshtoken": refreshToken,
	})
}
func (h *Handler) refreshTokens(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newResponceError(c, h.log, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newResponceError(c, h.log, http.StatusUnauthorized, "invalid header")
		return
	}
	userMail, err := h.services.Authorization.CheckTokenInDB(headerParts[1])
	if err != nil {
		newResponceError(c, h.log, http.StatusUnauthorized, "The provided refresh token doesn't exist. Please sign in again to obtain a new token pair.")
		return
	}
	_, err = h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newResponceError(c, h.log, http.StatusUnauthorized, err.Error())
		return
	}
	accessToken, refreshToken, err := h.services.Authorization.GenerateTokenPair(userMail)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, "Error creating signed tokens.")
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"accesstoken":  accessToken,
		"refreshtoken": refreshToken,
	})
}
