package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/gofiber/fiber/v2"
)

func (c *APIController) GetUsers(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Get("Authorization"))
	if err != nil || userID != 0 {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	users, err := c.userService.GetUsers(ctx.Context())
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(users)
}

func (c *APIController) GetUser(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Get("Authorization"))
	if err != nil || userID < 0 {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	user, err := c.userService.GetUser(ctx.Context(), userID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(user)
}

func (c *APIController) CreateUser(ctx *fiber.Ctx) error {
	var req request.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		log.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	id, err := c.userService.CreateUser(ctx.Context(), entity.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})

	if err != nil {
		log.Println(err)
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(id)
}

func (c *APIController) UpdateUser(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Get("Authorization"))
	if err != nil || userID < 0 {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	var req request.UpdateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = c.userService.UpdateUser(ctx.Context(), entity.User{
		UserID:   userID,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.UpdateNoteResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (c *APIController) DeleteUser(ctx *fiber.Ctx) error {
	userID, err := strconv.Atoi(ctx.Get("Authorization"))
	if err != nil || userID < 0 {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	err = c.userService.DeleteUser(ctx.Context(), userID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.DeleteNoteResponse{
		Message: response.MessageSuccsessfully,
	})
}
