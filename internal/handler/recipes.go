package handler

import (
	"github.com/Perceverance7/recipes/internal/models"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) getAllRecipes(c *gin.Context){
	fullRecipes, err := h.services.Recipe.GetAllRecipes()
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, "Не удалось получить рецепты")
		return
	}

	// Возвращаем JSON ответ
	c.JSON(http.StatusOK, fullRecipes)
}

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