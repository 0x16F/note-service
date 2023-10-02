package user

import (
	"context"
	"strings"

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

func (r *cachedRepository) Create(ctx context.Context, user *User) error {
	// Создаем пользователя
	return r.repo.Create(ctx, user)
}

func (r *cachedRepository) Fetch(ctx context.Context, userId uuid.UUID) (*User, error) {
	// Проверяем существует ли пользователь в кэше
	exists := true

	encoded, err := r.client.Get(ctx, getUserKey(userId)).Bytes()
	if err == redis.Nil {
		exists = false
	} else if err != nil {
		return nil, err
	}

	// Если его нет в кэше, то берем информацию из БД и добавляем в кэш
	if !exists {
		// Получаем пользователя из БД
		user, err := r.repo.Fetch(ctx, userId)
		if err != nil {
			return nil, err
		}

		// Кодируем структуру в строку
		encoded, err := json.MarshalContext(ctx, user)
		if err != nil {
			return nil, err
		}

		// Кешируем пользователя
		if err := r.client.Set(ctx, getUserKey(user.Id), string(encoded), DEFAULT_TTL).Err(); err != nil {
			return nil, err
		}

		return user, nil
	}

	// Декорируем строку в структуру
	user := User{}

	if err := json.Unmarshal(encoded, &user); err != nil {
		return nil, err
	}

	// Обновляем время кэша
	if err := r.client.Expire(ctx, getUserKey(userId), DEFAULT_TTL).Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *cachedRepository) FetchLogin(ctx context.Context, login string) (*User, error) {
	// Проверяем существует ли пользователь в кэше
	exists := true

	// Преобразуем логин в нижний регистр
	login = strings.ToLower(login)

	encoded, err := r.client.Get(ctx, getUserLoginKey(login)).Bytes()
	if err == redis.Nil {
		exists = false
	} else if err != nil {
		return nil, err
	}

	// Если его нет в кэше, то берем информацию из БД и добавляем в кэш
	if !exists {
		// Получаем пользователя из БД
		user, err := r.repo.FetchLogin(ctx, login)
		if err != nil {
			return nil, err
		}

		// Кодируем структуру в строку
		encoded, err := json.MarshalContext(ctx, user)
		if err != nil {
			return nil, err
		}

		// Кешируем пользователя
		if err := r.client.Set(ctx, getUserLoginKey(user.Login), string(encoded), DEFAULT_TTL).Err(); err != nil {
			return nil, err
		}

		return user, nil
	}

	// Декорируем строку в структуру
	user := User{}

	if err := json.Unmarshal(encoded, &user); err != nil {
		return nil, err
	}

	// Обновляем время кэша
	if err := r.client.Expire(ctx, getUserLoginKey(user.Login), DEFAULT_TTL).Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *cachedRepository) Update(ctx context.Context, user *User) error {
	// Обновляем данные пользователя в БД
	if err := r.repo.Update(ctx, user); err != nil {
		return err
	}

	// Удаляем данные из кэша
	return r.client.Del(ctx, getUserKey(user.Id)).Err()
}
