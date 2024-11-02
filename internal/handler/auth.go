package handler

import (
	"net/http"

	"github.com/Perceverance7/recipes/internal/models"

	"github.com/gin-gonic/gin"
)

type signUpInput struct{
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignUp
// @Tags auth
// @Description Create account
// @Accept json
// @Produce json
// @Param input body signUpInput true "username and password"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context){
	var input signUpInput

	if err := c.BindJSON(&input); err != nil{
		newErrorResponce(c, http.StatusBadRequest, "invalid input body")
		return
	}

	var userInfo models.User
	userInfo.Username = input.Username
	userInfo.Password = input.Password

	id, err := h.services.Authorization.CreateUser(userInfo)
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

// @Summary SignIn
// @Description Authorization
// @Tags auth
// @Accept json
// @Produce json
// @Param input body signInInput true "username and password"
// @Success 200 {object} string
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
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