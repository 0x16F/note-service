package note

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func newDatabaseRepo(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, note *Note) error {
	return r.db.WithContext(ctx).Model(&Note{}).Create(note).Error
}

func (r *repository) Fetch(ctx context.Context, noteId uuid.UUID) (*Note, error) {
	note := Note{}

	if err := r.db.WithContext(ctx).Model(&Note{}).Where("id = $1", noteId.String()).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoteIsNotExists
		}

		return nil, err
	}

	return &note, nil
}

func (r *repository) FetchAll(ctx context.Context, userId uuid.UUID) ([]*Note, error) {
	notes := make([]*Note, 0)

	if err := r.db.WithContext(ctx).Model(&Note{}).Where("author_id = $1", userId.String()).Scan(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *repository) Delete(ctx context.Context, noteId uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&Note{}).Where("id = $1", noteId.String()).Delete(nil).Error
}

func (r *repository) Update(ctx context.Context, note *NoteDTO) error {
	return r.db.WithContext(ctx).Save(note).Error
}