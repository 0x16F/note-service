package auth

import (
	"net/http"
	"notes-manager/src/controller/web/headers"
	"notes-manager/src/controller/web/responses"
	"notes-manager/src/internal/session"
	"notes-manager/src/internal/user"
	"notes-manager/src/usecase/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func New(repo *repository.Repository) *Router {
	return &Router{
		validator: validator.New(),
		repo:      repo,
		headers:   headers.New(),
	}
}

// @Summary      Sign in
// @Description  sign in account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request	body	LoginRequest  true  "login params"
// @Success      200
// @Failure      400  {object}  responses.Error
// @Failure      404  {object}  responses.Error
// @Failure      500  {object}  responses.Error
// @Router       /v0/auth/login [post]
func (r *Router) Login(c *fiber.Ctx) error {
	request := LoginRequest{}
	if err := c.BodyParser(&request); err != nil {
		return responses.BadRequest("bad json", err.Error())
	}

	if err := r.validator.StructCtx(c.Context(), &request); err != nil {
		return responses.BadRequest("failed to validate all fields in struct", err.Error())
	}

	u, err := r.repo.Users.FetchLogin(c.Context(), request.Login)
	if err != nil {
		if err == user.ErrUserIsNotExists {
			return responses.New(http.StatusNotFound, "user is not exists", err.Error())
		}

		logrus.Error(err)
		return responses.System("failed to fetch user by login", err)
	}

	ok, err := u.ValidatePassword(request.Password)
	if err != nil {
		logrus.Error(err)
		return responses.System("failed to validate password", err)
	}

	if !ok {
		return responses.New(http.StatusUnauthorized, "bad password or login", nil)
	}

	s := session.New(u.Id, u.Role)

	if err := r.repo.Sessions.Create(c.Context(), s); err != nil {
		logrus.Error(err)
		return responses.System("failed to create new session", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session-id",
		Value:    s.Id.String(),
		Expires:  s.LastActivity.Add(session.SessionTTL),
		HTTPOnly: true,
	})

	return nil
}

// @Summary      Sign up
// @Description  create an account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request	body	RegisterRequest  true  "register params"
// @Success      200
// @Failure      400  {object}  responses.Error
// @Failure      500  {object}  responses.Error
// @Router       /v0/auth/register [post]
func (r *Router) Register(c *fiber.Ctx) error {
	request := RegisterRequest{}
	if err := c.BodyParser(&request); err != nil {
		return responses.BadRequest("bad json", err.Error())
	}

	if err := r.validator.StructCtx(c.Context(), &request); err != nil {
		return responses.BadRequest("failed to validate all fields in struct", err.Error())
	}

	u, err := user.New(request.Login, request.Password)
	if err != nil {
		logrus.Error(err)
		return responses.System("failed to init new user struct", err.Error())
	}

	if err := r.repo.Users.Create(c.Context(), u); err != nil {
		logrus.Error(err)
		return responses.System("failed to create new user", err.Error())
	}

	return nil
}

// @Summary      Logout
// @Description  logout from account and destroy the session
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      500  {object}  responses.Error
// @Router       /v0/auth/logout [post]
func (r *Router) Logout(c *fiber.Ctx) error {
	s := r.headers.GetSession(c)

	c.ClearCookie("session-id")

	if err := r.repo.Sessions.Delete(c.Context(), s.SessionId); err != nil {
		if err == session.ErrSessionIsNotExists {
			return nil
		}

		logrus.Error(err)
		return responses.System("failed to delete session", err.Error())
	}

	return nil
}
