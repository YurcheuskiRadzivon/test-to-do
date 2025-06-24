package user

import (
	"context"
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	GetUser(ctx context.Context, userID int) (string, string, error)
	UpdateUser(ctx context.Context, user entity.User) error
	DeleteUser(ctx context.Context, userID int) error
}

type AuthManager interface {
	GetUserID(ctx *fiber.Ctx) (int, error)
}

type EncryptManager interface {
	EncodePassword(password string) (string, error)
}

type FileMetaService interface {
	GetFileMetasIDByUserID(ctx context.Context, userID int) ([]int, error)
	GetFileMetaURI(ctx context.Context, id int) (string, error)
}

type FileManager interface {
	DeleteFile(ctx *fiber.Ctx, path string) error
}

type UserController interface {
	GetUser(ctx *fiber.Ctx) error
	UpdateUser(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
}

type UserControl struct {
	userService     UserService
	authManager     AuthManager
	encryptManager  EncryptManager
	fileMetaService FileMetaService
	fileManager     FileManager
}

func NewUserControl(
	userService UserService,
	authManager AuthManager,
	encryptManager EncryptManager,
	fileMetaService FileMetaService,
	fileManager FileManager,
) *UserControl {
	return &UserControl{
		userService:     userService,
		authManager:     authManager,
		encryptManager:  encryptManager,
		fileMetaService: fileMetaService,
		fileManager:     fileManager,
	}
}

func (uc *UserControl) GetUser(ctx *fiber.Ctx) error {

	userID, err := uc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	username, email, err := uc.userService.GetUser(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	return ctx.Status(http.StatusOK).JSON(response.UserData{
		Username: username,
		Email:    email,
	})
}

func (uc *UserControl) UpdateUser(ctx *fiber.Ctx) error {
	userID, err := uc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	var req request.OperationUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	hashedPassword, err := uc.encryptManager.EncodePassword(req.Password)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err = uc.userService.UpdateUser(ctx.Context(), entity.User{
		UserID:   userID,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (uc *UserControl) DeleteUser(ctx *fiber.Ctx) error {
	userID, err := uc.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err = uc.userService.DeleteUser(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}
