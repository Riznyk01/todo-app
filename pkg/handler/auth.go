package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
	"strings"
	todoapp "todo-app"
)

// @Summary SignUp
// @Description create account
// @Tags Authentication
// @ID create-account
// @Accept json
// @Procedure json
// @Param input body todoapp.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {

	var input todoapp.User
	err := c.BindJSON(&input)
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		newResponseError(c, http.StatusBadRequest, "Invalid email address.")
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

//"required,email,max=64"`
//"required,min=8,max=64"`

// TODO: Implement authentication for both username and email
// @Summary SignIn
// @Description login
// @Tags Authentication
// @ID login
// @Accept json
// @Procedure json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	err := c.BindJSON(&input)
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	if _, err := mail.ParseAddress(input.Email); err != nil {
		newResponseError(c, http.StatusBadRequest, "Invalid email address.")
		return
	}
	exist, err := h.services.Authorization.ExistsUser(input.Email)
	if err != nil {
		newResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}
	if !exist {
		newResponseError(c, http.StatusUnauthorized, "Account not found.")
		return
	}
	if err := h.services.Authorization.CheckUserPassword(input.Email, input.Password); err != nil {
		newResponseError(c, http.StatusUnauthorized, "Incorrect password. Please try again.")
		return
	}
	accessToken, refreshToken, err := h.services.Authorization.GenerateTokenPair(input.Email)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, "Error creating signed tokens.")
		return
	}
	c.JSON(http.StatusOK, tokenResponse{accessToken, refreshToken})
}

// @Summary Refresh access and refresh tokens
// @Description Refreshes the access and refresh tokens based on the provided refresh token
// @Tags Authentication
// @ID refresh tokens
// @Procedure json
// @Param Authorization header string true "Bearer {refresh_token}"
// @Success 200 {object} tokenResponse
// @Failure 401,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/refresh-tokens [post]
func (h *Handler) refreshTokens(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newResponseError(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newResponseError(c, http.StatusUnauthorized, "invalid header")
		return
	}
	userMail, err := h.services.Authorization.CheckTokenInDB(headerParts[1])
	if err != nil {
		newResponseError(c, http.StatusUnauthorized, "The provided refresh token doesn't exist. Please sign in again to obtain a new token pair.")
		return
	}
	_, err = h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}
	accessToken, refreshToken, err := h.services.Authorization.GenerateTokenPair(userMail)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, "Error creating signed tokens.")
		return
	}
	c.JSON(http.StatusOK, tokenResponse{accessToken, refreshToken})
}
