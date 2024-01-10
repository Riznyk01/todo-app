package service

import (
	"github.com/sirupsen/logrus"
	todo_app "todo-app"
	"todo-app/pkg/repository"
)

type TodoListService struct {
	log  *logrus.Logger
	repo repository.TodoList
}

func NewTodoListService(log *logrus.Logger, repo repository.TodoList) *TodoListService {
	return &TodoListService{
		log:  log,
		repo: repo,
	}
}
func (s *TodoListService) Create(userId int, list todo_app.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}
func (s *TodoListService) GetAll(userId int) ([]todo_app.TodoList, error) {
	return s.repo.GetAll(userId)
}
func (s *TodoListService) GetById(userId, listId int) (todo_app.TodoList, error) {
	return s.repo.GetById(userId, listId)
}
func (s *TodoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}
