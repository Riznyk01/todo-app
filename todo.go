package todo_app

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
