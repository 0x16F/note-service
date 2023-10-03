package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=model.go -destination=mocks/service.go

const (
	sessionsUserIdKeyBase = "ns:sessions:user:"
	sessionsKeyBase       = "ns:sessions:"
	MaxSessions           = 5
	SessionTTL            = time.Hour * 24 * 30
)

type Session struct {
	Id           uuid.UUID `json:"id" redis:"id"`
	UserId       uuid.UUID `json:"user_id" redis:"user_id"`
	Role         string    `json:"role" redis:"role"`
	CreatedAt    time.Time `json:"created_at" redis:"created_at"`
	LastActivity time.Time `json:"last_activity" redis:"last_activity"`
}

// Repository defines the interface for session-related operations.
type Repository interface {
	Create(ctx context.Context, session *Session) error                 // Create a new session
	Delete(ctx context.Context, sessionId uuid.UUID) error              // Delete a session by its ID
	Update(ctx context.Context, session *Session) error                 // Update an existing session's details
	Fetch(ctx context.Context, sessionId uuid.UUID) (*Session, error)   // Fetch a session by its ID
	FetchAll(ctx context.Context, userId uuid.UUID) ([]*Session, error) // Fetch all sessions associated with a user
}

var ErrSessionIsNotExists = errors.New("session is not exists")
