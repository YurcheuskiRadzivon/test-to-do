package http

import (
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, noteService *service.NoteService) {
	APIController := &APIController{
		app:         app,
		noteService: noteService,
	}
	_ = APIController
	app.Get("/note/:id", APIController.GetNote)
	app.Get("/notes", APIController.GetNotes)
	app.Post("/note", APIController.CreateNote)
	app.Delete("/note/:id", APIController.DeleteNote)
	app.Put("/note/:id", APIController.UpdateNote)

}
