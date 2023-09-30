package routes

import (
	"notes-manager/src/controller/web/routes/auth"
	"notes-manager/src/controller/web/routes/notes"
	"notes-manager/src/usecase/repository"
)

func New(repo *repository.Repository) *Routes {
	return &Routes{
		Auth:  auth.New(repo),
		Notes: notes.New(repo),
	}
}
