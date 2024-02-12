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
	Create(userId int, list todoapp.TodoList) (int, error)
	GetAll(userId int) ([]todoapp.TodoList, error)
	GetById(userId, listId int) (todoapp.TodoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, list todoapp.UpdateTodoList) error
}

type TodoItem interface {
	Create(userId, listId int, item todoapp.TodoItem) (int, error)
	GetAllItems(userId, listId int) ([]todoapp.TodoItem, error)
	GetItemById(userId, itemId int) (todoapp.TodoItem, error)
	Delete(userId, itemId int) error
	Update(userId, itemId int, item todoapp.UpdateTodoItem) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(log *logrus.Logger, config *todoapp.Config, repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(log, config, repos.Authorization),
		TodoList:      NewTodoListService(log, repos.TodoList),
		TodoItem:      NewTodoItemService(log, repos.TodoItem, repos.TodoList),
	}
}
