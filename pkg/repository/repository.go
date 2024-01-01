package repository

import (
	"github.com/jmoiron/sqlx"
	todoapp "todo-app"
)

type Authorization interface {
	CreateUser(user todoapp.User) (int, error)
	GetUser(username string) (todoapp.User, error)
	UserExists(username string) (bool, error)
}

type TodoList interface {
}

type TodoItem interface {
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSql(db),
	}
}
