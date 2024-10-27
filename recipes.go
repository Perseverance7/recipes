package recipes

// Структура для таблицы Units
type Unit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Структура для таблицы Ingredients
type Ingredient struct {
	ID     int `json:"id"`
	Name   string `json:"name"`
	UnitID int `json:"unit_id"` 
}

// Структура для таблицы Recipes
type Recipe struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Instructions string `json:"instructions"`
	UserID      int    `json:"user_id"` 
}

// Структура для таблицы RecipeIngredients (связь между рецептами и ингредиентами)
type RecipeIngredient struct {
	RecipeID    int     `json:"recipe_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity    float64 `json:"quantity"` 
}

// Структура для таблицы UserIngredients (ингредиенты пользователя)
type UserIngredient struct {
	UserID      int     `json:"user_id"`      
	IngredientID int     `json:"ingredient_id"` 
	Quantity    float64 `json:"quantity"`     
}

// Структура для таблицы SavedRecipes (сохраненные рецепты пользователя)
type SavedRecipe struct {
	UserID   int       `json:"user_id"`   
	RecipeID int       `json:"recipe_id"` 
}
