package repository

import (
	"github.com/Eagoker/recipes"
	
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RecipesPostgres struct{
	db *sqlx.DB
}

func NewRecipesPostgres(db *sqlx.DB) *RecipesPostgres{
	return &RecipesPostgres{db: db}
}

func (r *RecipesPostgres) GetRecipesByIngredients(ingredients []string) (*[]recipes.SimplifiedRecipe, error) {
	var findedRecipes []recipes.SimplifiedRecipe
	
	query := `
		SELECT r.id, r.name, r.user_id
		FROM recipes AS r
		JOIN recipe_ingredients AS ri ON r.id = ri.recipe_id
		JOIN ingredients AS i ON ri.ingredient_id = i.id
		GROUP BY r.id, r.name, r.user_id
		HAVING COUNT(DISTINCT i.name) = $1
		AND ARRAY_AGG(DISTINCT i.name)::text[] @> $2::text[];
	`

	numIngredients := len(ingredients)

	rows, err := r.db.Query(query, numIngredients, pq.Array(ingredients))
	if err != nil {
		return nil, err
	}
	defer rows.Close()


	for rows.Next() {
		var recipeId int
		var recipeName string
		var userId int

		err := rows.Scan(&recipeId, &recipeName, &userId)
		if err != nil {
			return nil, err 
		}

		findedRecipes = append(findedRecipes, recipes.SimplifiedRecipe{
			ID: recipeId, 
			Name: recipeName,
			UserID: userId,
		})

	}

	return &findedRecipes, nil
}

func (r *RecipesPostgres) CreateRecipe(recipe recipes.Recipe, ingredients []recipes.Ingredient) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				// Обработка ошибки при коммите
			}
		}
	}()

	// 1. Добавляем новый рецепт в таблицу recipes
	var recipeID int
	query := fmt.Sprintf("INSERT INTO %s (name, instructions, user_id) VALUES ($1, $2, $3) RETURNING id", recipesTable)
	err = tx.QueryRow(query, recipe.Name, recipe.Instructions, recipe.UserID).Scan(&recipeID)
	if err != nil {
		return 0, err
	}

	// 2. Проверяем ингредиенты и добавляем их, если они отсутствуют
	for _, ingredient := range ingredients {
		var ingredientID int

		// Проверяем, существует ли 
		query := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", ingredientsTable)
		err = tx.QueryRow(query,ingredient.Name).Scan(&ingredientID)

		if err == sql.ErrNoRows {
			// Если ингредиент не существует, добавляем его
			query := fmt.Sprintf("INSERT INTO %s (name, unit_id) VALUES ($1, $2) RETURNING id", ingredientsTable)
			err = tx.QueryRow(query, ingredient.Name, ingredient.UnitID).Scan(&ingredientID)

			if err != nil {
				return 0, err
			}

		} else if err != nil {
			return 0, err
		}

		// Добавляем ингредиент в recipe_ingredients
		query = fmt.Sprintf("INSERT INTO %s (recipe_id, ingredient_id, quantity) VALUES ($1, $2, $3)", recipesIngredientsTable)
		_, err = tx.Exec(
			query,
			recipeID, ingredientID, ingredient.Quantity,
		)

		if err != nil {
			return 0, err
		}
	}

	return recipeID, nil
}

