package http

import (
	"fmt"

	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewRoute(
	app *fiber.App,
	noteService *service.NoteService,
	userService *service.UserService,
	cfg *config.Config,
	jwtS *jwtservice.JWTService,
) {
	APIController := &APIController{
		app:         app,
		noteService: noteService,
		userService: userService,
		jwtS:        jwtS,
		cfg:         cfg,
	}

	app.Static("/swagger/swagger.yaml", "./docs/swagger.yaml")
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          fmt.Sprintf("%s/swagger/swagger.yaml", cfg.APP.DOMAIN),
		DeepLinking:  true,
		DocExpansion: "none",
		Title:        "API Documentation",
	}))

	authGroup := app.Group("/auth")
	{
		authGroup.Post("/login", APIController.Login)
	}

	adminGroup := app.Group("/admin")
	{
		adminGroup.Use(APIController.AuthAdminMiddleware)
		adminGroup.Get("/users", APIController.GetUsers)
		adminGroup.Post("/user", APIController.CreateUser)
	}

	userGroup := app.Group("/account")
	{
		userGroup.Use(APIController.AuthMiddleware)
		userGroup.Get("/user", APIController.GetUser)
		userGroup.Delete("/user", APIController.DeleteUser)
		userGroup.Put("/user", APIController.UpdateUser)
	}

	noteGroup := app.Group("/manage")
	{
		noteGroup.Use(APIController.AuthMiddleware)
		noteGroup.Get("/note/:id", APIController.GetNote)
		noteGroup.Get("/notes", APIController.GetNotes)
		noteGroup.Post("/note", APIController.CreateNote)
		noteGroup.Delete("/note/:id", APIController.DeleteNote)
		noteGroup.Put("/note/:id", APIController.UpdateNote)
	}

}
