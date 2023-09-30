package repository

import (
	"notes-manager/src/internal/note"
	"notes-manager/src/internal/session"
	"notes-manager/src/internal/user"
)

type Repository struct {
	Users    user.Repository
	Sessions session.Repository
	Notes    note.Repository
}
