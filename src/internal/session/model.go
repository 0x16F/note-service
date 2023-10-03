package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=model.go -destination=mocks/service.go

type Session struct {
	Id           uuid.UUID `json:"id" redis:"id"`
	UserId       uuid.UUID `json:"user_id" redis:"user_id"`
	Role         string    `json:"role" redis:"role"`
	CreatedAt    time.Time `json:"created_at" redis:"created_at"`
	LastActivity time.Time `json:"last_activity" redis:"last_activity"`
}

type Repository interface {
	Create(ctx context.Context, session *Session) error
	Delete(ctx context.Context, sessionId uuid.UUID) error
	Update(ctx context.Context, session *Session) error
	Fetch(ctx context.Context, sessionId uuid.UUID) (*Session, error)
	FetchAll(ctx context.Context, userId uuid.UUID) ([]*Session, error)
}

const MAX_SESSIONS = 5
const SESSION_TTL = time.Hour * 24 * 30

var ErrSessionIsNotExists = errors.New("session is not exists")
