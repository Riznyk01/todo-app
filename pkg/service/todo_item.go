package service

import (
	"github.com/sirupsen/logrus"
	todo_app "todo-app"
	"todo-app/pkg/repository"
)

type TodoItemService struct {
	log      *logrus.Logger
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(log *logrus.Logger, repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
	return &TodoItemService{
		log:      log,
		repo:     repo,
		listRepo: listRepo,
	}
}
func (s *TodoItemService) Create(userId, listId int, item todo_app.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		//List doesn't exists or doesn't belongs to user
		return 0, err
	}
	return s.repo.Create(listId, item)
}
func (s *TodoItemService) GetAllItems(userId, listId int) ([]todo_app.TodoItem, error) {
	return s.repo.GetAllItems(userId, listId)
}
func (s *TodoItemService) GetItemById(userId, itemId int) (todo_app.TodoItem, error) {
	return s.repo.GetItemById(userId, itemId)
}
func (s *TodoItemService) Delete(userId, itemId int) error {
	_, err := s.repo.GetItemById(userId, itemId)
	if err != nil {
		//Item doesn't exists or doesn't belongs to user
		return err
	}
	return s.repo.Delete(userId, itemId)
}
