package file

import (
	"context"
	"net/http"
	"strconv"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/gofiber/fiber/v2"
)

type FileMetaService interface {
	GetFileMetaURI(ctx context.Context, id int) (string, error)
	DeleteFileMetaByID(ctx context.Context, id int) error
}

type FileManager interface {
	DownloadFile(ctx *fiber.Ctx, path string) error
	DeleteFile(ctx *fiber.Ctx, path string) error
}

type FileController interface {
	DownloadFile(ctx *fiber.Ctx) error
	//UploadFile(ctx *fiber.Ctx) error
	DeleteFile(ctx *fiber.Ctx) error
}

type FileControl struct {
	fileMetaService FileMetaService
	fileManager     FileManager
}

func NewFileControl(fileMetaService FileMetaService, fileManager FileManager) *FileControl {
	return &FileControl{
		fileMetaService: fileMetaService,
		fileManager:     fileManager,
	}
}

func (fc *FileControl) DownloadFile(ctx *fiber.Ctx) error {
	fileID, err := strconv.Atoi(ctx.Params("file_id"))
	if err != nil || fileID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	uri, err := fc.fileMetaService.GetFileMetaURI(ctx.Context(), fileID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	if err := fc.fileManager.DownloadFile(ctx, uri); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return fc.fileManager.DownloadFile(ctx, uri)
}

func (fc *FileControl) DeleteFile(ctx *fiber.Ctx) error {
	fileID, err := strconv.Atoi(ctx.Params("file_id"))
	if err != nil || fileID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	uri, err := fc.fileMetaService.GetFileMetaURI(ctx.Context(), fileID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	if err := fc.fileManager.DeleteFile(ctx, uri); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	if err := fc.fileMetaService.DeleteFileMetaByID(ctx.Context(), fileID); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}
