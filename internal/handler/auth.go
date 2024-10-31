package handler

import (
	"net/http"

	"github.com/Perceverance7/recipes/internal/models"

	"github.com/gin-gonic/gin"
)
func (h *Handler) signUp(c *gin.Context){
	var input models.User

	if err := c.BindJSON(&input); err != nil{
		newErrorResponce(c, http.StatusBadRequest, "invalid input body")
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil{
		newErrorResponce(c, http.StatusConflict, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type signInInput struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context){
	var input signInInput
	if err := c.BindJSON(&input); err != nil{
		newErrorResponce(c, http.StatusBadRequest, "invalid input body")
		return
	}

	token, err := h.services.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponce(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}