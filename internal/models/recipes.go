package models

// Структура для таблицы Units
type Unit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Структура для таблицы Ingredients
type Ingredient struct {
	ID     int `json:"id"`
	Name   string `json:"name" binding:"required"`
	UnitID int `json:"unit_id" binding:"required"`
	Quantity float32 `json:"quantity" binding:"required"` 
}

// Структура для таблицы Recipes
type Recipe struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required"`
	Instructions string `json:"instructions" binding:"required"`
	UserID      int    `json:"user_id"` 
}

// Структура для таблицы RecipeIngredients (связь между рецептами и ингредиентами)
type RecipeIngredient struct {
	RecipeID    int     `json:"recipe_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity    float64 `json:"quantity"` 
}

// Структура для таблицы SavedRecipes (сохраненные рецепты пользователя)
type SavedRecipe struct {
	UserID   int       `json:"user_id"`   
	RecipeID int       `json:"recipe_id"` 
}

// --------------------------

type SimplifiedRecipe struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

type FullRecipe struct {
	Recipe        Recipe            `json:"recipe"`
	Ingredients  []Ingredient      `json:"ingredients"`
}
