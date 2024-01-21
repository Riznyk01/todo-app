package todo_app

import "errors"

type TodoList struct {
	Id    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title" binding:"required"`
	//Description string `json:"description"`
}
type UsersList struct {
	Id     int
	UserId int
	ListId int
}
type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Favorite    bool   `json:"favorite" db:"favorite"`
	Description string `json:"description" db:"description" binding:"required"`
	Done        bool   `json:"done" db:"done"`
}
type ListsItem struct {
	Id     int
	ListId int
	ItemId int
}
type UpdateTodoList struct {
	Title *string `json:"title"`
}

func (i UpdateTodoList) Validate() error {
	if i.Title == nil {
		return errors.New("update structure has no values")
	}
	return nil
}

type UpdateTodoItem struct {
	Favorite    *bool   `json:"favorite"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i UpdateTodoItem) Validate() error {
	if i.Favorite == nil && i.Description == nil && i.Done == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
