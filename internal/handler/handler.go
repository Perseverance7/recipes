package handler

import (
	"github.com/Eagoker/recipes/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct{
	services *service.Service
}

func NewHandler(services *service.Service) *Handler{
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine{
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)

	}
	api := router.Group("/api", h.userIdentity)
	{
		recipes := api.Group("/recipes")
		{
			recipes.POST("/create", h.createRecipe)
			recipes.GET("/saved", h.getSavedRecipes)
			recipes.GET("/by-ingredients", h.getRecipesByingredients)
			recipes.GET("/", h.getAllRecipes)
			recipes.POST("/:id", h.SaveRecipeToProfile)
			recipes.GET("/:id", h.getRecipeById)
			recipes.PUT("/:id", h.updateRecipe)
			recipes.DELETE("/:id", h.deleteRecipe)
		}
	}

	return router
}

