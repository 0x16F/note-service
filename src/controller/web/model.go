package web

import (
	"notes-manager/src/controller/web/routes"

	"github.com/gofiber/fiber/v2"
)

type Web interface {
	Start(port uint16) error
}

type web struct {
	app    *fiber.App
	routes *routes.Routes
}
