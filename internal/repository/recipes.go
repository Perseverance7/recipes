package repository

import (
	"database/sql"
	"fmt"

	"github.com/Eagoker/recipes"
	"github.com/jmoiron/sqlx"
)

type RecipesPostgres struct{
	db *sqlx.DB
}

func NewRecipesPostgres(db *sqlx.DB) *RecipesPostgres{
	return &RecipesPostgres{db: db}
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
		r.name, r.user_id
	FROM 
		%s AS r;
	`, recipesTable)

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipeName string
		var userID int

		err := rows.Scan(&recipeName, &userID)
		if err != nil {
			return nil, err
		}

		simplifiedRecipes = append(simplifiedRecipes, recipes.SimplifiedRecipe{
			Name:   recipeName,
			UserID: userID,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &simplifiedRecipes, nil
}



