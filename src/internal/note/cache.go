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

func (r *RepositoryCache) Create(ctx context.Context, note *Note) error {
	if err := r.repo.Create(ctx, note); err != nil {
		return err
	}

	// Удаляем из кэша все записи о заметках пользователя
	return r.client.Del(ctx, getNotesKey(note.AuthorId)).Err()
}

func (r *RepositoryCache) Fetch(ctx context.Context, noteId uuid.UUID) (*Note, error) {
	// Проверяем существует ли запись в кэше
	exists := true

	encoded, err := r.client.Get(ctx, getNoteKey(noteId)).Bytes()
	if err == redis.Nil {
		exists = false
	} else if err != nil {
		return nil, err
	}

	// Если её нет в кэше, то берем информацию из БД и добавляем в кэш
	if !exists {
		// Получаем пользователя из БД
		note, err := r.repo.Fetch(ctx, noteId)
		if err != nil {
			return nil, err
		}

		// Кодируем структуру в строку
		encoded, err := json.MarshalContext(ctx, note)
		if err != nil {
			return nil, err
		}

		// Кешируем запись
		if err := r.client.Set(ctx, getNoteKey(noteId), string(encoded), DEFAULT_TTL).Err(); err != nil {
			return nil, err
		}

		return note, nil
	}

	// Декорируем строку в структуру
	note := Note{}

	if err := json.Unmarshal(encoded, &note); err != nil {
		return nil, err
	}

	// Обновляем время кэша
	if err := r.client.Expire(ctx, getNoteKey(noteId), DEFAULT_TTL).Err(); err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *RepositoryCache) FetchAll(ctx context.Context, userId uuid.UUID) ([]*Note, error) {
	// Attempt to fetch notes from cache
	encoded, err := r.client.Get(ctx, getNotesKey(userId)).Bytes()

	// If the key doesn't exist in cache
	if err == redis.Nil {
		// Fetch notes from the database
		notes, err := r.repo.FetchAll(ctx, userId)
		if err != nil {
			return nil, err
		}

		// Encode the notes slice to JSON
		encoded, err := json.MarshalContext(ctx, notes)
		if err != nil {
			return nil, err
		}

		// Cache the encoded notes
		if err := r.client.Set(ctx, getNotesKey(userId), string(encoded), DEFAULT_TTL).Err(); err != nil {
			return nil, err
		}

		return notes, nil
	} else if err != nil { // Handle other errors from cache retrieval
		return nil, err
	}

	// Decode the cached value into a slice of notes
	notes := make([]*Note, 0)
	if err := json.Unmarshal(encoded, &notes); err != nil {
		return nil, err
	}

	// Refresh cache expiration time
	if err := r.client.Expire(ctx, getNotesKey(userId), DEFAULT_TTL).Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *RepositoryCache) Delete(ctx context.Context, noteId uuid.UUID) error {
	note, err := r.repo.Fetch(ctx, noteId)
	if err != nil {
		return err
	}

	if err := r.repo.Delete(ctx, noteId); err != nil {
		return err
	}

	// Удаляем из кэша все записи о заметках пользователя и закешированную заметку
	return r.client.Del(ctx, getNotesKey(note.AuthorId), getNoteKey(noteId)).Err()
}

func (r *RepositoryCache) Update(ctx context.Context, note *NoteDTO) error {
	n, err := r.Fetch(ctx, note.Id)
	if err != nil {
		return err
	}

	if err := r.repo.Update(ctx, note); err != nil {
		return err
	}

	// Удаляем из кэша все записи о заметках пользователя и закешированную заметку
	return r.client.Del(ctx, getNotesKey(n.AuthorId), getNoteKey(n.Id)).Err()
}
