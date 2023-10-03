package session

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// На пользователя создается хэш мапа со всеми его сессиями
// [hash] ns:sessions:users:<user_id>

// При создании сессии её id добавляется в список сессий пользователя и после этого создается ключ и на него навешивается TTL
// [hash] ns:sessions:<session_id> id <session_id> user_id <user_id> ...

// При таком подходе мы получаем мгновенное чтение данных и их перезапись.

type repository struct {
	client *redis.Client
}

func NewRepo(client *redis.Client) Repository {
	return &repository{
		client: client,
	}
}

func (s *repository) Update(ctx context.Context, session *Session) error {
	_, err := s.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		// Обновляем поля сессии
		if err := session.HSet(p); err != nil {
			return errors.Join(err, errors.New("failed to update session fields"))
		}

		// Обновляем время истечения сессии
		if err := p.Expire(ctx, getSessionKey(session.Id), SESSION_TTL).Err(); err != nil {
			return errors.Join(err, errors.New("failed to update session expire time"))
		}

		// Обновляем время истечения списка с сессиями
		if err := p.Expire(ctx, getSessionsKey(session.UserId), SESSION_TTL).Err(); err != nil {
			return errors.Join(err, errors.New("failed to update sessions list expire time"))
		}

		return nil
	})

	return err
}

func (s *repository) Create(ctx context.Context, session *Session) error {
	// Проверяем кол-во активных сессий пользователя
	count, err := s.client.LLen(ctx, getSessionsKey(session.UserId)).Result()
	if err == redis.Nil {
		count = 0
	} else if err != nil {
		return errors.Join(err, errors.New("failed to get count of user sessions"))
	}

	// Если количество активных сессий равняется или превышает максимальное значение,
	// то удаляем сессию, которая не использовалась дольше всего
	for count >= MAX_SESSIONS {
		// Получаем наименее активную сессию
		oldestSessionId, err := s.FetchOldest(ctx, session.UserId)
		if err != nil {
			return errors.Join(err, errors.New("failed to get oldest session"))
		}

		// Удаляем её из списка сессий пользователя и её саму
		_, err = s.client.Pipelined(ctx, func(p redis.Pipeliner) error {
			if err := p.LRem(ctx, getSessionsKey(session.UserId), 1, oldestSessionId.String()).Err(); err != nil {
				return errors.Join(err, errors.New("failed to delete sessions from user session hashmap"))
			}

			if err := p.Del(ctx, getSessionKey(oldestSessionId)).Err(); err != nil {
				return errors.Join(err, errors.New("failed to delete session"))
			}

			return nil
		})

		if err != nil {
			return err
		}

		count--
	}

	// Создаем сессию и добавляем её идентификатор в список сессий пользователя
	_, err = s.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		if err := session.HSet(p); err != nil {
			return err
		}

		if err := p.Expire(ctx, getSessionKey(session.Id), SESSION_TTL).Err(); err != nil {
			return errors.Join(err, errors.New("failed to set expire on session"))
		}

		if err := p.LPush(ctx, getSessionsKey(session.UserId), session.Id.String()).Err(); err != nil {
			return errors.Join(err, errors.New("failed to add session is user sessions list"))
		}

		if err := p.Expire(ctx, getSessionsKey(session.UserId), SESSION_TTL).Err(); err != nil {
			return errors.Join(err, errors.New("failed to set expire on session"))
		}

		return nil
	})

	return err
}

func (s *repository) Delete(ctx context.Context, sessionId uuid.UUID) error {
	_, err := s.client.Pipelined(ctx, func(p redis.Pipeliner) error {
		session, err := s.Fetch(ctx, sessionId)
		if err != nil {
			return ErrSessionIsNotExists
		}

		if err := p.LRem(ctx, getSessionsKey(session.UserId), 1, session.Id.String()).Err(); err != nil {
			return err
		}

		return p.Del(ctx, getSessionKey(session.Id)).Err()
	})

	return err
}

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

func (s *repository) FetchOldest(ctx context.Context, userId uuid.UUID) (uuid.UUID, error) {
	sessionsIds, err := s.client.LRange(ctx, getSessionsKey(userId), 0, MAX_SESSIONS-1).Result()
	if err != nil {
		return uuid.Nil, err
	}

	maxLastActivity := time.Now().UTC()
	var oldestSessionId uuid.UUID

	for _, sessionId := range sessionsIds {
		id, err := uuid.Parse(sessionId)
		if err != nil {
			return uuid.Nil, err
		}

		lastActivity, err := s.client.HGet(ctx, getSessionKey(id), "last_activity").Time()
		if err == redis.Nil {
			if err := s.client.Del(ctx, getSessionKey(id)).Err(); err != nil {
				return uuid.Nil, err
			}
		} else if err != nil {
			return uuid.Nil, err
		}

		if lastActivity.Unix() < maxLastActivity.Unix() {
			maxLastActivity = lastActivity
			oldestSessionId = id
		}
	}

	return oldestSessionId, nil
}

func (s *repository) FetchAll(ctx context.Context, userId uuid.UUID) ([]*Session, error) {
	sessionsIds, err := s.client.LRange(ctx, getSessionsKey(userId), 0, MAX_SESSIONS-1).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]*Session, 0)

	for _, sessionId := range sessionsIds {
		id, err := uuid.Parse(sessionId)
		if err != nil {
			return nil, err
		}

		session, err := s.Fetch(ctx, id)
		if err == redis.Nil {
			if err := s.client.LRem(ctx, getSessionsKey(userId), 1, sessionId).Err(); err != nil {
				return nil, err
			}

			continue
		} else if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}
