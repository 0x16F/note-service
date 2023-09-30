package auth

import (
	"notes-manager/src/usecase/repository"

	"github.com/go-playground/validator/v10"
)

type Router struct {
	validator *validator.Validate
	repo      *repository.Repository
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}
