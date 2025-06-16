package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/gofiber/fiber/v2"
)

const (
	fileIDParam = "file_id"
)

type FileMetaService interface {
	FileMetasExistsByIDAndUserID(ctx context.Context, id int, userID int) (bool, error)
}

type UserService interface {
	UserExistsByID(ctx context.Context, userID int) (bool, error)
}

type AuthManager interface {
	GetUserID(ctx *fiber.Ctx) (int, error)
	Validate(ctx *fiber.Ctx) error
}

type AuthMiddleware interface {
	AuthUserMiddleware(ctx *fiber.Ctx) error
	AuthAdminMiddleware(ctx *fiber.Ctx) error
	AuthFileActionMiddleware(ctx *fiber.Ctx) error
}

type AuthMW struct {
	fileMetaService FileMetaService
	authManager     AuthManager
	userService     UserService
	cfg             *config.Config
}

func NewAuthMW(
	fileMetaService FileMetaService,
	authManager AuthManager,
	userService UserService,
	cfg *config.Config,
) *AuthMW {
	return &AuthMW{
		fileMetaService: fileMetaService,
		authManager:     authManager,
		userService:     userService,
		cfg:             cfg,
	}
}

func (am *AuthMW) AuthUserMiddleware(ctx *fiber.Ctx) error {
	err := am.authManager.Validate(ctx)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	userID, err := am.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	exist, err := am.userService.UserExistsByID(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	if exist == false {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	return ctx.Next()
}

func (am *AuthMW) AuthAdminMiddleware(ctx *fiber.Ctx) error {
	err := am.authManager.Validate(ctx)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	userID, err := am.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	if userID != am.cfg.ADMIN.ID {
		return response.ErrorResponse(ctx, http.StatusForbidden, response.ErrInvalidToken)
	}

	return ctx.Next()
}

func (am *AuthMW) AuthFileActionMiddleware(ctx *fiber.Ctx) error {
	err := am.authManager.Validate(ctx)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	userID, err := am.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	fileID, err := strconv.Atoi(ctx.Params(fileIDParam))
	if err != nil || fileID == 0 {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	exist, err := am.fileMetaService.FileMetasExistsByIDAndUserID(ctx.Context(), fileID, userID)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	if exist == false {
		return response.ErrorResponse(ctx, http.StatusForbidden, response.ErrInvalidToken)
	}

	return ctx.Next()
}
