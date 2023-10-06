package notes

import (
	"notes-manager/src/controller/web/headers"
	"notes-manager/src/internal/note"
	"notes-manager/src/usecase/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Router struct {
	repo      *repository.Repository
	validator *validator.Validate
	headers   headers.Getter
}

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required,max=32" example:"some title"`
	Content string `json:"content" example:"some content"`
}

type UpdateNoteRequest struct {
	NoteId    uuid.UUID `json:"note_id" validate:"required" example:"07f3c5a1-70ea-4e3f-b9b5-110d29891673"`
	Title     string    `json:"title" validate:"required,max=32" example:"new title"`
	Content   string    `json:"content" validate:"required" example:"new content"`
	IsPrivate bool      `json:"is_private" example:"false"`
}

type FetchAllParams struct {
	Public  bool `query:"public"`
	PerPage int  `query:"per_page" validate:"min=1,max=50"`
	Page    int  `query:"page" validate:"min=1"`
}

type UserNotesResponse struct {
	Notes []*note.Note `json:"notes"`
	Pages int64        `json:"pages,omitempty"`
}

type NoteResponse struct {
	Note *note.Note `json:"note"`
}
