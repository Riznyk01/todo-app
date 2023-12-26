package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	todoapp "todo-app"
)

type AuthSql struct {
	db *sqlx.DB
}

func NewAuthSql(db *sqlx.DB) *AuthSql {
	return &AuthSql{db: db}
}

func (r *AuthSql) CreateUser(user todoapp.User) (int, error) {
	var id int
	q := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(q, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
