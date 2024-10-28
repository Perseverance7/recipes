package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	usersTable = "users"
	unitsTable = "units"
	ingredientsTable = "ingredients"
	recipesTable = "recipes"
	recipesIngredientsTable = "recipe_ingredients"
	userIngredientsTable = "user_ingredients"
	savedRecipesTable = "saved_recipes"
)

type Config struct{
	Host string
	Port string
	User string
	Password string
	DbName string
	SslMode string
}

func NewPostgresDb(cfg *Config) (*sqlx.DB, error){
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", 
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SslMode))
	if err != nil{
		return nil, err
	}

	return db, nil
}