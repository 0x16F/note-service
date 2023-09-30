package repository

import (
	"notes-manager/src/internal/note"
	mock_note "notes-manager/src/internal/note/mocks"
	"notes-manager/src/internal/session"
	mock_session "notes-manager/src/internal/session/mocks"
	"notes-manager/src/internal/user"
	mock_user "notes-manager/src/internal/user/mocks"

	"github.com/golang/mock/gomock"
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

func NewMock(ctrl *gomock.Controller) *MockRepository {
	return &MockRepository{
		Users:    mock_user.NewMockRepository(ctrl),
		Sessions: mock_session.NewMockRepository(ctrl),
		Notes:    mock_note.NewMockRepository(ctrl),
	}
}

func (mr *MockRepository) GetRepository() *Repository {
	return &Repository{
		Users:    mr.Users,
		Sessions: mr.Sessions,
		Notes:    mr.Notes,
	}
}
