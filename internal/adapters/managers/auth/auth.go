package authmanage

import (
	"errors"
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/internal/adapters/http/response"
	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type AuthManager interface {
	GetUserID(ctx *fiber.Ctx) (int, error)
	Validate(ctx *fiber.Ctx) error
	CreateAuthResponse(ctx *fiber.Ctx, id int) error
}

type AuthManage struct {
	jwtS *jwtservice.JWTService
}

func NewAuthManage(jwtS *jwtservice.JWTService) *AuthManage {
	return &AuthManage{
		jwtS: jwtS,
	}
}

func (am *AuthManage) GetUserID(ctx *fiber.Ctx) (int, error) {
	token := ctx.Get(jwtservice.HeaderAuthorization)
	userID, err := am.jwtS.GetUserID(token)
	if err != nil {
		return 0, errors.New(jwtservice.StatusInvalidToken)
	}
	return userID, nil
}
func (am *AuthManage) Validate(ctx *fiber.Ctx) error {
	tokenString := ctx.Get(jwtservice.HeaderAuthorization)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return am.jwtS.GetJwtSecretKey(), nil
	})
	if err != nil || !token.Valid {
		return errors.New("INVALID_OR_EXPIRED_TOKEN")

	}

	return nil
}
func (am *AuthManage) CreateAuthResponse(ctx *fiber.Ctx, id int) error {
	payload := jwt.MapClaims{
		"user_id": id,
	}

	token, err := am.jwtS.CreateToken(payload)
	if err != nil {
		return response.ErrorResponse(ctx, http.StatusBadRequest, response.ErrJWT)
	}

	return ctx.Status(http.StatusOK).JSON(response.LoginResponse{
		Token: token,
	})
}
