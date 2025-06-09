package admin

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

type AdminController interface {
	GetUsers(ctx *fiber.Ctx) error
	CreateUser(ctx *fiber.Ctx) error
}

type AdminControl struct {
	userService *service.UserService
	jwtS        *jwtservice.JWTService
}

func NewAdminControl(
	userService *service.UserService,
	jwtS *jwtservice.JWTService,
) *AdminControl {
	return &AdminControl{
		userService: userService,
		jwtS:        jwtS,
	}
}

func (ac *AdminControl) GetUsers(ctx *fiber.Ctx) error {
	token := ctx.Get(jwtservice.HeaderAuthorization)

	userID, err := ac.jwtS.GetUserID(token)
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

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

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
