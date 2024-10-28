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

		// Проверяем, существует ли ингредиент
		err = tx.QueryRow(
			"SELECT id FROM ingredients WHERE name = $1",
			ingredient.Name,
		).Scan(&ingredientID)

		if err == sql.ErrNoRows {
			// Если ингредиент не существует, добавляем его
			err = tx.QueryRow(
				"INSERT INTO ingredients (name, unit_id) VALUES ($1, $2) RETURNING id",
				ingredient.Name, ingredient.UnitID,
			).Scan(&ingredientID)
			if err != nil {
				return 0, err
			}
		} else if err != nil {
			return 0, err
		}

		// Добавляем ингредиент в recipe_ingredients
		_, err = tx.Exec(
			"INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity) VALUES ($1, $2, $3)",
			recipeID, ingredientID, ingredient.Quantity,
		)
		if err != nil {
			return 0, err
		}
	}

	return recipeID, nil
}


