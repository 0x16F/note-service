package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

const DEFAULT_ROLE = "user"
const DEFAULT_TTL = time.Hour

type User struct {
	Id           uuid.UUID `json:"id" gorm:"column:id"`
	Login        string    `json:"login" gorm:"column:login" validate:"required,min=3,max=32,alphanum"`
	Password     string    `json:"password" gorm:"column:password"`
	Salt         string    `json:"salt" gorm:"column:salt"`
	Role         string    `json:"role" gorm:"column:role"`
	RegisteredAt time.Time `json:"registered_at" gorm:"column:registered_at"`
	LastLoginAt  time.Time `json:"last_login_at" gorm:"column:last_login_at"`
}

type Repository interface {
	Create(ctx context.Context, user *User) error
	Fetch(ctx context.Context, userId uuid.UUID) (*User, error)
	FetchLogin(ctx context.Context, login string) (*User, error)
	Update(ctx context.Context, user *User) error
}

var ErrUserIsNotExists = errors.New("user is not exists")