func (r *RecipesPostgres) GetAllRecipes() (*[]recipes.SimplifiedRecipe, error) {
	var simplifiedRecipes []recipes.SimplifiedRecipe

	// Запрос на выборку всех рецептов с их именами и user ID
	query := fmt.Sprintf(`
	SELECT 
		r.id, r.name, r.user_id
	FROM 
		%s AS r;
	`, recipesTable)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipeId int
		var recipeName string
		var userID int

		err := rows.Scan(&recipeId, &recipeName, &userID)
		if err != nil {
			return nil, err
		}

		simplifiedRecipes = append(simplifiedRecipes, recipes.SimplifiedRecipe{
			ID: recipeId,
			Name:   recipeName,
			UserID: userID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &simplifiedRecipes, nil
}

func (r *RecipesPostgres) GetRecipeById(id int) (recipes.FullRecipe, error) {
	var recipe recipes.FullRecipe

	// Запрос для получения рецепта и его ингредиентов
	query := fmt.Sprintf(`
		SELECT r.id, r.name, r.instructions, r.user_id,
			   i.id, i.name, i.unit_id, ri.quantity
		FROM %s AS r
		LEFT JOIN %s AS ri ON ri.recipe_id = r.id
		LEFT JOIN %s AS i ON i.id = ri.ingredient_id
		WHERE r.id = $1
	`, recipesTable, recipesIngredientsTable, ingredientsTable)

	rows, err := r.db.Query(query, id)
	if err != nil {
		return recipe, err
	}
	defer rows.Close()

	// Инициализация переменных для хранения данных
	var (
		recipeID             int
		recipeName           string
		recipeInstructions    string
		userID               int
		ingredientID         sql.NullInt64
		ingredientName       sql.NullString
		unitID               sql.NullInt64
		quantity             sql.NullFloat64
	)

	// Проходим по всем строкам результата
	for rows.Next() {
		err := rows.Scan(&recipeID, &recipeName, &recipeInstructions, &userID, &ingredientID, &ingredientName, &unitID, &quantity)
		if err != nil {
			return recipe, err
		}

		// Заполняем рецепт (без проверки на ID)
		recipe.Recipe = recipes.Recipe{
			ID:          recipeID,
			Name:        recipeName,
			Instructions: recipeInstructions,
			UserID:      userID,
		}

		// Если есть ингредиент, добавляем его
		if ingredientID.Valid {
			ingredient := recipes.Ingredient{
				ID:       int(ingredientID.Int64),
				Name:     ingredientName.String,
				UnitID:   int(unitID.Int64),
				Quantity: float32(quantity.Float64),
			}
			recipe.Ingredients = append(recipe.Ingredients, ingredient)
		}
	}

	if err := rows.Err(); err != nil {
		return recipe, err
	}

	// Отладка: выводим полученные значения
	fmt.Printf("Recipe: %+v\n", recipe)

	return recipe, nil
}

func (r *RecipesPostgres) SaveRecipeToProfile(userId, recipeId int) error {
	// Запрос для вставки в таблицу saved_recipes
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, recipe_id) 
		VALUES ($1, $2)
		ON CONFLICT (user_id, recipe_id) DO NOTHING
	`, savedRecipesTable)

	_, err := r.db.Exec(query, userId, recipeId)
	return err
}

func (r *RecipesPostgres) GetSavedRecipes(userId int) ([]string,error) {
	query := fmt.Sprintf(`
		SELECT r.name
		FROM %s AS r
		JOIN %s AS sr ON sr.recipe_id = r.id
		WHERE sr.user_id = $1

	`, recipesTable, savedRecipesTable)

	var recipes []string
	err := r.db.Select(&recipes, query, userId)
	if err != nil{
		return []string{""}, err
	}
	return recipes, nil
}

func (r *RecipesPostgres) UpdateRecipe(userID, recipeID int, updatedRecipe recipes.FullRecipe) error {
	// Проверяем, является ли пользователь создателем рецепта
	var creatorID int
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE id = $1", recipesTable)
	err := r.db.QueryRow(query, recipeID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("recipe not found")
		}
		return err
	}

	// Если текущий пользователь не создатель рецепта, возвращаем ошибку
	if creatorID != userID {
		return fmt.Errorf("user is not the owner of the recipe")
	}

	// Начинаем транзакцию для атомарности операций
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Обновляем основную информацию о рецепте
	query = fmt.Sprintf(`UPDATE %s
						 SET name = $1, instructions = $2
						 WHERE id = $3`, recipesTable)

	_, err = tx.Exec(query, updatedRecipe.Recipe.Name, updatedRecipe.Recipe.Instructions, recipeID)
	if err != nil {
		return err
	}

	// Удаляем старые ингредиенты рецепта
	query = fmt.Sprintf(`DELETE FROM %s
						 WHERE recipe_id = $1`, recipesIngredientsTable)

	_, err = tx.Exec(query, recipeID)
	if err != nil {
		return err
	}

	// Вставляем новые ингредиенты с проверкой на существование
	for _, ingredient := range updatedRecipe.Ingredients {
		var ingredientID int

		// Проверяем, существует ли ингредиент
		checkIngredientQuery := fmt.Sprintf("SELECT id FROM %s WHERE name = $1", ingredientsTable)
		err = tx.QueryRow(checkIngredientQuery, ingredient.Name).Scan(&ingredientID)

		if err == sql.ErrNoRows {
			// Если ингредиент не существует, добавляем его
			insertIngredientQuery := fmt.Sprintf("INSERT INTO %s (name, unit_id) VALUES ($1, $2) RETURNING id", ingredientsTable)
			err = tx.QueryRow(insertIngredientQuery, ingredient.Name, ingredient.UnitID).Scan(&ingredientID)
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Добавляем ингредиент в recipe_ingredients
		insertRecipeIngredientQuery := fmt.Sprintf(`
			INSERT INTO %s (recipe_id, ingredient_id, quantity) 
			VALUES ($1, $2, $3)`, recipesIngredientsTable)
		_, err = tx.Exec(insertRecipeIngredientQuery, recipeID, ingredientID, ingredient.Quantity)
		if err != nil {
			return err
		}
	}

	// Завершаем транзакцию
	return tx.Commit()
}

func (r *RecipesPostgres) DeleteRecipe(userID, recipeID int) error {
	var creatorID int
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE id = $1", recipesTable)
	err := r.db.QueryRow(query, recipeID).Scan(&creatorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("recipe not found")
		}
		return err
	}

	// Если текущий пользователь не создатель рецепта, возвращаем ошибку
	if creatorID != userID {
		return fmt.Errorf("user is not the owner of the recipe")
	}

	query = fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = $1
		`, recipesTable)
	
	_, err = r.db.Exec(query, recipeID)
	if err != nil {
		return err
	}

	return nil

}


