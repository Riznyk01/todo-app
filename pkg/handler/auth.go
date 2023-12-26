package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	todoapp "todo-app"
)

func (h *Handler) signUp(c *gin.Context) {
	var input todoapp.User
	err := c.BindJSON(&input)
	if err != nil {
		newResponceError(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newResponceError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {

}
