package routes

import (
	"notes-manager/src/controller/web/routes/auth"
	"notes-manager/src/controller/web/routes/notes"
)

type Routes struct {
	Auth  *auth.Router
	Notes *notes.Router
}
