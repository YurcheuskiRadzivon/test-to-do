package http

import (
	"log"
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/gofiber/fiber/v2"
)

func (c *APIController) GetNotes(ctx *fiber.Ctx) error {
	notes, err := c.noteService.GetNotes(ctx.Context())
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(response.GetNotesResponse{
		Notes: notes,
	})
}

func (c *APIController) GetNote(ctx *fiber.Ctx) error {
	var req request.GetNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	note, err := c.noteService.GetNote(ctx.Context(), req.NoteID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(response.GetNoteResponse{
		Note: note,
	})
}

func (c *APIController) CreateNote(ctx *fiber.Ctx) error {
	var req request.CreateNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	err := c.noteService.CreateNote(ctx.Context(), entity.Note{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		log.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(response.CreateNoteResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (c *APIController) UpdateNote(ctx *fiber.Ctx) error {
	var req request.UpdateNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	err := c.noteService.UpdateNote(ctx.Context(), entity.Note{
		NoteID:      req.ID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(response.UpdateNoteResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (c *APIController) DeleteNote(ctx *fiber.Ctx) error {
	var req request.DeleteNoteRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	err := c.noteService.DeleteNote(ctx.Context(), req.NoteID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(response.DeleteNoteResponse{
		Message: response.MessageSuccsessfully,
	})
}
