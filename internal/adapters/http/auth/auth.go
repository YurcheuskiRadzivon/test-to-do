package auth

import (
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	authmanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/auth"
	encryptmanage "github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/managers/encrypt"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/gofiber/fiber/v2"
)

const (
	UserIDHeader      = "user_id"
	ErrNotImplemented = "NOT_IMPLEMENTED"
)

type AuthController interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

type AuthControl struct {
	userService    *service.UserService
	authManager    authmanage.AuthManager
	encryptManager encryptmanage.EncryptManager
}

func NewAuthControl(
	userService *service.UserService,
	authManager authmanage.AuthManager,
	encryptManager encryptmanage.EncryptManager,
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
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ac.authManager.CreateAuthResponse(ctx, userID)

}

func (ac *AuthControl) Register(ctx *fiber.Ctx) error {
	return response.ErrorResponse(ctx, http.StatusBadRequest, ErrNotImplemented)
}
