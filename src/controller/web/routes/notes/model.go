package notes

import (
	"notes-manager/src/usecase/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Router struct {
	repo      *repository.Repository
	validator *validator.Validate
}

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required,max=32"`
	Content string `json:"content"`
}

type UpdateNoteRequest struct {
	NoteId  uuid.UUID `json:"note_id" validate:"required"`
	Title   string    `json:"title" validate:"required,max=32"`
	Content string    `json:"content" validate:"required"`
}
