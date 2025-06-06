package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
)

func (c *APIController) GetNotes(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	userID, err := c.jwtS.GetUserID(token)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}
	notes, err := c.noteService.GetNotes(ctx.Context(), userID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(notes)
}

func (c *APIController) GetNote(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	userID, err := c.jwtS.GetUserID(token)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	noteID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || noteID == 0 {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	note, err := c.noteService.GetNote(ctx.Context(), noteID, userID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(note)
}

func (c *APIController) CreateNote(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	userID, err := c.jwtS.GetUserID(token)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	var req request.OperationNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = c.noteService.CreateNote(ctx.Context(), entity.Note{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AuthorID:    userID,
	})

	if err != nil {
		log.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (c *APIController) UpdateNote(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	userID, err := c.jwtS.GetUserID(token)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	noteID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || noteID == 0 {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	var req request.OperationNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = c.noteService.UpdateNote(ctx.Context(), entity.Note{
		NoteID:      noteID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		AuthorID:    userID,
	})

	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (c *APIController) DeleteNote(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	userID, err := c.jwtS.GetUserID(token)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	noteID, err := strconv.Atoi(ctx.Params("id"))
	if err != nil || noteID == 0 {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = c.noteService.DeleteNote(ctx.Context(), noteID, userID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}
