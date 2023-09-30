package auth

import (
	"net/http"
	"notes-manager/src/controller/web/headers"
	"notes-manager/src/controller/web/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func (r *Router) IsAuthorized(c *fiber.Ctx) error {
	sessionId, err := uuid.Parse(c.Cookies("session-id"))
	if err != nil {
		return responses.New(http.StatusUnauthorized, "session id is not exists", err.Error())
	}

	session, err := r.repo.Sessions.Fetch(sessionId)
	if err != nil {
		if err == redis.Nil {
			return responses.New(http.StatusUnauthorized, "session is expired or not exists", err.Error())
		}

		logrus.Error(err)
		return responses.System(nil, err.Error())
	}

	session.UpdateActivity()

	if err := r.repo.Sessions.Update(session); err != nil {
		logrus.Error(err)
		return responses.System("failed to update session last activity time", err.Error())
	}

	headers.SetUser(c, session)

	return c.Next()
}
