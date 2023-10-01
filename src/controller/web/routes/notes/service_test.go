package notes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"notes-manager/src/controller/web/responses"
	"notes-manager/src/internal/note"
	"notes-manager/src/usecase/repository"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type fields struct {
	repo      *repository.MockRepository
	validator *validator.Validate
}

func TestRouter_FetchAll(t *testing.T) {
	type Test struct {
		name         string
		route        string
		expectedCode int
		callback     func(f fields)
	}

	tests := []Test{
		{
			name:         "200",
			route:        "/notes",
			expectedCode: http.StatusOK,
			callback: func(f fields) {
				f.repo.Notes.EXPECT().FetchAll(gomock.Any(), gomock.Any()).Return([]*note.Note{}, nil)
			},
		},
		{
			name:         "500",
			route:        "/notes",
			expectedCode: http.StatusInternalServerError,
			callback: func(f fields) {
				f.repo.Notes.EXPECT().FetchAll(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrEmptySlice)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				repo:      repository.NewMock(ctrl),
				validator: validator.New(),
			}

			r := &Router{
				repo:      f.repo.GetRepository(),
				validator: f.validator,
			}

			app := fiber.New(fiber.Config{
				ErrorHandler: responses.CustomErrorHandler(),
			})

			route := New(r.repo)

			app.Get("/notes", route.FetchAll)

			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			req.Header.Set("Content-Type", "application/json; charset=utf8")

			tt.callback(f)

			res, err := app.Test(req, -1)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			require.Equal(t, tt.expectedCode, res.StatusCode, string(body))
		})
	}
}

func TestRouter_Fetch(t *testing.T) {
	type Test struct {
		name         string
		route        string
		expectedCode int
		callback     func(f fields)
	}

	tests := []Test{
		{
			name:         "200",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusOK,
			callback: func(f fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&note.Note{}, nil)
			},
		},
		{
			name:         "500",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusInternalServerError,
			callback: func(f fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrEmptySlice)
			},
		},
		{
			name:         "404",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusNotFound,
			callback: func(f fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, note.ErrNoteIsNotExists)
			},
		},
		{
			name:         "403",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusForbidden,
			callback: func(f fields) {
				n := note.Note{
					Id:       uuid.New(),
					AuthorId: uuid.New(),
				}

				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&n, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				repo:      repository.NewMock(ctrl),
				validator: validator.New(),
			}

			r := &Router{
				repo:      f.repo.GetRepository(),
				validator: f.validator,
			}

			app := fiber.New(fiber.Config{
				ErrorHandler: responses.CustomErrorHandler(),
			})

			route := New(r.repo)

			app.Get("/notes/:note_id", route.Fetch)

			req := httptest.NewRequest(http.MethodGet, tt.route, nil)
			req.Header.Set("Content-Type", "application/json; charset=utf8")

			tt.callback(f)

			res, err := app.Test(req, -1)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			require.Equal(t, tt.expectedCode, res.StatusCode, string(body))
		})
	}
}

func TestRouter_Create(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Router{
				repo:      tt.fields.repo.GetRepository(),
				validator: tt.fields.validator,
			}
			if err := r.Create(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Router.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRouter_Delete(t *testing.T) {
	type Test struct {
		name         string
		route        string
		expectedCode int
		callback     func(f fields)
	}

	tests := []Test{
		{
			name:         "200",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusOK,
			callback: func(f fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil)
				f.repo.Notes.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:         "500",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusInternalServerError,
			callback: func(f fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(gorm.ErrEmptySlice)
			},
		},
		{
			name:         "404",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusNotFound,
			callback: func(f fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, note.ErrNoteIsNotExists)
			},
		},
		{
			name:         "403",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusForbidden,
			callback: func(f fields) {
				n := note.Note{
					Id:       uuid.New(),
					AuthorId: uuid.New(),
				}

				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&n, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := fields{
				repo:      repository.NewMock(ctrl),
				validator: validator.New(),
			}

			r := &Router{
				repo:      f.repo.GetRepository(),
				validator: f.validator,
			}

			app := fiber.New(fiber.Config{
				ErrorHandler: responses.CustomErrorHandler(),
			})

			route := New(r.repo)

			app.Delete("/notes/:note_id", route.Delete)

			req := httptest.NewRequest(http.MethodDelete, tt.route, nil)
			req.Header.Set("Content-Type", "application/json; charset=utf8")

			tt.callback(f)

			res, err := app.Test(req, -1)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				assert.Fail(t, err.Error())
			}

			require.Equal(t, tt.expectedCode, res.StatusCode, string(body))
		})
	}
}
