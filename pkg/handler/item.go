package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	todo_app "todo-app"
)

func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, "Invalid list id param")
		return
	}
	var input todo_app.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, "Invalid list id param")
		return
	}
	items, err := h.services.TodoItem.GetAllItems(userId, listId)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) getItemById(c *gin.Context) {

}

func (h *Handler) updateItem(c *gin.Context) {

}

func (h *Handler) deleteItem(c *gin.Context) {

}
