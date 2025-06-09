package note

import (
	"log"
	"net/http"
	"strconv"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	authmanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/auth"
	encryptmanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/encrypt"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
)

type NoteController interface {
	GetNotes(ctx *fiber.Ctx) error
	GetNote(ctx *fiber.Ctx) error
	CreateNote(ctx *fiber.Ctx) error
	UpdateNote(ctx *fiber.Ctx) error
	DeleteNote(ctx *fiber.Ctx) error
}

type NoteControl struct {
	noteService    *service.NoteService
	authManager    authmanage.AuthManager
	encryptManager encryptmanage.EncryptManager
}

func NewNoteControl(
	noteService *service.NoteService,
	authManager authmanage.AuthManager,
	encryptManager encryptmanage.EncryptManager,

) *NoteControl {
	return &NoteControl{
		noteService:    noteService,
		authManager:    authManager,
		encryptManager: encryptManager,
	}
}

func (nc *NoteControl) GetNotes(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}
	notes, err := nc.noteService.GetNotes(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(notes)
}

func (nc *NoteControl) GetNote(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	noteID, err := strconv.Atoi(ctx.Params(jwtservice.ParamID))
	if err != nil || noteID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	note, err := nc.noteService.GetNote(ctx.Context(), noteID, userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(note)
}

func (nc *NoteControl) CreateNote(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	var req request.OperationNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = nc.noteService.CreateNote(ctx.Context(), entity.Note{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AuthorID:    userID,
	})

	if err != nil {
		log.Println(err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (nc *NoteControl) UpdateNote(ctx *fiber.Ctx) error {
	userID, err := nc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
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
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	noteID, err := strconv.Atoi(ctx.Params(jwtservice.ParamID))
	if err != nil || noteID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = nc.noteService.DeleteNote(ctx.Context(), noteID, userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}
