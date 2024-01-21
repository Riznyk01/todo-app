package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	todo_app "todo-app"
)

// @Summary Create a new todo list for the user
// @Description Create a new todo list for the authenticated user.
// @Security ApiKeyAuth
// @Tags Lists
// @ID create-list
// @Accept json
// @Procedure json
// @Param input body todo_app.TodoList true "Information about the todo list (e.g., title, description, etc.)"
// @Success 200 {integer} integer 1 "Successfully created a new todo list. Returns the ID of the created list."
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	var input todo_app.TodoList
	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

type AllListsResponse struct {
	Data []todo_app.TodoList `json:"data"`
}

// @Summary Get all todo lists for the user
// @Description Retrieve all todo lists belonging to the authenticated user.
// @Security ApiKeyAuth
// @Tags Lists
// @ID get all lists
// @Produce json
// @Success 200 {object} AllListsResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, AllListsResponse{
		Data: lists,
	})
}

// @Summary Get users todo list by id
// @Description Retrieve todo list belonging to the authenticated user by todo list id.
// @Security ApiKeyAuth
// @Tags Lists
// @ID get list by id
// @Produce json
// @Param id path int true "ID of the todo list to retrieve"
// @Success 200 {object} AllListsResponse
// @Failure 400,401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id} [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}
	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

// @Summary Update users todo list by id
// @Description Updates todo list belonging to the authenticated user by todo list id.
// @Security ApiKeyAuth
// @Tags Lists
// @ID update list by id
// @Accept json
// @Produce json
// @Param input body todo_app.UpdateTodoList true "UpdateTodoList object with list title to update"
// @Param id path int true "ID of the todo list to update"
// @Success 200 {object} statusResponse
// @Failure 400,401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, "invalid id.")
		return
	}
	var input todo_app.UpdateTodoList
	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = input.Validate()
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.TodoList.Update(userId, id, input)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "Successfully updated the todo list."})
}

// @Summary Delete users todo list by id
// @Description Deletes todo list belonging to the authenticated user by todo list id.
// @Security ApiKeyAuth
// @Tags Lists
// @ID delete list by id
// @Produce json
// @Param id path int true "ID of the todo list to delete"
// @Success 200 {object} statusResponse
// @Failure 400,401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, "Invalid id.")
		return
	}
	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "Successfully deleted the todo list."})
}
