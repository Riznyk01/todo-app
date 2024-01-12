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
	Id int `json:"id" db:"id"`
	//Title       string `json:"title"`
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
		return errors.New("the title is empty, please provide a valid title")
	}
	return nil
}
