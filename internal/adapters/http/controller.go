package http

import (
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/gofiber/fiber/v2"
)

type APIController struct {
	app         *fiber.App
	noteService *service.NoteService
}
