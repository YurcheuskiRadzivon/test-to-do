package middleware

import (
	"context"
	"log"
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
		return response.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	userID, err := am.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	exist, err := am.userService.UserExistsByID(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	if exist == false {
		log.Printf("Faled to auth user: %v - exist", exist)
		return response.ErrorResponse(ctx, http.StatusUnauthorized, response.ErrInvalidToken)
	}

	return ctx.Next()
}

func (am *AuthMW) AuthAdminMiddleware(ctx *fiber.Ctx) error {
	err := am.authManager.Validate(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	userID, err := am.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	if userID != am.cfg.ADMIN.ID {
		log.Printf("Faled, user is not admin: %v - userID", userID)
		return response.ErrorResponse(ctx, http.StatusForbidden, response.ErrInvalidToken)
	}

	return ctx.Next()
}

func (am *AuthMW) AuthFileActionMiddleware(ctx *fiber.Ctx) error {
	err := am.authManager.Validate(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	userID, err := am.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
	}

	fileID, err := strconv.Atoi(ctx.Params(fileIDParam))
	if err != nil || fileID <= 0 {
		log.Printf("Faled to get file id or invalid file id: %v - fileID, %v - err", fileID, err)
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	exist, err := am.fileMetaService.FileMetasExistsByIDAndUserID(ctx.Context(), fileID, userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if exist == false {
		log.Printf("Faled, user is not have rights to work with files: %v - userID", userID)
		return response.ErrorResponse(ctx, http.StatusForbidden, response.ErrInvalidToken)
	}

	return ctx.Next()
}
