package middleware

import (
	"log"
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/config"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type AuthMiddleware interface {
	AuthUserMiddleware(ctx *fiber.Ctx) error
	AuthAdminMiddleware(ctx *fiber.Ctx) error
}

type AuthMW struct {
	userService *service.UserService
	jwtS        *jwtservice.JWTService
	cfg         *config.Config
}

func NewAuthMW(
	userService *service.UserService,
	jwtS *jwtservice.JWTService,
	cfg *config.Config,
) *AuthMW {
	return &AuthMW{
		userService: userService,
		jwtS:        jwtS,
		cfg:         cfg,
	}
}

func (am *AuthMW) AuthUserMiddleware(ctx *fiber.Ctx) error {
	tokenString := ctx.Get(jwtservice.HeaderAuthorization)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return am.jwtS.GetJwtSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	userID, err := am.jwtS.GetUserID(tokenString)
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
	tokenString := ctx.Get(jwtservice.HeaderAuthorization)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return am.jwtS.GetJwtSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	userID, err := am.jwtS.GetUserID(tokenString)
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
