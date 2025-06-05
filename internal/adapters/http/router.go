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
	noteGroup := app.Group("/manage")
	{
		noteGroup.Get("/note/:id", APIController.GetNote)
		noteGroup.Get("/notes", APIController.GetNotes)
		noteGroup.Post("/note", APIController.CreateNote)
		noteGroup.Delete("/note/:id", APIController.DeleteNote)
		noteGroup.Put("/note/:id", APIController.UpdateNote)
	}

	userGroup := app.Group("/account")
	{
		userGroup.Get("/user", APIController.GetUser)
		userGroup.Get("/users", APIController.GetUsers)
		userGroup.Post("/user", APIController.CreateUser)
		userGroup.Delete("/user", APIController.DeleteUser)
		userGroup.Put("/user", APIController.UpdateUser)
	}

}
