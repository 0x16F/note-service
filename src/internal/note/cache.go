package note

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type cachedRepository struct {
	repo   Repository
	client *redis.Client
}

func NewRepo(db *gorm.DB, client *redis.Client) Repository {
	return &cachedRepository{
		repo:   newDatabaseRepo(db),
		client: client,
	}
}

func (r *cachedRepository) Create(ctx context.Context, note *Note) error {
	if err := r.repo.Create(ctx, note); err != nil {
		return err
	}

	// Удаляем из кэша все записи о заметках пользователя
	return r.client.Del(ctx, getNotesKey(note.AuthorId)).Err()
}

func (r *cachedRepository) Fetch(ctx context.Context, noteId uuid.UUID) (*Note, error) {
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
		encoded, err := json.Marshal(note)
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

func (r *cachedRepository) FetchAll(ctx context.Context, userId uuid.UUID) ([]*Note, error) {
	// Проверяем существует ли запись в кэше
	exists := true

	encoded, err := r.client.Get(ctx, getNotesKey(userId)).Bytes()
	if err == redis.Nil {
		exists = false
	} else if err != nil {
		return nil, err
	}

	// Если её нет в кэше, то берем информацию из БД и добавляем в кэш
	if !exists {
		// Получаем пользователя из БД
		notes, err := r.repo.FetchAll(ctx, userId)
		if err != nil {
			return nil, err
		}

		// Кодируем структуру в строку
		encoded, err := json.Marshal(notes)
		if err != nil {
			return nil, err
		}

		// Кешируем запись
		if err := r.client.Set(ctx, getNotesKey(userId), string(encoded), DEFAULT_TTL).Err(); err != nil {
			return nil, err
		}

		return notes, nil
	}

	// Декорируем строку в слайс структур
	notes := make([]*Note, 0)

	if err := json.Unmarshal(encoded, &notes); err != nil {
		return nil, err
	}

	// Обновляем время кэша
	if err := r.client.Expire(ctx, getNotesKey(userId), DEFAULT_TTL).Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *cachedRepository) Delete(ctx context.Context, noteId uuid.UUID) error {
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

func (r *cachedRepository) Update(ctx context.Context, note *NoteDTO) error {
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
