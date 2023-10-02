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

func (r *RepositoryDB) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Model(&User{}).Create(user).Error
}

func (r *RepositoryDB) Fetch(ctx context.Context, userId uuid.UUID) (*User, error) {
	u := User{}

	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = $1", userId.String()).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserIsNotExists
		}

		return nil, err
	}

	return &u, nil
}

func (r *RepositoryDB) FetchLogin(ctx context.Context, login string) (*User, error) {
	u := User{}

	if err := r.db.WithContext(ctx).Model(&User{}).Where("login = $1", strings.ToLower(login)).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserIsNotExists
		}

		return nil, err
	}

	return &u, nil
}

func (r *RepositoryDB) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Model(&User{}).Updates(user).Error
}
