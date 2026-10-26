package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	todo "github.com/slemeshaev/todo/internal/models"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			return 0, ErrDuplicate
		}

		return 0, err
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(username string) (todo.User, error) {
	var user todo.User
	query := fmt.Sprintf("SELECT id, password_hash AS password FROM %s WHERE username=$1", usersTable)
	err := r.db.Get(&user, query, username)

	if errors.Is(err, sql.ErrNoRows) {
		return todo.User{}, ErrNotFound
	}

	return user, err
}
