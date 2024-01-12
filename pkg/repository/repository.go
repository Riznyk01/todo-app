package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	todoapp "todo-app"
)

type Authorization interface {
	CreateUser(user todoapp.User) (int, error)
	GetUser(email string) (todoapp.User, error)
	ExistsUser(email string) (bool, error)
	UpdateRefreshTokenInDB(email, newRefreshToken string) error
	CheckRefreshTokenInDB(refreshTokenString string) (string, error)
}

type TodoList interface {
	Create(userId int, list todoapp.TodoList) (int, error)
	GetAll(userId int) ([]todoapp.TodoList, error)
	GetById(userId, listId int) (todoapp.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, list todoapp.UpdateTodoList) error
}

type TodoItem interface {
	Create(listId int, item todoapp.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]todoapp.TodoItem, error)
	GetItemById(userId, itemId int) (todoapp.TodoItem, error)
	Delete(userId, itemId int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(log *logrus.Logger, db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSql(
			log,
			db),
		TodoList: NewTodoListPostgres(
			log,
			db),
		TodoItem: NewTodoItemPostgres(
			log,
			db),
	}
}
