package repository

import "github.com/jmoiron/sqlx"

type Authorization interface{

}

type Recipe interface{

}

type Repository struct{
	Authorization
	Recipe
}

func NewRepository(db *sqlx.DB) *Repository{
	return &Repository{}
}