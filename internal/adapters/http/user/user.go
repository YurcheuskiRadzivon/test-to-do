package user

import (
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/request"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/entity"
	"github.com/YurcheuskiRadzivon/test-to-do/internal/core/service"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	GetUser(ctx *fiber.Ctx) error
	UpdateUser(ctx *fiber.Ctx) error
	DeleteUser(ctx *fiber.Ctx) error
}

type UserControl struct {
	userService *service.UserService
	jwtS        *jwtservice.JWTService
}

func NewUserControl(
	userService *service.UserService,
	jwtS *jwtservice.JWTService,
) *UserControl {
	return &UserControl{
		userService: userService,
		jwtS:        jwtS,
	}
}

func (uc *UserControl) GetUser(ctx *fiber.Ctx) error {
	token := ctx.Get(jwtservice.HeaderAuthorization)

	userID, err := uc.jwtS.GetUserID(token)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	username, email, err := uc.userService.GetUser(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}
	return ctx.Status(http.StatusOK).JSON(response.UserData{
		Username: username,
		Email:    email,
	})
}

func (uc *UserControl) UpdateUser(ctx *fiber.Ctx) error {
	token := ctx.Get(jwtservice.HeaderAuthorization)

	userID, err := uc.jwtS.GetUserID(token)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	var req request.OperationUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	err = uc.userService.UpdateUser(ctx.Context(), entity.User{
		UserID:   userID,
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	})

	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}

func (uc *UserControl) DeleteUser(ctx *fiber.Ctx) error {
	token := ctx.Get(jwtservice.HeaderAuthorization)
	userID, err := uc.jwtS.GetUserID(token)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	err = uc.userService.DeleteUser(ctx.Context(), userID)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrInvalidRequest)
	}

	return ctx.Status(http.StatusOK).JSON(response.MessageResponse{
		Message: response.MessageSuccsessfully,
	})
}
