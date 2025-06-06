package http

import (
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	userIDHeader = "user_id"
)

func (c *APIController) Login(ctx *fiber.Ctx) error {
	var req request.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	userID, hashedPassword, err := c.userService.GetUserLoginParams(ctx.Context(), req.Username)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(req.Password),
	)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	payload := jwt.MapClaims{
		userIDHeader: userID,
	}

	token, err := c.jwtS.CreateToken(payload)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrJWT)
	}

	return ctx.Status(http.StatusOK).JSON(response.LoginResponse{
		Token: token,
	})
}
