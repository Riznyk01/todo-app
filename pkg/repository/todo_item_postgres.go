package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strings"
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
func (r *TodoItemPostgres) GetItemById(userId, itemId int) (todo_app.TodoItem, error) {
	fc := "Repository. todo_item_postgres. GetItemById"
	var item todo_app.TodoItem
	getItemQuery := fmt.Sprintf("SELECT ti.id,ti.favorite,ti.description, ti.done FROM %s ti INNER JOIN %s li ON li.item_id = ti.id INNER JOIN %s ul ON ul.list_id=li.list_id WHERE ti.id=$1 AND ul.user_id=$2",
		todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Get(&item, getItemQuery, itemId, userId); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return todo_app.TodoItem{}, err
	}
	return item, nil
}
func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	fc := "Repository. todo_item_postgres. Delete"

	delUsersItemQuery := fmt.Sprintf("DELETE FROM %s ti USING %s ul, %s li WHERE ti.id=li.item_id AND li.list_id=ul.list_id AND ul.user_id=$1 AND ti.id=$2",
		todoItemsTable, usersListsTable, listsItemsTable)
	if _, err := r.db.Exec(delUsersItemQuery, userId, itemId); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return err
	}
	return nil
}
func (r *TodoItemPostgres) Update(userId, itemId int, item todo_app.UpdateTodoItem) error {
	fc := "Repository. todo_item_postgres. Update"

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if item.Favorite != nil {
		setValues = append(setValues, fmt.Sprintf("favorite=$%d", argId))
		args = append(args, *item.Favorite)
		argId++
	}
	if item.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *item.Description)
		argId++
	}
	if item.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *item.Done)
		argId++
	}

	setStr := strings.Join(setValues, ", ")

	updateUserItemQuery := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id=li.item_id AND ul.list_id=li.list_id AND ul.user_id=$%d AND ti.id=$%d",
		todoItemsTable, setStr, listsItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userId, itemId)
	if _, err := r.db.Exec(updateUserItemQuery, args...); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return err
	}
	return nil
}
