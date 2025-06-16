package auth

import (
	"context"
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/gofiber/fiber/v2"
)

const (
	UserIDHeader = "user_id"
)

type UserService interface {
	GetUserLoginParams(ctx context.Context, username string) (int, string, error)
}

type AuthManager interface {
	CreateAuthResponse(ctx *fiber.Ctx, id int) error
}

type EncryptManager interface {
	CheckPassword(password, hashedPassword string) error
}

type AuthController interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

type AuthControl struct {
	userService    UserService
	authManager    AuthManager
	encryptManager EncryptManager
}

func NewAuthControl(
	userService UserService,
	authManager AuthManager,
	encryptManager EncryptManager,
) *AuthControl {
	return &AuthControl{
		userService:    userService,
		authManager:    authManager,
		encryptManager: encryptManager,
	}
}

func (ac *AuthControl) Login(ctx *fiber.Ctx) error {
	var req request.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	userID, hashedPassword, err := ac.userService.GetUserLoginParams(ctx.Context(), req.Username)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	if ac.encryptManager.CheckPassword(req.Password, hashedPassword); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrIvalidPassword)
	}

	return ac.authManager.CreateAuthResponse(ctx, userID)
}

func (ac *AuthControl) Register(ctx *fiber.Ctx) error {
	return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrNotImplemented)
}
