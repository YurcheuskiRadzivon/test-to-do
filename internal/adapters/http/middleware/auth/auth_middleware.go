package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
)

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
}

type AuthMW struct {
	authManager AuthManager
	userService UserService
	cfg         *config.Config
}

func NewAuthMW(
	authManager AuthManager,
	userService UserService,
	cfg *config.Config,
) *AuthMW {
	return &AuthMW{
		authManager: authManager,
		userService: userService,
		cfg:         cfg,
	}
}

func (am *AuthMW) AuthUserMiddleware(ctx *fiber.Ctx) error {
	err := am.authManager.Validate(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	userID, err := am.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	exist, err := am.userService.UserExistsByID(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	if exist == false {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	return ctx.Next()
}

func (am *AuthMW) AuthAdminMiddleware(ctx *fiber.Ctx) error {
	err := am.authManager.Validate(ctx)

	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	userID, err := am.authManager.GetUserID(ctx)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	log.Println(am.cfg.ADMIN.ID, userID)
	if userID != am.cfg.ADMIN.ID {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "You are not admin",
		})
	}

	return ctx.Next()
}
