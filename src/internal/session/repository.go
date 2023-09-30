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

func (s *repository) Update(session *Session) error {
	_, err := s.client.Pipelined(context.Background(), func(p redis.Pipeliner) error {
		// Обновляем поля сессии
		if err := session.HSet(p); err != nil {
			return errors.Join(err, errors.New("failed to update session fields"))
		}

		// Обновляем время истечения сессии
		if err := p.Expire(context.Background(), getSessionKey(session.Id), SESSION_TTL).Err(); err != nil {
			return errors.Join(err, errors.New("failed to update session expire time"))
		}

		// Обновляем время истечения списка с сессиями
		if err := p.Expire(context.Background(), getSessionsKey(session.UserId), SESSION_TTL).Err(); err != nil {
			return errors.Join(err, errors.New("failed to update sessions list expire time"))
		}

		return nil
	})

	return err
}

func (s *repository) Create(session *Session) error {
	// Проверяем кол-во активных сессий пользователя
	count, err := s.client.LLen(context.Background(), getSessionsKey(session.UserId)).Result()
	if err == redis.Nil {
		count = 0
	} else if err != nil {
		return errors.Join(err, errors.New("failed to get count of user sessions"))
	}

	// Если количество активных сессий равняется или превышает максимальное значение,
	// то удаляем сессию, которая не использовалась дольше всего
	if count >= MAX_SESSIONS {
		// Получаем наименее активную сессию
		oldestSessionId, err := s.FetchOldest(session.UserId)
		if err != nil {
			return errors.Join(err, errors.New("failed to get oldest session"))
		}

		// Удаляем её из списка сессий пользователя и её саму
		_, err = s.client.Pipelined(context.Background(), func(p redis.Pipeliner) error {
			if err := p.LRem(context.Background(), getSessionsKey(session.UserId), 1, oldestSessionId.String()).Err(); err != nil {
				return errors.Join(err, errors.New("failed to delete sessions from user session hashmap"))
			}

			if err := p.Del(context.Background(), getSessionKey(oldestSessionId)).Err(); err != nil {
				return errors.Join(err, errors.New("failed to delete session"))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	// Создаем сессию и добавляем её идентификатор в список сессий пользователя
	_, err = s.client.Pipelined(context.Background(), func(p redis.Pipeliner) error {
		if err := session.HSet(p); err != nil {
			return err
		}

		if err := p.Expire(context.Background(), getSessionKey(session.Id), SESSION_TTL).Err(); err != nil {
			return errors.Join(err, errors.New("failed to set expire on session"))
		}

		if err := p.LPush(context.Background(), getSessionsKey(session.UserId), session.Id.String()).Err(); err != nil {
			return errors.Join(err, errors.New("failed to add session is user sessions list"))
		}

		if err := p.Expire(context.Background(), getSessionsKey(session.UserId), SESSION_TTL).Err(); err != nil {
			return errors.Join(err, errors.New("failed to set expire on session"))
		}

		return nil
	})

	return err
}

func (s *repository) Delete(sessionId uuid.UUID) error {
	_, err := s.client.Pipelined(context.Background(), func(p redis.Pipeliner) error {
		session, err := s.Fetch(sessionId)
		if err != nil {
			return ErrSessionIsNotExists
		}

		if err := p.LRem(context.Background(), getSessionsKey(session.UserId), 1, session.Id.String()).Err(); err != nil {
			return err
		}

		return p.Del(context.TODO(), getSessionKey(session.Id)).Err()
	})

	return err
}

func (s *repository) Fetch(sessionId uuid.UUID) (*Session, error) {
	session := Session{}

	if err := s.client.HGetAll(context.TODO(), getSessionKey(sessionId)).Scan(&session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *repository) FetchOldest(userId uuid.UUID) (uuid.UUID, error) {
	sessionsIds, err := s.client.LRange(context.Background(), getSessionsKey(userId), 0, MAX_SESSIONS-1).Result()
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

		lastActivity, err := s.client.HGet(context.Background(), getSessionKey(id), "last_activity").Time()
		if err == redis.Nil {
			if err := s.client.Del(context.Background(), getSessionKey(id)).Err(); err != nil {
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

func (s *repository) FetchAll(userId uuid.UUID) ([]*Session, error) {
	sessionsIds, err := s.client.LRange(context.Background(), getSessionsKey(userId), 0, MAX_SESSIONS-1).Result()
	if err != nil {
		return nil, err
	}

	sessions := make([]*Session, 0)

	for _, sessionId := range sessionsIds {
		id, err := uuid.Parse(sessionId)
		if err != nil {
			return nil, err
		}

		session, err := s.Fetch(id)
		if err == redis.Nil {
			if err := s.client.LRem(context.Background(), getSessionsKey(userId), 1, sessionId).Err(); err != nil {
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
