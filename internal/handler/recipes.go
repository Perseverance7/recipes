package handler

import (
	"github.com/Perceverance7/recipes/internal/models"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary createRecipe
// @Description creating recipe
// @Security BearerAuth
// @Tags api/recipes
// @Accept json
// @Produce json
// @Param input body models.FullRecipe true "json рецепта"
// @Success 201 {object} string
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/recipes/create [post]
func (h *Handler) createRecipe(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	var input models.FullRecipe
	if err := c.BindJSON(&input); err != nil {
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	input.Recipe.UserID = userId

	recipeID, err := h.services.Recipe.CreateRecipe(input.Recipe, input.Ingredients)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"recipe_id": recipeID})

}

// @Summary getAllRecipes
// @Description getting all recipes
// @Security BearerAuth
// @Tags api/recipes
// @Accept json
// @Produce json
// @Success 200 {array} models.SimplifiedRecipe `json:"data"`
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/recipes [get]
func (h *Handler) getAllRecipes(c *gin.Context){
	fullRecipes, err := h.services.Recipe.GetAllRecipes()
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, "Не удалось получить рецепты")
		return
	}

	// Возвращаем JSON ответ
	c.JSON(http.StatusOK, fullRecipes)
}

// @Summary saveRecipeToProfile
// @Description saving recipe to profile
// @Security BearerAuth
// @Tags api/recipes
// @Accept json
// @Produce json
// @Param id path int true "ID рецепта"
// @Success 200 {object} string
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/recipes/{id} [post]
func (h *Handler) SaveRecipeToProfile(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	recipeId, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		newErrorResponce(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.Recipe.SaveRecipeToProfile(userId, recipeId); err != nil{
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "recipe saved to profile successfully"})
}

// @Summary getRecipeById
// @Description getting recipe by id
// @Security BearerAuth
// @Tags api/recipes
// @Accept json
// @Produce json
// @Success 200 {object} models.FullRecipe `json:"data"`
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/recipes/{id} [get]
func (h *Handler) getRecipeById(c *gin.Context){

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		newErrorResponce(c, http.StatusBadRequest, "invalid id param")
		return
	}

	recipe, err := h.services.Recipe.GetRecipeById(id)
	if err != nil{
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	if recipe == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "recipe not found"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// @Summary getSavedRecipes
// @Description getting saved recipes
// @Security BearerAuth
// @Tags api/recipes
// @Accept json
// @Produce json
// @Success 200 {array} models.SavedRecipes `json:"data"`
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/recipes/saved [get]
func (h *Handler) getSavedRecipes(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	recipes, err := h.services.Recipe.GetSavedRecipes(userId)
	if err != nil{
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, recipes)
}

// @Summary updateRecipe
// @Description updating recipe(only creator can update recipe he created)
// @Security BearerAuth
// @Tags api/recipes
// @Accept json
// @Produce json
// @Param id path int true "ID рецепта"
// @Success 200 {object} string
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/recipes/{id} [put]
func (h *Handler) updateRecipe(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	recipeId, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		newErrorResponce(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.FullRecipe
	if err := c.BindJSON(&input); err != nil{
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Recipe.UpdateRecipe(userId, recipeId, input); err != nil{
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "recipe updated succesfully"})

}

// @Summary deleteRecipe
// @Description deleting recipe(only creator can delete recipe he created)
// @Security BearerAuth
// @Tags api/recipes
// @Accept json
// @Produce json
// @Param id path int true "ID рецепта"
// @Success 200 {object} string
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/recipes/{id} [delete]
func (h *Handler) deleteRecipe(c *gin.Context){
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	recipeId, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		newErrorResponce(c, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.services.Recipe.DeleteRecipe(userId, recipeId); err != nil{
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "recipe delete successful"})
}

// @Summary getRecipeByIngredients
// @Description getting recipe by ingredients
// @Security BearerAuth
// @Tags api/recipes
// @Accept json
// @Produce json
// @Param input body map[string]string true "Список ингредиентов"
// @Success 200 {array} models.SimplifiedRecipe `json:"data"`
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/recipes/by-ingredients [post]
func (h *Handler) getRecipesByingredients(c *gin.Context){
	var input map[string]string
	if err := c.BindJSON(&input); err != nil{
		newErrorResponce(c, http.StatusBadRequest, err.Error())
		return
	}

	ingredients, ok := input["ingredients"]
	if !ok {
		newErrorResponce(c, http.StatusBadRequest, "ingredient is required")
		return
	}

	recipes, err := h.services.Recipe.GetRecipesByIngredients(ingredients); 
	if err != nil{
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, recipes)
}