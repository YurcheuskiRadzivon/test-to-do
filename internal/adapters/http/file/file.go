package file

import (
	"context"

	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/gofiber/fiber/v2"
)

const (
	fileIDParam   = "file_id"
	NameFormFiles = "[]files"
	noteIDParam   = "note_id"
	contentType   = "Content-Type"
)

type NoteService interface {
	GetNote(ctx context.Context, noteID int, authorID int) (entity.Note, error)
}

type AuthManager interface {
	GetUserID(ctx *fiber.Ctx) (int, error)
}

type FileMetaService interface {
	GetFileMetaURI(ctx context.Context, id int) (string, error)
	DeleteFileMetaByID(ctx context.Context, id int) error
	CreateFileMeta(ctx context.Context, fileMeta entity.FileMeta) error
}

type FileManager interface {
	DownloadFile(ctx *fiber.Ctx, path string) error
	DeleteFile(ctx *fiber.Ctx, path string) error
	UploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, error)
}

type FileController interface {
	DownloadFile(ctx *fiber.Ctx) error
	UploadFiles(ctx *fiber.Ctx) error
	DeleteFile(ctx *fiber.Ctx) error
}

type FileControl struct {
	noteService     NoteService
	fileMetaService FileMetaService
	fileManager     FileManager
	authManager     AuthManager
}

func NewFileControl(fileMetaService FileMetaService, fileManager FileManager, authManager AuthManager, noteService NoteService) *FileControl {
	return &FileControl{
		noteService:     noteService,
		fileMetaService: fileMetaService,
		fileManager:     fileManager,
		authManager:     authManager,
	}
}

func (fc *FileControl) DownloadFile(ctx *fiber.Ctx) error {
	fileID, err := strconv.Atoi(ctx.Params(fileIDParam))
	if err != nil || fileID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidFileID)
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
	fileID, err := strconv.Atoi(ctx.Params(fileIDParam))
	if err != nil || fileID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidFileID)
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

func (fc *FileControl) UploadFiles(ctx *fiber.Ctx) error {
	noteID, err := strconv.Atoi(ctx.Params(noteIDParam))
	if err != nil || noteID <= 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidNoteID)
	}
	userID, err := fc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidToken)
	}

	_, err = fc.noteService.GetNote(ctx.Context(), noteID, userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	files := form.File[NameFormFiles]

	uriList, err := fc.fileManager.UploadFiles(ctx, files)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	for i := range uriList {
		err := fc.fileMetaService.CreateFileMeta(ctx.Context(), entity.FileMeta{
			ContentType: files[i].Header.Get(contentType),
			OwnerType:   entity.OwnerNote,
			OwnerID:     noteID,
			UserID:      userID,
			URI:         uriList[i],
		})
		if err != nil {
			return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
		}
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})

}
