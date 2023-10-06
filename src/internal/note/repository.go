package note

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RepositoryDB struct {
	db *gorm.DB
}

func NewDatabaseRepo(db *gorm.DB) Repository {
	return &RepositoryDB{
		db: db,
	}
}

// Create adds a new note to the database.
func (r *RepositoryDB) Create(ctx context.Context, note *Note) error {
	return r.db.WithContext(ctx).Model(&Note{}).Create(note).Error
}

// Fetch retrieves a note by its ID from the database.
// Returns ErrNoteIsNotExists if the note does not exist.
func (r *RepositoryDB) Fetch(ctx context.Context, noteId uuid.UUID) (*Note, error) {
	var note Note

	if err := r.db.WithContext(ctx).Model(&Note{}).Where("id = $1", noteId.String()).First(&note).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoteIsNotExists
		}

		return nil, err
	}

	return &note, nil
}

// FetchAll retrieves all notes of a specific user from the database.
func (r *RepositoryDB) FetchAll(ctx context.Context, userId uuid.UUID) ([]*Note, error) {
	var notes []*Note

	if err := r.db.WithContext(ctx).Model(&Note{}).Where("author_id = $1", userId.String()).Find(&notes).Error; err != nil {
		return nil, err
	}

	return notes, nil
}

// FetchPublic retrieves all public notes, and count of elements from the database.
func (r *RepositoryDB) FetchPublic(ctx context.Context, page int, limit int) ([]*Note, int64, error) {
	notes := make([]*Note, 0)
	offset := 0

	if page == 1 {
		page = 0
		offset = 0
	} else {
		offset = page * limit
	}

	var count int64 = 0
	if err := r.db.WithContext(ctx).Model(&Note{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Model(&Note{}).Limit(limit).Offset(offset).Scan(&notes).Error; err != nil {
		return nil, 0, err
	}

	return notes, count, nil
}

// Delete removes a note by its ID from the database.
func (r *RepositoryDB) Delete(ctx context.Context, noteId uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&Note{}).Where("id = $1", noteId.String()).Delete(nil).Error
}

// Update modifies an existing note in the database.
func (r *RepositoryDB) Update(ctx context.Context, note *NoteDTO) error {
	return r.db.WithContext(ctx).Save(note).Error
}
