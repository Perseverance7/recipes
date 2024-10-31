package repository

import (
	"github.com/Perceverance7/recipes/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface{
	CreateUser(user models.User) (int, error)
	GetUserSalt(username string) (string, error)
	GetUser(username, password string) (models.User, error)
}

type Recipe interface{
	CreateRecipe(recipe models.Recipe, ingredients []models.Ingredient) (int, error) 
	GetAllRecipes() (*[]models.SimplifiedRecipe, error)
	GetRecipeById(id int) (models.FullRecipe, error)
	SaveRecipeToProfile(userId, recipeId int) error
	GetSavedRecipes(userId int) ([]string,error)
	UpdateRecipe(userID, recipeID int, updatedRecipe models.FullRecipe) error
	DeleteRecipe(userID, recipeID int) error 
	GetRecipesByIngredients(ingredients []string) (*[]models.SimplifiedRecipe, error)
}

type Repository struct{
	Authorization
	Recipe
}

func NewRepository(db *sqlx.DB) *Repository{
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Recipe: NewRecipesPostgres(db),
	}
}