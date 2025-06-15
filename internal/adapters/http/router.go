package http

import (
	"fmt"

	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/admin"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/auth"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/file"
	middleware "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/middleware/auth"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/note"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewRoute(
	app *fiber.App,
	noteController note.NoteController,
	userController user.UserController,
	adminController admin.AdminController,
	authController auth.AuthController,
	authMiddleware middleware.AuthMiddleware,
	fileController file.FileController,
	cfg *config.Config,
) {
	APIC := &APIController{
		app:             app,
		noteController:  noteController,
		userController:  userController,
		adminController: adminController,
		authController:  authController,
		authMiddleware:  authMiddleware,
		fileController:  fileController,
		cfg:             cfg,
	}

	//app.Static("/uploadfiles/", "./uploaded_files/")
	/*app.Get("/uploadfiles/:filename", func(c *fiber.Ctx) error {
		filename := c.Params("filename")
		return c.SendFile("./uploaded_files/" + filename)
	})*/

	app.Static("/swagger/swagger.yaml", "./docs/swagger.yaml")
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          fmt.Sprintf("%s/swagger/swagger.yaml", cfg.APP.DOMAIN),
		DeepLinking:  true,
		DocExpansion: "none",
		Title:        "API Documentation",
	}))

	authGroup := app.Group("/auth")
	{
		authGroup.Post("/login", APIC.authController.Login)
	}

	adminGroup := app.Group("/admin")
	{
		adminGroup.Use(APIC.authMiddleware.AuthAdminMiddleware)
		adminGroup.Get("/users", APIC.adminController.GetUsers)
		adminGroup.Post("/user", APIC.adminController.CreateUser)
	}

	userGroup := app.Group("/account")
	{
		userGroup.Use(APIC.authMiddleware.AuthUserMiddleware)
		userGroup.Get("/user", APIC.userController.GetUser)
		userGroup.Delete("/user", APIC.userController.DeleteUser)
		userGroup.Put("/user", APIC.userController.UpdateUser)
	}

	noteGroup := app.Group("/manage")
	{
		noteGroup.Use(APIC.authMiddleware.AuthUserMiddleware)

		noteGroup.Get("/note/:id", APIC.noteController.GetNote)
		noteGroup.Get("/notes", APIC.noteController.GetNotes)
		noteGroup.Post("/note", APIC.noteController.CreateNote)
		noteGroup.Delete("/note/:id", APIC.noteController.DeleteNote)
		noteGroup.Put("/note/:id", APIC.noteController.UpdateNote)

		noteGroup.Get("/file/:file_id", APIC.fileController.DownloadFile)
		noteGroup.Delete("/file/:file_id", APIC.fileController.DeleteFile)
	}

}
