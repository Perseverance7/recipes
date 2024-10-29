package service

import (
	"github.com/Eagoker/recipes"
	"github.com/Eagoker/recipes/internal/repository"
)

type Authorization interface{
	CreateUser(user recipes.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error) 
}

type Recipe interface{
	CreateRecipe(recipe recipes.Recipe, ingredients []recipes.Ingredient) (int, error)
	GetAllRecipes() (*[]recipes.SimplifiedRecipe, error)
	GetRecipeById(id int) (recipes.FullRecipe, error)
	SaveRecipeToProfile(userId, recipeId int) error
	GetSavedRecipes(userId int) ([]string,error)
	UpdateRecipe(userID, recipeID int, updatedRecipe recipes.FullRecipe) error
}

type Service struct{
	Authorization
	Recipe
}

func NewService(repo *repository.Repository) *Service{
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Recipe: NewRecipesService(repo.Recipe),
	}
}