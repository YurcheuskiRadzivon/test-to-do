package http

import (
	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/admin"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/auth"
	middleware "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/middleware/auth"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/note"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/user"
	"github.com/gofiber/fiber/v2"
)

type APIController struct {
	app             *fiber.App
	noteController  note.NoteController
	userController  user.UserController
	adminController admin.AdminController
	authController  auth.AuthController
	authMiddleware  middleware.AuthMiddleware
	cfg             *config.Config
}
