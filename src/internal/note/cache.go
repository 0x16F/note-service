package note

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RepositoryCache struct {
	repo   Repository
	client *redis.Client
}

func NewRepo(db *gorm.DB, client *redis.Client) Repository {
	return &RepositoryCache{
		repo:   NewDatabaseRepo(db),
		client: client,
	}
}

// invalidateNoteCache deletes the cached note and the user's note list from the cache.
func (r *RepositoryCache) invalidateNoteCache(ctx context.Context, authorId, noteId uuid.UUID) error {
	return r.client.Del(ctx, getNotesKey(authorId), getNoteKey(noteId)).Err()
}

// Create adds a new note to the database and invalidates the cache for the author's notes.
func (r *RepositoryCache) Create(ctx context.Context, note *Note) error {
	if err := r.repo.Create(ctx, note); err != nil {
		return err
	}

	return r.invalidateNoteCache(ctx, note.AuthorId, note.Id)
}

// Fetch retrieves a note by its ID. It first checks the cache, and if not found, fetches from the database and updates the cache.
func (r *RepositoryCache) Fetch(ctx context.Context, noteId uuid.UUID) (*Note, error) {
	encoded, err := r.client.Get(ctx, getNoteKey(noteId)).Bytes()

	// If the key doesn't exist in the cache
	if err == redis.Nil {
		note, err := r.repo.Fetch(ctx, noteId)
		if err != nil {
			return nil, err
		}

		// Cache the note
		encodedNote, err := json.MarshalContext(ctx, note)
		if err != nil {
			return nil, err
		}

		if err := r.client.Set(ctx, getNoteKey(noteId), string(encodedNote), defaultTTL).Err(); err != nil {
			return nil, err
		}

		return note, nil
	} else if err != nil { // Handle other errors from cache retrieval
		return nil, err
	}

	// Decode the cached value
	note := &Note{}
	if err := json.Unmarshal(encoded, note); err != nil {
		return nil, err
	}

	// Refresh cache expiration time
	if err := r.client.Expire(ctx, getNoteKey(noteId), defaultTTL).Err(); err != nil {
		return nil, err
	}

	return note, nil
}

// FetchAll retrieves all notes of a user. It first checks the cache, and if not found, fetches from the database and updates the cache.
func (r *RepositoryCache) FetchAll(ctx context.Context, userId uuid.UUID) ([]*Note, error) {
	encoded, err := r.client.Get(ctx, getNotesKey(userId)).Bytes()

	// If the key doesn't exist in the cache
	if err == redis.Nil {
		notes, err := r.repo.FetchAll(ctx, userId)
		if err != nil {
			return nil, err
		}

		// Cache the notes list
		encodedNotes, err := json.MarshalContext(ctx, notes)
		if err != nil {
			return nil, err
		}

		if err := r.client.Set(ctx, getNotesKey(userId), string(encodedNotes), defaultTTL).Err(); err != nil {
			return nil, err
		}

		return notes, nil
	} else if err != nil { // Handle other errors from cache retrieval
		return nil, err
	}

	// Decode the cached value into notes list
	notes := make([]*Note, 0)
	if err := json.Unmarshal(encoded, &notes); err != nil {
		return nil, err
	}

	// Refresh cache expiration time
	if err := r.client.Expire(ctx, getNotesKey(userId), defaultTTL).Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// Delete removes a note by its ID and invalidates the cache for the note and the author's notes.
func (r *RepositoryCache) Delete(ctx context.Context, noteId uuid.UUID) error {
	note, err := r.Fetch(ctx, noteId)
	if err != nil {
		return err
	}

	if err := r.repo.Delete(ctx, noteId); err != nil {
		return err
	}

	return r.invalidateNoteCache(ctx, note.AuthorId, note.Id)
}

// Update modifies a note's details in the database and invalidates the cache for the note and the author's notes.
func (r *RepositoryCache) Update(ctx context.Context, note *NoteDTO) error {
	n, err := r.Fetch(ctx, note.Id)
	if err != nil {
		return err
	}

	if err := r.repo.Update(ctx, note); err != nil {
		return err
	}

	return r.invalidateNoteCache(ctx, n.AuthorId, n.Id)
}
