package response

import (
	"github.com/gofiber/fiber/v2"
)

const (
	ErrInvalidRequest   = "INVALID_REQUEST"
	ErrUnknown          = "UKNOWN_ERROR"
	ErrJWT              = "ENCRYPTION_ERROR"
	ErrInvalidOwnerType = "INVALID_OWNER_TYPE"
	ErrNotEnoughRights  = "NOT_ENOUGH_RIGHTS"
	ErrNotImplemented   = "NOT_IMPLEMENTED"
	ErrIvalidPassword   = "INVALID_PASSWORD"
	ErrInvalidFileID    = "INVALID_FILE_ID"
	ErrInvalidNoteID    = "INVALID_NOTE_ID"
	ErrInvalidToken     = "INVALID_OR_EXPIRED_TOKEN"
)

type Error struct {
	Message string `json:"message" example:"message"`
}

func ErrorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(Error{Message: msg})
}
