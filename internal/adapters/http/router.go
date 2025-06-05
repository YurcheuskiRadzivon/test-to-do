package http

import (
	"fmt"

	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewRoute(app *fiber.App, noteService *service.NoteService, userService *service.UserService, cfg *config.Config) {
	APIController := &APIController{
		app:         app,
		noteService: noteService,
		userService: userService,
	}

	app.Static("/swagger/swagger.yaml", "./docs/swagger.yaml")
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          fmt.Sprintf("%s/swagger/swagger.yaml", cfg.APP.DOMAIN),
		DeepLinking:  true,
		DocExpansion: "none",
		Title:        "API Documentation",
	}))

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
