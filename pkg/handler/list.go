package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	todo_app "todo-app"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	var input todo_app.TodoList
	if err := c.BindJSON(&input); err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type AllListsResponse struct {
	Data []todo_app.TodoList `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, AllListsResponse{
		Data: lists,
	})
}

func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, "invalid id")
		return
	}
	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, "invalid id")
		return
	}
	var input todo_app.TodoList
	if err := c.BindJSON(&input); err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.TodoList.Update(userId, id, input)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponceError(c, h.log, http.StatusBadRequest, "Invalid id")
		return
	}
	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		newResponceError(c, h.log, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "ok"})
}
