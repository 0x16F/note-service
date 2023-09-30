package repository

import (
	"notes-manager/src/internal/note"
	mock_note "notes-manager/src/internal/note/mocks"
	"notes-manager/src/internal/session"
	mock_session "notes-manager/src/internal/session/mocks"
	"notes-manager/src/internal/user"
	mock_user "notes-manager/src/internal/user/mocks"
)

type Repository struct {
	Users    user.Repository
	Sessions session.Repository
	Notes    note.Repository
}

type MockRepository struct {
	Users    *mock_user.MockRepository
	Sessions *mock_session.MockRepository
	Notes    *mock_note.MockRepository
}
