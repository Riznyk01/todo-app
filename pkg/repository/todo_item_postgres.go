package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	todo_app "todo-app"
)

type TodoItemPostgres struct {
	log *logrus.Logger
	db  *sqlx.DB
}

func NewTodoItemPostgres(log *logrus.Logger, db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{
		log: log,
		db:  db,
	}
}
func (r *TodoItemPostgres) Create(listId int, item todo_app.TodoItem) (int, error) {
	fc := "Repository. todo_item_postgres. Create"

	tx, err := r.db.Begin()
	if err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return 0, err
	}
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (favorite, description, done) values ($1,$2,$3) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Favorite, item.Description, item.Done)
	if err := row.Scan(&itemId); err != nil {
		tx.Rollback()
		r.log.Errorf("%s: %v", fc, err)
		return 0, err
	}
	createListsItemQuery := fmt.Sprintf("INSERT INTO %s (item_id, list_id) values ($1,$2)", listsItemsTable)
	_, err = tx.Exec(createListsItemQuery, itemId, listId)
	if err != nil {
		tx.Rollback()
		r.log.Errorf("%s: %v", fc, err)
		return 0, err
	}
	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAllItems(userId, listId int) ([]todo_app.TodoItem, error) {
	fc := "Repository. todo_item_postgres. GetAllItems"
	var items []todo_app.TodoItem
	getListsItemQuery := fmt.Sprintf("SELECT ti.id,ti.favorite,ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id INNER JOIN %s ul ON ul.list_id=li.list_id WHERE ul.user_id=$1 AND li.list_id=$2",
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Select(&items, getListsItemQuery, userId, listId); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return []todo_app.TodoItem{}, err
	}
	return items, nil
}
func (r *TodoItemPostgres) GetItemById(itemId int) (todo_app.TodoItem, error) {
	fc := "Repository. todo_item_postgres. GetItemById"
	var item todo_app.TodoItem
	getItemQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", todoItemsTable)
	if err := r.db.Get(&item, getItemQuery, itemId); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return todo_app.TodoItem{}, err
	}
	return item, nil
}
