package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"
	portuc "ishari-backend/internal/core/port/usecase"

	"github.com/gofiber/fiber/v2"
)

// Controllers holds all HTTP controllers to be registered.
type Controllers struct {
	Health  *controller.HealthController
	Book    *controller.BookController
	Chapter *controller.ChapterController
	User    *controller.UserController
	Auth    *controller.AuthController
}

// AuthDeps holds auth-related dependencies for route registration
type AuthDeps struct {
	AuthUC portuc.AuthUseCase
}

// RegisterRoutes wires all module routes under a common API grouping.
func RegisterRoutes(app *fiber.App, ctrls Controllers, authDeps *AuthDeps) {
	api := app.Group("/api")

	// Health routes (always public)
	if ctrls.Health != nil {
		RegisterHealthRoutes(api, ctrls.Health)
	}

	// Book and User routes (require authDeps for protected endpoints)
	if authDeps != nil && authDeps.AuthUC != nil {
		if ctrls.Book != nil {
			RegisterBookRoutes(api, ctrls.Book, authDeps.AuthUC)
		}
		if ctrls.Chapter != nil {
			RegisterChapterRoutes(api, ctrls.Chapter, authDeps.AuthUC)
		}
		if ctrls.User != nil {
			RegisterUserRoutes(api, ctrls.User, authDeps.AuthUC)
		}
		if ctrls.Auth != nil {
			RegisterAuthRoutes(api, ctrls.Auth, authDeps.AuthUC)
		}
	}
}
