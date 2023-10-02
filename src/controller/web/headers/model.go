package headers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

//go:generate mockgen -source=model.go -destination=mocks/service.go

const (
	SESSION_ID_KEY = "X-Session-Id"
	USER_ID_KEY    = "X-User-Id"
	USER_ROLE_KEY  = "X-User-Role"
)

type Session struct {
	SessionId uuid.UUID `json:"session_id"`
	UserId    uuid.UUID `json:"user_id"`
	Role      string    `json:"role"`
}

type Getter interface {
	GetSession(c *fiber.Ctx) *Session
}

type getter struct{}
