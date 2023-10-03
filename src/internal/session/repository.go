package session

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type repository struct {
	client *redis.Client
}

func NewRepo(client *redis.Client) Repository {
	return &repository{
		client: client,
	}
}

// Update updates an existing session.
func (s *repository) Update(ctx context.Context, session *Session) error {
	_, err := s.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := session.HSet(ctx, p); err != nil {
			return errors.Wrap(err, "failed to update session fields")
		}

		if err := p.Expire(ctx, getSessionKey(session.Id), SessionTTL).Err(); err != nil {
			return errors.Wrap(err, "failed to update session expire time")
		}

		if err := p.Expire(ctx, getSessionsKey(session.UserId), SessionTTL).Err(); err != nil {
			return errors.Wrap(err, "failed to update sessions list expire time")
		}

		return nil
	})

	return err
}

// Create creates a new session and manages the user's active sessions.
func (s *repository) Create(ctx context.Context, session *Session) error {
	count, err := s.client.LLen(ctx, getSessionsKey(session.UserId)).Result()
	if err != nil && err != redis.Nil {
		return errors.Wrap(err, "failed to get count of user sessions")
	}

	for count >= MaxSessions {
		oldestSessionId, err := s.FetchOldest(ctx, session.UserId)
		if err != nil {
			return errors.Wrap(err, "failed to get oldest session")
		}

		_, err = s.client.Pipelined(ctx, func(p redis.Pipeliner) error {
			if err := p.LRem(ctx, getSessionsKey(session.UserId), 1, oldestSessionId.String()).Err(); err != nil {
				return errors.Wrap(err, "failed to delete session from user session list")
			}

			return p.Del(ctx, getSessionKey(oldestSessionId)).Err()
		})

		if err != nil {
			return err
		}

		count--
	}

	_, err = s.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := session.HSet(ctx, p); err != nil {
			return err
		}

		if err := p.Expire(ctx, getSessionKey(session.Id), SessionTTL).Err(); err != nil {
			return errors.Wrap(err, "failed to set session expiry")
		}

		if err := p.LPush(ctx, getSessionsKey(session.UserId), session.Id.String()).Err(); err != nil {
			return errors.Wrap(err, "failed to add session to user's session list")
		}

		return p.Expire(ctx, getSessionsKey(session.UserId), SessionTTL).Err()
	})

	return err
}

// Delete removes a session.
func (s *repository) Delete(ctx context.Context, sessionId uuid.UUID) error {
	session, err := s.Fetch(ctx, sessionId)
	if err != nil {
		return ErrSessionIsNotExists
	}

	_, err = s.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := p.LRem(ctx, getSessionsKey(session.UserId), 1, session.Id.String()).Err(); err != nil {
			return err
		}

		return p.Del(ctx, getSessionKey(session.Id)).Err()
	})

	return err
}

// Fetch retrieves a session based on its ID.
func (s *repository) Fetch(ctx context.Context, sessionId uuid.UUID) (*Session, error) {
	session := Session{}
	result := s.client.HGetAll(ctx, getSessionKey(sessionId))

	if m, err := result.Result(); err != nil || len(m) == 0 {
		if err == nil {
			return nil, ErrSessionIsNotExists
		}

		return nil, err
	}

	if err := result.Scan(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

// FetchOldest retrieves the oldest session for a user.
func (s *repository) FetchOldest(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	sessionsIds, err := s.client.LRange(ctx, getSessionsKey(userId), 0, MaxSessions-1).Result()
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to fetch session IDs")
	}

	var oldestSessionId uuid.UUID
	maxLastActivity := time.Now().UTC()

	for _, sessionId := range sessionsIds {
		id, err := uuid.Parse(sessionId)
		if err != nil {
			return uuid.Nil, errors.Wrap(err, "failed to parse session ID")
		}

		lastActivity, err := s.client.HGet(ctx, getSessionKey(id), "last_activity").Time()
		if err == redis.Nil {
			if err := s.client.Del(ctx, getSessionKey(id)).Err(); err != nil {
				return uuid.Nil, errors.Wrap(err, "failed to delete session key")
			}
		} else if err != nil {
			return uuid.Nil, errors.Wrap(err, "failed to get last activity for session")
		}

		if lastActivity.Before(maxLastActivity) {
			maxLastActivity = lastActivity
			oldestSessionId = id
		}
	}

	return oldestSessionId, nil
}

// FetchAll retrieves all sessions for a user.
func (s *repository) FetchAll(ctx context.Context, userId uuid.UUID) ([]*Session, error) {
	sessionsIds, err := s.client.LRange(ctx, getSessionsKey(userId), 0, MaxSessions-1).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch session IDs")
	}

	var sessions []*Session

	for _, sessionId := range sessionsIds {
		id, err := uuid.Parse(sessionId)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse session ID")
		}

		session, err := s.Fetch(ctx, id)
		if err == redis.Nil {
			if err := s.client.LRem(ctx, getSessionsKey(userId), 1, sessionId).Err(); err != nil {
				return nil, errors.Wrap(err, "failed to remove session ID from list")
			}
			continue
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to fetch session by ID")
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}
