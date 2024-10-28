package handler

import (
	"github.com/Eagoker/recipes"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateRecipeRequest struct {
	Recipe        recipes.Recipe            `json:"recipe"`
	Ingredients  []recipes.Ingredient      `json:"ingredients"`
}

func (h *Handler) createRecipe(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		return
	}
	var input CreateRecipeRequest
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	input.Recipe.UserID = userId

	recipeID, err := h.services.Recipe.CreateRecipe(input.Recipe, input.Ingredients)
	if err != nil {
		NewErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"recipe_id": recipeID})

}

func (h *Handler) getAllRecipes(c *gin.Context){
	
}

func (h *Handler) getRecipeById(c *gin.Context){
	
}

func (h *Handler) updateRecipe(c *gin.Context){
	
}

func (h *Handler) deleteRecipe(c *gin.Context){
	
}