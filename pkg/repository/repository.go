package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	todoapp "todo-app"
)

type Authorization interface {
	CreateUser(user todoapp.User) (int, error)
	GetUser(email string) (todoapp.User, error)
	UserExists(email string) (bool, error)
	UpdateRefreshTokenInDB(email, newRefreshToken string) error
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

func NewRepository(log *logrus.Logger, db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSql(
			log,
			db),
	}
}
