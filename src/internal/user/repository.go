package user

import (
	"context"
	"strings"

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

// Create inserts a new user into the database.
func (r *RepositoryDB) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// Fetch retrieves a user from the database based on their user ID.
func (r *RepositoryDB) Fetch(ctx context.Context, userId uuid.UUID) (*User, error) {
	var u User

	err := r.db.WithContext(ctx).Where("id = $1", userId.String()).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserIsNotExists
		}
		return nil, err
	}

	return &u, nil
}

// FetchLogin retrieves a user from the database based on their login.
func (r *RepositoryDB) FetchLogin(ctx context.Context, login string) (*User, error) {
	var u User

	err := r.db.WithContext(ctx).Where("login = $1", strings.ToLower(login)).First(&u).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserIsNotExists
		}
		return nil, err
	}

	return &u, nil
}

// Update modifies the details of an existing user in the database.
func (r *RepositoryDB) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete removes a user from the database based on their user ID.
func (r *RepositoryDB) Delete(ctx context.Context, userId uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = $1", userId.String()).Delete(nil).Error
}
