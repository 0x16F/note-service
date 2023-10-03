package user

import (
	"context"
	"strings"

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

// Create adds a new user to the database.
func (r *RepositoryCache) Create(ctx context.Context, user *User) error {
	return r.repo.Create(ctx, user)
}

// Fetch retrieves a user either from the cache or the database.
func (r *RepositoryCache) Fetch(ctx context.Context, userId uuid.UUID) (*User, error) {
	// Attempt to fetch the user from cache
	encoded, err := r.client.Get(ctx, getUserKey(userId)).Bytes()

	// If the key doesn't exist in cache
	if err == redis.Nil {
		// Fetch the user from the database
		user, err := r.repo.Fetch(ctx, userId)
		if err != nil {
			return nil, err
		}

		// Encode the user struct to JSON
		encoded, err := json.MarshalContext(ctx, user)
		if err != nil {
			return nil, err
		}

		// Cache the encoded user
		if err := r.client.Set(ctx, getUserKey(user.Id), string(encoded), DefaultTTL).Err(); err != nil {
			return nil, err
		}

		return user, nil
	} else if err != nil { // Handle other errors from cache retrieval
		return nil, err
	}

	// Decode the cached value into a user struct
	user := User{}
	if err := json.Unmarshal(encoded, &user); err != nil {
		return nil, err
	}

	// Refresh cache expiration time
	if err := r.client.Expire(ctx, getUserKey(userId), DefaultTTL).Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RepositoryCache) FetchLogin(ctx context.Context, login string) (*User, error) {
	// Convert login to lowercase for consistency
	login = strings.ToLower(login)

	// Attempt to fetch user from cache by login
	encoded, err := r.client.Get(ctx, getUserLoginKey(login)).Bytes()

	// If the key doesn't exist in cache
	if err == redis.Nil {
		// Fetch the user from the database by login
		user, err := r.repo.FetchLogin(ctx, login)
		if err != nil {
			return nil, err
		}

		// Encode the user struct to JSON
		encoded, err := json.MarshalContext(ctx, user)
		if err != nil {
			return nil, err
		}

		// Cache the encoded user
		if err := r.client.Set(ctx, getUserLoginKey(user.Login), string(encoded), DefaultTTL).Err(); err != nil {
			return nil, err
		}

		return user, nil
	} else if err != nil { // Handle other errors from cache retrieval
		return nil, err
	}

	// Decode the cached value into a user struct
	user := User{}
	if err := json.Unmarshal(encoded, &user); err != nil {
		return nil, err
	}

	// Refresh cache expiration time
	if err := r.client.Expire(ctx, getUserLoginKey(user.Login), DefaultTTL).Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates the user in the database and invalidates the cache for that user.
func (r *RepositoryCache) Update(ctx context.Context, user *User) error {
	// Fetch the current user data (mainly to get the login for cache invalidation)
	u, err := r.Fetch(ctx, user.Id)
	if err != nil {
		return err
	}

	// Update user data in the database
	if err := r.repo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate the user's cache by ID and login
	return r.client.Del(ctx, getUserKey(u.Id), getUserLoginKey(u.Login)).Err()
}

// Delete removes the user from the database and invalidates the cache for that user.
func (r *RepositoryCache) Delete(ctx context.Context, userId uuid.UUID) error {
	// Fetch the user (mainly to get the login for cache invalidation)
	u, err := r.Fetch(ctx, userId)
	if err != nil {
		return err
	}

	// Delete the user from the database
	if err := r.repo.Delete(ctx, userId); err != nil {
		return err
	}

	// Invalidate the user's cache by ID and login
	return r.client.Del(ctx, getUserKey(userId), getUserLoginKey(u.Login)).Err()
}
