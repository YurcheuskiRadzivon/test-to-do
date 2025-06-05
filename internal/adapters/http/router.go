package http

import (
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/gofiber/fiber/v2"
)

func NewRoute(app *fiber.App, noteService *service.NoteService, userService *service.UserService) {
	APIController := &APIController{
		app:         app,
		noteService: noteService,
		userService: userService,
	}

	app.Get("/note/:id", APIController.GetNote)
	app.Get("/notes", APIController.GetNotes)
	app.Post("/note", APIController.CreateNote)
	app.Delete("/note/:id", APIController.DeleteNote)
	app.Put("/note/:id", APIController.UpdateNote)

	app.Get("/user", APIController.GetUser)
	app.Get("/users", APIController.GetUsers)
	app.Post("/user", APIController.CreateUser)
	app.Delete("/user", APIController.DeleteUser)
	app.Put("/user", APIController.UpdateUser)

}
