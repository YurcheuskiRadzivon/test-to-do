package auth

import (
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	userIDHeader      = "user_id"
	ErrNotImplemented = "NOT_IMPLEMENTED"
)

type AuthController interface {
	Login(ctx *fiber.Ctx) error
	Register(ctx *fiber.Ctx) error
}

type AuthControl struct {
	userService *service.UserService
	jwtS        *jwtservice.JWTService
}

func NewAuthControl(
	userService *service.UserService,
	jwtS *jwtservice.JWTService,
) *AuthControl {
	return &AuthControl{
		userService: userService,
		jwtS:        jwtS,
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

	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(req.Password),
	)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	payload := jwt.MapClaims{
		userIDHeader: userID,
	}

	token, err := ac.jwtS.CreateToken(payload)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrJWT)
	}

	return ctx.Status(http.StatusOK).JSON(response.LoginResponse{
		Token: token,
	})
}

func (ac *AuthControl) Register(ctx *fiber.Ctx) error {
	return response.ErrorResponse(ctx, http.StatusBadRequest, ErrNotImplemented)
}
