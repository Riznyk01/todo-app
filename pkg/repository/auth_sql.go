package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	todoapp "todo-app"
)

type AuthSql struct {
	log *logrus.Logger
	db  *sqlx.DB
}

func NewAuthSql(log *logrus.Logger, db *sqlx.DB) *AuthSql {
	return &AuthSql{
		log: log,
		db:  db,
	}
}

func (r *AuthSql) CreateUser(user todoapp.User) (int, error) {
	fc := "Repository. CreateUser"
	var id int
	q := fmt.Sprintf("INSERT INTO %s (email, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	row := r.db.QueryRow(q, user.Email, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return 0, err
	}
	return id, nil
}
func (r *AuthSql) ExistsUser(email string) (bool, error) {
	fc := "Repository. ExistsUser"
	var count int
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email=$1", usersTable)
	if err := r.db.Get(&count, q, email); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return false, err
	}
	return count > 0, nil
}
func (r *AuthSql) GetUser(email string) (todoapp.User, error) {
	fc := "Repository. GetUser"
	var user todoapp.User
	q := fmt.Sprintf("SELECT id, password_hash FROM %s WHERE email=$1", usersTable)
	if err := r.db.Get(&user, q, email); err != nil {
		r.log.Errorf("%s: %v", fc, err)
	}
	return user, nil
}
func (r *AuthSql) UpdateRefreshTokenInDB(email, newRefreshToken string) error {
	fc := "Repository. UpdateRefreshTokenInDB"

	q := fmt.Sprintf("UPDATE %s SET refresh_token=:refresh_token WHERE email=:email", usersTable)
	_, err := r.db.NamedExec(q, map[string]interface{}{
		"email":         email,
		"refresh_token": newRefreshToken,
	})

	if err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return err
	}
	return nil
}
func (r *AuthSql) CheckRefreshTokenInDB(refreshTokenString string) (string, error) {
	fc := "Repository. CheckRefreshTokenInDB"
	var mail string
	q := fmt.Sprintf("SELECT email FROM %s WHERE refresh_token=$1", usersTable)
	if err := r.db.Get(&mail, q, refreshTokenString); err != nil {
		r.log.Errorf("%s: %v", fc, err)
		return "", err
	}
	return mail, nil
}
