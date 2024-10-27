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
	api := router.Group("/api")
	{
		recipes := api.Group("/recipes")
		{
			recipes.POST("/", h.createRecipe)
			recipes.GET("/", h.getAllRecipes)
			recipes.GET("/:id", h.getRecipeById)
			recipes.PUT("/:id", h.updateRecipe)
			recipes.DELETE("/:id", h.deleteRecipe)
		}
	}

	return router
}