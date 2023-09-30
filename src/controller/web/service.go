package web

import (
	"fmt"
	"notes-manager/src/controller/web/routes"
	"notes-manager/src/usecase/repository"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func New(repo *repository.Repository) Web {
	cfg := fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	}

	w := &web{
		app:    fiber.New(cfg),
		routes: routes.New(repo),
	}

	w.SetupRoutes()

	return w
}

func (w *web) SetupRoutes() {
	v0 := w.app.Group("/v0")
	{
		auth := v0.Group("/auth")
		{
			auth.Post("/login", w.routes.Auth.Login)
			auth.Post("/register", w.routes.Auth.Register)
			auth.Post("/logout", w.routes.Auth.IsAuthorized, w.routes.Auth.Logout)
		}

		notes := v0.Group("/notes", w.routes.Auth.IsAuthorized)
		{
			notes.Get("", w.routes.Notes.FetchAll)
			notes.Get("/:note_id", w.routes.Notes.Fetch)
			notes.Delete("/:note_id", w.routes.Notes.Delete)
			notes.Post("", w.routes.Notes.Create)
			notes.Patch("", w.routes.Notes.Update)
		}
	}
}

func (w *web) Start(port uint16) error {
	return w.app.Listen(fmt.Sprintf(":%d", port))
}
