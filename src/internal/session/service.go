package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// getSessionsKey generates a unique Redis key for storing multiple sessions of a user.
func getSessionsKey(userId uuid.UUID) string {
	return fmt.Sprint(sessionsUserIdKeyBase, userId)
}

// getSessionKey generates a unique Redis key for a specific session.
func getSessionKey(sessionId uuid.UUID) string {
	return fmt.Sprint(sessionsKeyBase, sessionId)
}

// New initializes a new Session instance.
func New(userId uuid.UUID, role string) *Session {
	t := time.Now().UTC()

	return &Session{
		Id:           uuid.New(),
		UserId:       userId,
		Role:         role,
		CreatedAt:    t,
		LastActivity: t,
	}
}

// UpdateActivity updates the LastActivity timestamp of a session to the current time.
func (s *Session) UpdateActivity() {
	s.LastActivity = time.Now().UTC()
}

// HSet populates the Redis hash set with session details.
func (s *Session) HSet(ctx context.Context, p redis.Pipeliner) error {
	sessionIdStr := getSessionKey(s.Id)

	fields := map[string]interface{}{
		"id":            s.Id.String(),
		"user_id":       s.UserId.String(),
		"role":          s.Role,
		"created_at":    s.CreatedAt.UTC(),
		"last_activity": s.LastActivity.UTC(),
	}

	for field, value := range fields {
		if err := p.HSet(ctx, sessionIdStr, field, value).Err(); err != nil {
			return fmt.Errorf("failed to set session \"%s\": %w", field, err)
		}
	}

	return nil
}
