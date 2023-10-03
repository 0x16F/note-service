package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"notes-manager/src/controller/web/headers"
	mock_headers "notes-manager/src/controller/web/headers/mocks"
	"notes-manager/src/controller/web/responses"
	"notes-manager/src/internal/session"
	"notes-manager/src/internal/user"
	"notes-manager/src/pkg/jsonmap"
	"notes-manager/src/usecase/repository"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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

func TestRouter_Login(t *testing.T) {
	type Test struct {
		name         string
		expectedCode int
		body         string
		callback     func(f *fields)
	}

	u := user.New("login", "password")

	tests := []Test{
		{
			name:         "200",
			expectedCode: http.StatusOK,
			body: jsonmap.Map{
				"login":    "login",
				"password": "password",
			}.String(),
			callback: func(f *fields) {
				f.repo.Users.EXPECT().FetchLogin(gomock.Any(), gomock.Any()).Return(u, nil)
				f.repo.Sessions.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:         "400",
			expectedCode: http.StatusBadRequest,
			body: jsonmap.Map{
				"login":    "login",
				"password": "12345",
			}.String(),
			callback: func(f *fields) {},
		},
		{
			name:         "401",
			expectedCode: http.StatusUnauthorized,
			body: jsonmap.Map{
				"login":    "login",
				"password": "p4ssw0rd",
			}.String(),
			callback: func(f *fields) {
				f.repo.Users.EXPECT().FetchLogin(gomock.Any(), gomock.Any()).Return(u, nil)
			},
		},
		{
			name:         "404",
			expectedCode: http.StatusNotFound,
			body: jsonmap.Map{
				"login":    "login",
				"password": "p4ssw0rd",
			}.String(),
			callback: func(f *fields) {
				f.repo.Users.EXPECT().FetchLogin(gomock.Any(), gomock.Any()).Return(nil, user.ErrUserIsNotExists)
			},
		},
		{
			name:         "500",
			expectedCode: http.StatusInternalServerError,
			body: jsonmap.Map{
				"login":    "login",
				"password": "password",
			}.String(),
			callback: func(f *fields) {
				f.repo.Users.EXPECT().FetchLogin(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrInvalidData)
			},
		},
		{
			name:         "500",
			expectedCode: http.StatusInternalServerError,
			body: jsonmap.Map{
				"login":    "login",
				"password": "password",
			}.String(),
			callback: func(f *fields) {
				f.repo.Users.EXPECT().FetchLogin(gomock.Any(), gomock.Any()).Return(u, nil)
				f.repo.Sessions.EXPECT().Create(gomock.Any(), gomock.Any()).Return(redis.ErrClosed)
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

			app.Post("/login", r.Login)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.body))
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

func TestRouter_Register(t *testing.T) {
	type Test struct {
		name         string
		expectedCode int
		body         string
		callback     func(f *fields)
	}

	tests := []Test{
		{
			name:         "200",
			expectedCode: http.StatusOK,
			body: jsonmap.Map{
				"login":    "login",
				"password": "password",
			}.String(),
			callback: func(f *fields) {
				f.repo.Users.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:         "400",
			expectedCode: http.StatusBadRequest,
			body: jsonmap.Map{
				"login":    "login",
				"password": "12345",
			}.String(),
			callback: func(f *fields) {},
		},
		{
			name:         "500",
			expectedCode: http.StatusInternalServerError,
			body: jsonmap.Map{
				"login":    "login",
				"password": "password",
			}.String(),
			callback: func(f *fields) {
				f.repo.Users.EXPECT().Create(gomock.Any(), gomock.Any()).Return(gorm.ErrInvalidData)
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

			app.Post("/register", r.Register)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(tt.body))
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

func TestRouter_Logout(t *testing.T) {
	type Test struct {
		name         string
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
			name:         "200",
			expectedCode: http.StatusOK,
			callback: func(f *fields) {
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
				f.repo.Sessions.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			name:         "200",
			expectedCode: http.StatusOK,
			callback: func(f *fields) {
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
				f.repo.Sessions.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(session.ErrSessionIsNotExists)
			},
		},
		{
			name:         "500",
			expectedCode: http.StatusInternalServerError,
			callback: func(f *fields) {
				f.headers.EXPECT().GetSession(gomock.Any()).Return(&hs)
				f.repo.Sessions.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(redis.ErrClosed)
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

			app.Post("/logout", r.Logout)

			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
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
