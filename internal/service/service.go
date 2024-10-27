package service

import (
	"github.com/Eagoker/recipes"
	"github.com/Eagoker/recipes/internal/repository"
)

type Authorization interface{
	CreateUser(user recipes.User) (int, error)
	GenerateToken(username, password string) (string, error)
}

type Recipe interface{

}

type Service struct{
	Authorization
	Recipe
}

func NewService(repo *repository.Repository) *Service{
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}