package service

import "github.com/Eagoker/recipes/internal/repository"

type Authorization interface{

}

type Recipe interface{

}

type Service struct{
	Authorization
	Recipe
}

func NewService(repo *repository.Repository) *Service{
	return &Service{}
}