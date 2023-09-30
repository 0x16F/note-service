package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func getSessionsKey(userId uuid.UUID) string {
	return fmt.Sprintf("ns:sessions:user:%s", userId)
}

func getSessionKey(sessionId uuid.UUID) string {
	return fmt.Sprintf("ns:sessions:%s", sessionId)
}

func New(userId uuid.UUID, role string) *Session {
	return &Session{
		Id:           uuid.New(),
		UserId:       userId,
		Role:         role,
		CreatedAt:    time.Now().UTC(),
		LastActivity: time.Now().UTC(),
	}
}

func (s *Session) UpdateActivity() {
	s.LastActivity = time.Now().UTC()
}

func (s *Session) HSet(p redis.Pipeliner) error {
	sessionIdStr := getSessionKey(s.Id)

	if err := p.HSet(context.Background(), sessionIdStr, "id", s.Id.String()).Err(); err != nil {
		return errors.Join(err, errors.New("failed to set session \"id\""))
	}

	if err := p.HSet(context.Background(), sessionIdStr, "user_id", s.UserId.String()).Err(); err != nil {
		return errors.Join(err, errors.New("failed to set session \"user_id\""))
	}
	if err := p.HSet(context.Background(), sessionIdStr, "role", s.Role).Err(); err != nil {
		return errors.Join(err, errors.New("failed to set session \"role\""))
	}

	if err := p.HSet(context.Background(), sessionIdStr, "created_at", s.CreatedAt.UTC()).Err(); err != nil {
		return errors.Join(err, errors.New("failed to set session \"created_at\""))
	}

	if err := p.HSet(context.Background(), sessionIdStr, "last_activity", s.LastActivity.UTC()).Err(); err != nil {
		return errors.Join(err, errors.New("failed to set session \"last_activity\""))
	}

	return nil
}
