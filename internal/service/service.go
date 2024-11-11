package service

import (
	"context"

	"github.com/Perceverance7/recipes/internal/models"
	"github.com/Perceverance7/recipes/internal/repository"
	"github.com/redis/go-redis/v9"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type Recipe interface {
	CreateRecipe(recipe models.Recipe, ingredients []models.Ingredient) (int, error)
	GetAllRecipes(context.Context) (*[]models.SimplifiedRecipe, error)
	GetRecipeById(id int) (*models.FullRecipe, error)
	SaveRecipeToProfile(userId, recipeId int) error
	GetSavedRecipes(userId int) ([]models.SavedRecipes, error)
	UpdateRecipe(userID, recipeID int, updatedRecipe models.FullRecipe) error
	DeleteRecipe(userID, recipeID int) error
	GetRecipesByIngredients(ingredients string) (*[]models.SimplifiedRecipe, error)
	DeleteSavedRecipes(userId int, input []int) error
}

type Service struct {
	Authorization
	Recipe
}

func NewService(repo *repository.Repository, rdb *redis.Client) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		Recipe:        NewRecipesService(repo.Recipe, rdb),
	}
}
