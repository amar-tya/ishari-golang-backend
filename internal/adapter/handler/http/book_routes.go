package http

import (
	"ishari-backend/internal/adapter/handler/http/controller"

	"github.com/gofiber/fiber/v2"
)

func RegisterBookRoutes(router fiber.Router, h *controller.BookController) {
	router.Get("/books", h.ListBooks)
	router.Post("/books", h.CreateBook)
	router.Put("/books/:id", h.EditBook)
	router.Delete("/books/:id", h.DeleteBook)
	router.Get("/books/:id", h.GetBookById)
}
