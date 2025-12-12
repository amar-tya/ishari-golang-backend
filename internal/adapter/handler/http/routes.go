package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"

	"github.com/gofiber/fiber/v2"
)

// Controllers holds all HTTP controllers to be registered.
type Controllers struct {
	Health *controller.HealthController
	Book   *controller.BookController
}

// RegisterRoutes wires all module routes under a common API grouping.
func RegisterRoutes(app *fiber.App, ctrls Controllers) {
	v1 := app.Group("/api/v1")
	if ctrls.Health != nil {
		RegisterHealthRoutes(v1, ctrls.Health)
	}
	if ctrls.Book != nil {
		RegisterBookRoutes(v1, ctrls.Book)
	}
}
