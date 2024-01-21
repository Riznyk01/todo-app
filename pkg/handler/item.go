package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	todo_app "todo-app"
)

// @Summary Create a new todo item for a todo list belonging to the user.
// @Description Create a new todo item for a todo list belonging to the authenticated user.
// @Security ApiKeyAuth
// @Tags Items
// @ID create item
// @Accept json
// @Procedure json
// @Param id path int true "ID of the todo list for the item"
// @Param input body todo_app.TodoItem true "Information about the todo item (e.g., favorite, description, done.)"
// @Success 200 {integer} integer 1 "Successfully created a new todo item. Returns the ID of the created item."
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists/{id}/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, "Invalid list id param")
		return
	}
	var input todo_app.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary Get all todo list items for the user by list id.
// @Description Retrieve all todo list items belonging to the authenticated user by list id.
// @Security ApiKeyAuth
// @Tags Items
// @ID get all items
// @Produce json
// @Param id path int true "ID of the todo list for retrieve all the items of the list"
// @Success 200 {object} []todo_app.TodoItem
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/lists/{id}/items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, "Invalid list id param")
		return
	}
	allItems, err := h.services.TodoItem.GetAllItems(userId, listId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, allItems)
}

// @Summary Get todo item for the user by item id.
// @Description Retrieve a todo item belonging to the authenticated user by item id.
// @Security ApiKeyAuth
// @Tags Items
// @ID get todo item
// @Produce json
// @Param id path int true "ID of the todo item retrieve"
// @Success 200 {object} todo_app.TodoItem
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/items/{id} [get]
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, "Invalid item id param")
		return
	}
	item, err := h.services.TodoItem.GetItemById(userId, itemId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// @Summary Update todo item for the user by item id.
// @Description Update a todo item belonging to the authenticated user by item id.
// @Security ApiKeyAuth
// @Tags Items
// @ID update todo item
// @Produce json
// @Param id path int true "ID of the todo item to update"
// @Param input body todo_app.UpdateTodoItem true "UpdateTodoItem object with item favorite, description, done"
// @Success 200 {object} statusResponse
// @Failure 400,401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/items/{id} [put]
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var input todo_app.UpdateTodoItem
	if err := c.BindJSON(&input); err != nil {
		newResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = input.Validate()
	if err != nil {
		newResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.TodoItem.Update(userId, id, input)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "Successfully updated the todo item."})
}

// @Summary Delete users todo item by id.
// @Description Delete а todo item belonging to the authenticated user by todo item id.
// @Security ApiKeyAuth
// @Tags Items
// @ID delete item by id
// @Produce json
// @Param id path int true "ID of the todo item to delete"
// @Success 200 {object} statusResponse
// @Failure 400,401 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/items/{id} [delete]
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c, h.log)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newResponseError(c, http.StatusBadRequest, "Invalid id")
		return
	}
	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newResponseError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "Successfully deleted the todo item."})
}
