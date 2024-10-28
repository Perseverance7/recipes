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
