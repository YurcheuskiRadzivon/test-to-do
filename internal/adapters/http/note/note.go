package note

import (
	"context"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
)

const (
	nameFormTitle       = "title"
	nameFormDescription = "description"
	nameFormStatus      = "status"
	NameFormFiles       = "[]files"
	contentType         = "Content-Type"
)

type FileManager interface {
	UploadFiles(ctx *fiber.Ctx, files []*multipart.FileHeader) ([]string, error)
	DeleteFile(ctx *fiber.Ctx, path string) error
}

type AuthManager interface {
	GetUserID(ctx *fiber.Ctx) (int, error)
}

type FileMetaService interface {
	CreateFileMeta(ctx context.Context, fileMeta entity.FileMeta) error
	GetFileMetas(ctx context.Context) ([]entity.FileMeta, error)
	GetFileMetaIDByID(ctx context.Context, ownerType string, ownerID int) ([]int, error)
	GetFileMetaURI(ctx context.Context, id int) (string, error)
}

type NoteController interface {
	GetNotes(ctx *fiber.Ctx) error
	GetNote(ctx *fiber.Ctx) error
	CreateNote(ctx *fiber.Ctx) error
	UpdateNote(ctx *fiber.Ctx) error
	DeleteNote(ctx *fiber.Ctx) error
}

type NoteControl struct {
	fileMetaService FileMetaService
	fileManager     FileManager
	noteService     *service.NoteService
	authManager     AuthManager
}

func NewNoteControl(
	fileMetaService FileMetaService,
	noteService *service.NoteService,
	authManager AuthManager,
	fileManager FileManager,
) *NoteControl {
	return &NoteControl{
		noteService:     noteService,
		authManager:     authManager,
		fileManager:     fileManager,
		fileMetaService: fileMetaService,
	}
}

func (nc *NoteControl) GetNotes(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	res := make([]request.GetNoteReq, 0)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	notes, err := nc.noteService.GetNotes(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	for _, note := range notes {
		fileMetasID, err := nc.fileMetaService.GetFileMetaIDByID(ctx.Context(), string(entity.OwnerNote), note.NoteID)
		if err != nil {
			return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
		}
		res = append(res, request.GetNoteReq{
			Noteinfo:   note,
			FileIDList: fileMetasID,
		})
	}
	return ctx.Status(http.StatusOK).JSON(res)
}

func (nc *NoteControl) GetNote(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	noteID, err := strconv.Atoi(ctx.Params(jwtservice.ParamID))
	if err != nil || noteID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	note, err := nc.noteService.GetNote(ctx.Context(), noteID, userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	fileMetasID, err := nc.fileMetaService.GetFileMetaIDByID(ctx.Context(), string(entity.OwnerNote), note.NoteID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(request.GetNoteReq{
		Noteinfo:   note,
		FileIDList: fileMetasID,
	})
}

func (nc *NoteControl) CreateNote(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	title := ctx.FormValue(nameFormTitle)
	description := ctx.FormValue(nameFormDescription)
	status := ctx.FormValue(nameFormStatus)

	form, err := ctx.MultipartForm()
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	files := form.File[NameFormFiles]

	uriList, err := nc.fileManager.UploadFiles(ctx, files)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	noteID, err := nc.noteService.CreateNote(ctx.Context(), entity.Note{
		Title:       title,
		Description: description,
		Status:      status,
		AuthorID:    userID,
	})

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	for i := range uriList {
		err := nc.fileMetaService.CreateFileMeta(ctx.Context(), entity.FileMeta{
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

func (nc *NoteControl) UpdateNote(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	noteID, err := strconv.Atoi(ctx.Params(jwtservice.HeaderAuthorization))
	if err != nil || noteID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	var req request.OperationNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = nc.noteService.UpdateNote(ctx.Context(), entity.Note{
		NoteID:      noteID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AuthorID:    userID,
	})

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (nc *NoteControl) DeleteNote(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	noteID, err := strconv.Atoi(ctx.Params(jwtservice.ParamID))
	if err != nil || noteID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	fileMetasID, err := nc.fileMetaService.GetFileMetaIDByID(ctx.Context(), string(entity.OwnerNote), noteID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	for _, fileMetaID := range fileMetasID {
		uri, err := nc.fileMetaService.GetFileMetaURI(ctx.Context(), fileMetaID)
		if err != nil {
			return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
		}

		if err = nc.fileManager.DeleteFile(ctx, uri); err != nil {
			return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		}
	}

	err = nc.noteService.DeleteNote(ctx.Context(), noteID, userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}
