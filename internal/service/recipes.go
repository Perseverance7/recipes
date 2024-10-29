package service

import (
	"strings"

	"github.com/Eagoker/recipes"
	"github.com/Eagoker/recipes/internal/repository"
)

type RecipesService struct{
	repo repository.Recipe
}

func NewRecipesService(repo repository.Recipe) *RecipesService{
	return &RecipesService{repo: repo}
}

func (s *RecipesService) CreateRecipe(recipe recipes.Recipe, ingredients []recipes.Ingredient) (int, error) {
	recipe.Name = strings.ToLower(recipe.Name)
	recipe.Instructions = strings.ToLower(recipe.Instructions)

	for i := range ingredients {
		ingredients[i].Name = strings.ToLower(ingredients[i].Name)
	}

	return s.repo.CreateRecipe(recipe, ingredients)
}

func (s *RecipesService) GetAllRecipes() (*[]recipes.SimplifiedRecipe, error) {
	return s.repo.GetAllRecipes()
}

func (s *RecipesService) GetRecipeById(id int) (recipes.FullRecipe, error) {
	return s.repo.GetRecipeById(id)
}

func (s *RecipesService) SaveRecipeToProfile(userId, recipeId int) error {
	return s.repo.SaveRecipeToProfile(userId, recipeId)
}

func (s *RecipesService) GetSavedRecipes(userId int) ([]string, error) {
	return s.repo.GetSavedRecipes(userId)
}

func (s *RecipesService) UpdateRecipe(userID, recipeID int, updatedRecipe recipes.FullRecipe) error {
	return s.repo.UpdateRecipe(userID, recipeID, updatedRecipe)
}
