package repository

import (
	"notes-manager/src/internal/note"
	"notes-manager/src/internal/session"
	"notes-manager/src/internal/user"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func New(db *gorm.DB, client *redis.Client) *Repository {
	return &Repository{
		Users:    user.NewRepo(db, client),
		Sessions: session.NewRepo(client),
		Notes:    note.NewRepo(db, client),
	}
}
