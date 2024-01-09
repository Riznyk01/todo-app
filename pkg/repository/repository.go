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
}

type TodoItem interface {
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
	}
}
