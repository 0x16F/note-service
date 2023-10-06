package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=model.go -destination=mocks/service.go

const (
	userLoginKeyBase = "ns:users:login:"
	userKeyBase      = "ns:users:"
	defaultTTL       = time.Hour
	DefaultRole      = "user"
)

type User struct {
	Id           uuid.UUID `json:"id" gorm:"column:id"`
	Login        string    `json:"login" gorm:"column:login" validate:"required,min=3,max=32,alphanum"`
	Password     string    `json:"password" gorm:"column:password"`
	Salt         string    `json:"salt" gorm:"column:salt"`
	Role         string    `json:"role" gorm:"column:role"`
	RegisteredAt time.Time `json:"registered_at" gorm:"column:registered_at"`
	LastLoginAt  time.Time `json:"last_login_at" gorm:"column:last_login_at"`
}

// Repository defines the interface for user-related operations.
type Repository interface {
	Create(ctx context.Context, user *User) error                // Create adds a new user.
	Fetch(ctx context.Context, userId uuid.UUID) (*User, error)  // Fetch retrieves a user by their ID.
	FetchLogin(ctx context.Context, login string) (*User, error) // FetchLogin retrieves a user by their login.
	Update(ctx context.Context, user *User) error                // Update modifies an existing user.
	Delete(ctx context.Context, userId uuid.UUID) error          // Delete removes a user by their ID.
}

var ErrUserIsNotExists = errors.New("user is not exists")
