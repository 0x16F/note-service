package headers

import (
	"notes-manager/src/internal/session"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func SetUser(c *fiber.Ctx, s *session.Session) {
	c.Response().Header.Set(SESSION_ID_KEY, s.Id.String())
	c.Response().Header.Set(USER_ID_KEY, s.UserId.String())
	c.Response().Header.Set(USER_ROLE_KEY, s.Role)
}

func New() Getter {
	return &getter{}
}

func (g *getter) GetSession(c *fiber.Ctx) *Session {
	sessionIdStr := c.Response().Header.Peek(SESSION_ID_KEY)
	userIdStr := c.Response().Header.Peek(USER_ID_KEY)
	roleStr := c.Response().Header.Peek(USER_ROLE_KEY)

	parsedSessionId, _ := uuid.ParseBytes(sessionIdStr)
	parsedUserId, _ := uuid.ParseBytes(userIdStr)

	return &Session{
		SessionId: parsedSessionId,
		UserId:    parsedUserId,
		Role:      string(roleStr),
	}
}
