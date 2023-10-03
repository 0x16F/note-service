package note

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=model.go -destination=mocks/service.go

const DEFAULT_TTL = time.Hour

type Note struct {
	Id        uuid.UUID `json:"id" gorm:"column:id" example:"6b30e5df-5add-42e1-be60-62b6f98afab1"`
	AuthorId  uuid.UUID `json:"author_id" gorm:"column:author_id" example:"07f3c5a1-70ea-4e3f-b9b5-110d29891673"`
	Title     string    `json:"title" gorm:"column:title" example:"some title"`
	Content   string    `json:"content" gorm:"column:content" example:"some content"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at" example:"2023-10-03T14:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at" example:"2023-10-04T04:20:00Z"`
}

type NoteDTO struct {
	Id        uuid.UUID `json:"id" gorm:"column:id;primaryKey"`
	Title     string    `json:"title" gorm:"column:title"`
	Content   string    `json:"content" gorm:"column:content"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type Repository interface {
	Create(ctx context.Context, note *Note) error
	Fetch(ctx context.Context, noteId uuid.UUID) (*Note, error)
	FetchAll(ctx context.Context, userId uuid.UUID) ([]*Note, error)
	Delete(ctx context.Context, noteId uuid.UUID) error
	Update(ctx context.Context, note *NoteDTO) error
}

var ErrNoteIsNotExists = errors.New("note is not exists")
