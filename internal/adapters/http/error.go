package http

import (
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/gofiber/fiber/v2"
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(response.Error{Message: msg})
}
