package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	todo_app "todo-app"
)

type TodoListPostgres struct {
	log *logrus.Logger
	db  *sqlx.DB
}

func NewTodoListPostgres(log *logrus.Logger, db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{
		log: log,
		db:  db,
	}
}
func (r *TodoListPostgres) Create(userId int, list todo_app.TodoList) (int, error) {
	fc := "Repository. todo_list_postgres. Create"

	tx, err := r.db.Begin()
	if err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return 0, err
	}
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title) values ($1) RETURNING id", todoListsTable)
	row := r.db.QueryRow(createListQuery, list.Title)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		r.log.Errorf("%s: %v", fc, err)
		return 0, err
	}
	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) values ($1,$2) RETURNING id", usersListsTable)
	_, err = tx.Exec(createUsersListQuery, userId, id)
	if err != nil {
		tx.Rollback()
		r.log.Errorf("%s: %v", fc, err)
		return 0, err
	}
	return id, tx.Commit()
}
func (r *TodoListPostgres) GetAll(userId int) ([]todo_app.TodoList, error) {
	fc := "Repository. todo_list_postgres. GetAll"
	var lists []todo_app.TodoList
	getUserListsQuery := fmt.Sprintf("SELECT tl.id, tl.title FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id=$1",
		todoListsTable, usersListsTable)
	if err := r.db.Select(&lists, getUserListsQuery, userId); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return []todo_app.TodoList{}, err
	}
	return lists, nil
}
func (r *TodoListPostgres) GetById(userId, listId int) (todo_app.TodoList, error) {
	fc := "Repository. todo_list_postgres. GetById"
	var list todo_app.TodoList
	getUserListsQuery := fmt.Sprintf("SELECT tl.id, tl.title FROM %s tl INNER JOIN %s ul ON tl.id = ul.list_id WHERE ul.user_id=$1 AND ul.list_id=$2",
		todoListsTable, usersListsTable)
	if err := r.db.Get(&list, getUserListsQuery, userId, listId); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return todo_app.TodoList{}, err
	}
	return list, nil
}
