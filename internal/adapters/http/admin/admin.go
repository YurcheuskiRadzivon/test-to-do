package admin

import (
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
)

type AuthManager interface {
	GetUserID(ctx *fiber.Ctx) (int, error)
}

type EncryptManager interface {
	EncodePassword(password string) (string, error)
}

type AdminController interface {
	GetUsers(ctx *fiber.Ctx) error
	CreateUser(ctx *fiber.Ctx) error
}

type AdminControl struct {
	userService    *service.UserService
	authManager    AuthManager
	encryptManager EncryptManager
}

func NewAdminControl(
	userService *service.UserService,
	authManager AuthManager,
	encryptManager EncryptManager,
) *AdminControl {
	return &AdminControl{
		userService:    userService,
		authManager:    authManager,
		encryptManager: encryptManager,
	}
}

func (ac *AdminControl) GetUsers(ctx *fiber.Ctx) error {
	userID, err := ac.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	if userID != 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	users, err := ac.userService.GetUsers(ctx.Context())
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(users)
}

func (ac *AdminControl) CreateUser(ctx *fiber.Ctx) error {
	var req request.OperationUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	hashedPassword, err := ac.encryptManager.EncodePassword(req.Password)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	_, err = ac.userService.CreateUser(ctx.Context(), entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	})

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.CreateUserResponse{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
}
