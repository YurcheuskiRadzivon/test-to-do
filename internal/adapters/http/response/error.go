package response

import (
	"github.com/gofiber/fiber/v2"
)

const (
	ErrInvalidRequest  = "INVALID_REQUEST"
	ErrCodeInvalid     = "INVALID_CODE"
	ErrPassHashInvalid = "PASSWORD_HASH_INVALID"
	ErrSignInFailed    = "SIGN_IN_FAILED"
	ErrUnknown         = "UKNOWN_ERROR"
	ErrJWT             = "ENCRYPTION_ERROR"
)

type Error struct {
	Message string `json:"message" example:"message"`
}

func ErrorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(Error{Message: msg})
}
