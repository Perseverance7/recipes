package service

import (
	"github.com/Perceverance7/recipes/internal/models"
	"github.com/Perceverance7/recipes/internal/repository"

	"strings"
	"errors"
	"regexp"
)

type RecipesService struct{
	repo repository.Recipe
}

func NewRecipesService(repo repository.Recipe) *RecipesService{
	return &RecipesService{repo: repo}
}

func (s *RecipesService) CreateRecipe(recipe models.Recipe, ingredients []models.Ingredient) (int, error) {
	// Проверка, чтобы поля Name и Instructions не были пустыми
	if strings.TrimSpace(recipe.Name) == "" || strings.TrimSpace(recipe.Instructions) == "" {
		return 0, errors.New("recipe name and instructions cannot be empty")
	}

	// Проверка, чтобы был хотя бы один ингредиент
	if len(ingredients) < 1 {
		return 0, errors.New("recipe must have at least one ingredient")
	}

	// Приведение имени рецепта и инструкции к нижнему регистру
	recipe.Name = strings.ToLower(recipe.Name)
	recipe.Instructions = strings.ToLower(recipe.Instructions)

	// Приведение имени каждого ингредиента к нижнему регистру
	for i := range ingredients {
		ingredients[i].Name = strings.ToLower(ingredients[i].Name)
	}

	return s.repo.CreateRecipe(recipe, ingredients)
}


func (s *RecipesService) GetAllRecipes() (*[]models.SimplifiedRecipe, error) {
	return s.repo.GetAllRecipes()
}

func (s *RecipesService) GetRecipeById(id int) (*models.FullRecipe, error) {
	recipe, err := s.repo.GetRecipeById(id)
	if err != nil {
		return &models.FullRecipe{}, err
	}

	// Проверка, если рецепт пустой
	if recipe.Recipe.ID == 0 {
		return nil, nil // Возвращаем nil вместо пустой структуры
	}

	return &recipe, nil
}

func (s *RecipesService) SaveRecipeToProfile(userId, recipeId int) error {
	return s.repo.SaveRecipeToProfile(userId, recipeId)
}

func (s *RecipesService) GetSavedRecipes(userId int) ([]models.SavedRecipes, error) {
	return s.repo.GetSavedRecipes(userId)
}

func (s *RecipesService) UpdateRecipe(userID, recipeID int, updatedRecipe models.FullRecipe) error {
	if len(updatedRecipe.Ingredients) < 1{
		return errors.New("recipe must have at least one ingredient")
	}
	return s.repo.UpdateRecipe(userID, recipeID, updatedRecipe)
}

func (s *RecipesService) DeleteRecipe(userID, recipeID int) error {
	return s.repo.DeleteRecipe(userID, recipeID)
}

func (s *RecipesService) GetRecipesByIngredients(ingredients string) (*[]models.SimplifiedRecipe, error) {
	ingredientsArr := extractWords(ingredients)
	return s.repo.GetRecipesByIngredients(ingredientsArr)

}

func extractWords(input string) []string {
	// Разбиваем строку на части по запятой
	ingredients := strings.Split(input, ",")

	var words []string
	re := regexp.MustCompile(`[a-zA-Zа-яА-ЯёЁ]+`)

	for _, ingredient := range ingredients {
		// Очищаем каждый ингредиент от лишних пробелов и символов
		ingredient = strings.TrimSpace(ingredient)
		// Находим все слова, подходящие под регулярное выражение
		if re.MatchString(ingredient) {
			// Преобразуем ингредиент в нижний регистр
			words = append(words, strings.ToLower(ingredient))
		}
	}

	return words
}
