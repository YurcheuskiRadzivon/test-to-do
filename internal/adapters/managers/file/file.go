package file

import "github.com/gofiber/fiber/v2"

type FileManager interface {
	UploadFileByNoteID(ctx *fiber.Ctx) error
	GetFile(ctx *fiber.Ctx) error
}

type FileManage struct {
}
