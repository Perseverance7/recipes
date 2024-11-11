package repository

import (
	"fmt"

	"github.com/Perceverance7/recipes/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash, salt) VALUES ($1, $2, $3) RETURNING id", usersTable)
	row := a.db.QueryRow(query, user.Username, user.Password, user.Salt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (a *AuthPostgres) GetUserSalt(username string) (string, error) {
	var salt string
	query := fmt.Sprintf("SELECT salt FROM %s WHERE username=$1", usersTable)
	row := a.db.QueryRow(query, username)
	if err := row.Scan(&salt); err != nil {
		return "", err
	}

	return salt, nil
}

func (a *AuthPostgres) GetUser(username, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := a.db.Get(&user, query, username, password)

	return user, err
}
