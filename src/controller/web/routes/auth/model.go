package auth

import (
	"notes-manager/src/controller/web/headers"
	"notes-manager/src/usecase/repository"

	"github.com/go-playground/validator/v10"
)

type Router struct {
	validator *validator.Validate
	repo      *repository.Repository
	headers   headers.Getter
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=32,alphanum" example:"login"`
	Password string `json:"password" validate:"required,min=8,max=64" example:"password"`
}

type RegisterRequest struct {
	Login    string `json:"login" validate:"required,min=3,max=32,alphanum" example:"login"`
	Password string `json:"password" validate:"required,min=8,max=64" example:"password"`
}
