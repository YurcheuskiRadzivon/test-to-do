package http

import (
	"log"
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func (c *APIController) GetUsers(ctx *fiber.Ctx) error {
	token := ctx.Get("Authorization")

	userID, err := c.jwtS.GetUserID(token)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	if userID != 0 {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	users, err := c.userService.GetUsers(ctx.Context())
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(users)
}

func (c *APIController) CreateUser(ctx *fiber.Ctx) error {
	var req request.OperationUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Println(1, err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		log.Println(2, err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	_, err = c.userService.CreateUser(ctx.Context(), entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	})

	if err != nil {
		log.Println(3, err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(response.CreateUserResponse{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
}
