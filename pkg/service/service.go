package service

import (
	"github.com/sirupsen/logrus"
	todoapp "todo-app"
	"todo-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user todoapp.User) (int, error)
	GenerateTokenPair(email string) (string, string, error)
	ExistsUser(email string) (bool, error)
	ParseToken(token string) (int, error)
	CheckTokenInDB(refreshTokenString string) (string, error)
	CheckUserPassword(email, password string) error
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(log *logrus.Logger, repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(
			log,
			repos.Authorization),
	}
}
