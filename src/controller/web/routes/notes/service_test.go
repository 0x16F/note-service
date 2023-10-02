package notes

import (
	"io"
	"net/http"
	"net/http/httptest"
	"notes-manager/src/controller/web/headers"
	mock_headers "notes-manager/src/controller/web/headers/mocks"
	"notes-manager/src/controller/web/responses"
	"notes-manager/src/internal/note"
	"notes-manager/src/internal/session"
	"notes-manager/src/pkg/jsonmap"
	"notes-manager/src/usecase/repository"
	"strings"
	"testing"
	"time"

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
	headers   *mock_headers.MockGetter
}

func newFields(ctrl *gomock.Controller) *fields {
	return &fields{
		repo:      repository.NewMock(ctrl),
		validator: validator.New(),
		headers:   mock_headers.NewMockGetter(ctrl),
	}
}

func TestRouter_FetchAll(t *testing.T) {
	type Test struct {
		name         string
		route        string
		expectedCode int
		callback     func(f *fields)
	}

	tests := []Test{
		{
			name:         "200",
			route:        "/notes",
			expectedCode: http.StatusOK,
			callback: func(f *fields) {
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&headers.Session{})
				f.repo.Notes.EXPECT().FetchAll(gomock.Any(), gomock.Any()).Return([]*note.Note{}, nil)
			},
		},
		{
			name:         "500",
			route:        "/notes",
			expectedCode: http.StatusInternalServerError,
			callback: func(f *fields) {
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&headers.Session{})
				f.repo.Notes.EXPECT().FetchAll(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrEmptySlice)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := newFields(ctrl)

			r := &Router{
				repo:      f.repo.GetRepository(),
				validator: f.validator,
				headers:   f.headers,
			}

			app := fiber.New(fiber.Config{
				ErrorHandler: responses.CustomErrorHandler(),
			})

			app.Get("/notes", r.FetchAll)

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
		callback     func(f *fields)
	}

	hs := headers.Session{
		SessionId: uuid.New(),
		UserId:    uuid.New(),
		Role:      "user",
	}

	n := note.Note{
		Id:        uuid.MustParse("dbfbcf1a-aaec-4791-bed4-77150532014a"),
		Title:     "new title",
		Content:   "new content",
		AuthorId:  hs.UserId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	tests := []Test{
		{
			name:         "200",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusOK,
			callback: func(f *fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&n, nil)
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
			},
		},
		{
			name:         "500",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusInternalServerError,
			callback: func(f *fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrEmptySlice)
			},
		},
		{
			name:         "404",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusNotFound,
			callback: func(f *fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, note.ErrNoteIsNotExists)
			},
		},
		{
			name:         "403",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusForbidden,
			callback: func(f *fields) {
				n := note.Note{
					Id:       uuid.New(),
					AuthorId: uuid.New(),
				}

				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&n, nil)
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := newFields(ctrl)

			r := &Router{
				repo:      f.repo.GetRepository(),
				validator: f.validator,
				headers:   f.headers,
			}

			app := fiber.New(fiber.Config{
				ErrorHandler: responses.CustomErrorHandler(),
			})

			app.Get("/notes/:note_id", r.Fetch)

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
	type Test struct {
		name         string
		route        string
		body         string
		expectedCode int
		callback     func(f *fields)
	}

	hs := headers.Session{
		SessionId: uuid.New(),
		UserId:    uuid.New(),
		Role:      "user",
	}

	tests := []Test{
		{
			name:  "200",
			route: "/notes",
			body: jsonmap.Map{
				"title":   "some title",
				"content": "some content",
			}.String(),
			expectedCode: http.StatusOK,
			callback: func(f *fields) {
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
				f.repo.Notes.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:  "500",
			route: "/notes",
			body: jsonmap.Map{
				"title":   "some title",
				"content": "some content",
			}.String(),
			expectedCode: http.StatusInternalServerError,
			callback: func(f *fields) {
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
				f.repo.Notes.EXPECT().Create(gomock.Any(), gomock.Any()).Return(gorm.ErrEmptySlice)
			},
		},
		{
			name:  "400",
			route: "/notes",
			body: jsonmap.Map{
				"content": "some content",
			}.String(),
			callback:     func(f *fields) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := newFields(ctrl)

			r := &Router{
				repo:      f.repo.GetRepository(),
				validator: f.validator,
				headers:   f.headers,
			}

			app := fiber.New(fiber.Config{
				ErrorHandler: responses.CustomErrorHandler(),
			})

			app.Post("/notes", r.Create)

			req := httptest.NewRequest(http.MethodPost, tt.route, strings.NewReader(tt.body))
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

func TestRouter_Delete(t *testing.T) {
	type Test struct {
		name         string
		route        string
		expectedCode int
		callback     func(f *fields)
	}

	s := session.New(uuid.New(), "user")

	hs := headers.Session{
		SessionId: s.Id,
		UserId:    s.UserId,
		Role:      s.Role,
	}

	n := note.Note{
		Id:        uuid.MustParse("dbfbcf1a-aaec-4791-bed4-77150532014a"),
		Title:     "new title",
		Content:   "new content",
		AuthorId:  hs.UserId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	tests := []Test{
		{
			name:         "200",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusOK,
			callback: func(f *fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&n, nil)
				f.headers.EXPECT().GetSession(gomock.All()).Return(&hs)
				f.repo.Notes.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:         "500",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusInternalServerError,
			callback: func(f *fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrEmptySlice)
			},
		},
		{
			name:         "404",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusNotFound,
			callback: func(f *fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(nil, note.ErrNoteIsNotExists)
			},
		},
		{
			name:         "403",
			route:        "/notes/dbfbcf1a-aaec-4791-bed4-77150532014a",
			expectedCode: http.StatusForbidden,
			callback: func(f *fields) {
				n := note.Note{
					Id:       uuid.New(),
					AuthorId: uuid.New(),
				}

				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&n, nil)
				f.headers.EXPECT().GetSession(gomock.All()).Return(&hs)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := newFields(ctrl)

			r := &Router{
				repo:      f.repo.GetRepository(),
				validator: f.validator,
				headers:   f.headers,
			}

			app := fiber.New(fiber.Config{
				ErrorHandler: responses.CustomErrorHandler(),
			})

			app.Delete("/notes/:note_id", r.Delete)

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

func TestRouter_Update(t *testing.T) {
	type Test struct {
		name         string
		route        string
		body         string
		expectedCode int
		callback     func(f *fields)
	}

	s := session.New(uuid.New(), "user")

	hs := headers.Session{
		SessionId: s.Id,
		UserId:    s.UserId,
		Role:      s.Role,
	}

	n := note.Note{
		Id:        uuid.MustParse("dbfbcf1a-aaec-4791-bed4-77150532014a"),
		Title:     "new title",
		Content:   "new content",
		AuthorId:  hs.UserId,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	tests := []Test{
		{
			name:  "200",
			route: "/notes",
			body: jsonmap.Map{
				"note_id": "dbfbcf1a-aaec-4791-bed4-77150532014a",
				"title":   "new title",
				"content": "new content",
			}.String(),
			expectedCode: http.StatusOK,
			callback: func(f *fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&n, nil)
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
				f.repo.Notes.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:  "500",
			route: "/notes",
			body: jsonmap.Map{
				"note_id": "dbfbcf1a-aaec-4791-bed4-77150532014a",
				"title":   "new title",
				"content": "new content",
			}.String(),
			expectedCode: http.StatusInternalServerError,
			callback: func(f *fields) {
				f.repo.Notes.EXPECT().Fetch(gomock.Any(), gomock.Any()).Return(&n, nil)
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
				f.repo.Notes.EXPECT().Update(gomock.Any(), gomock.Any()).Return(gorm.ErrEmptySlice)
			},
		},
		{
			name:  "400",
			route: "/notes",
			body: jsonmap.Map{
				"title":   "new title",
				"content": "new content",
			}.String(),
			callback:     func(f *fields) {},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			f := newFields(ctrl)

			r := &Router{
				repo:      f.repo.GetRepository(),
				validator: f.validator,
				headers:   f.headers,
			}

			app := fiber.New(fiber.Config{
				ErrorHandler: responses.CustomErrorHandler(),
			})

			app.Patch("/notes", r.Update)

			req := httptest.NewRequest(http.MethodPatch, tt.route, strings.NewReader(tt.body))
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
