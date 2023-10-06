package auth

import (
	"net/http"
	"notes-manager/src/controller/web/headers"
	"notes-manager/src/controller/web/responses"
	"notes-manager/src/internal/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (r *Router) IsAuthorized(c *fiber.Ctx) error {
	sessionId, err := uuid.Parse(c.Cookies("session-id"))
	if err != nil {
		return responses.New(http.StatusUnauthorized, "session id is not exists", err.Error())
	}

	s, err := r.repo.Sessions.Fetch(c.Context(), sessionId)
	if err != nil {
		if err == session.ErrSessionIsNotExists {
			return responses.New(http.StatusUnauthorized, "session is expired or not exists", err.Error())
		}

		logrus.Error(err)
		return responses.System("failed to fetch session", err.Error())
	}

	s.UpdateActivity()

	if err := r.repo.Sessions.Update(c.Context(), s); err != nil {
		logrus.Error(err)
		return responses.System("failed to update session last activity time", err.Error())
	}

	headers.SetUser(c, s)

	return c.Next()
}
