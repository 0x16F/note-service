package headers

import "github.com/google/uuid"

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
