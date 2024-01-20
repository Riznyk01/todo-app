package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
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
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newResponseError(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
}
func getUserId(c *gin.Context, log *logrus.Logger) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newResponseError(c, http.StatusInternalServerError, "User id not found")
		return 0, errors.New("User id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		newResponseError(c, http.StatusInternalServerError, "User id not found")
		return 0, errors.New("User id not found")
	}
	return idInt, nil
}
