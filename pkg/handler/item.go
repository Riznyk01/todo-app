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
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, "Invalid item id param")
		return
	}
	item, err := h.services.TodoItem.GetItemById(userId, itemId)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, "invalid id")
		return
	}
	var input todo_app.UpdateTodoItem
	if err := c.BindJSON(&input); err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, err.Error())
		return
	}
	err = input.Validate()
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.TodoItem.Update(userId, id, input)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, "Invalid id")
		return
	}
	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
