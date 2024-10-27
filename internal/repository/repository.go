package repository

import (
	"github.com/Eagoker/recipes"
	"github.com/jmoiron/sqlx"
)

type Authorization interface{
	CreateUser(user recipes.User) (int, error)
	GetUserSalt(username string) (string, error)
	GetUser(username, password string) (recipes.User, error)
}

type Recipe interface{

}

type Repository struct{
	Authorization
	Recipe
}

func NewRepository(db *sqlx.DB) *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}