package http

import (
	"log"
	"net/http"

	"github.com/YurcheuskiRadzivon/test-to-do/pkg/jwtservice"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func (c *APIController) AuthMiddleware(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return c.jwtS.GetJwtSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	userID, err := c.jwtS.GetUserID(tokenString)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	exist, err := c.userService.UserExistsByID(ctx.Context(), userID)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	if exist == false {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	return ctx.Next()
}

func (c *APIController) AuthAdminMiddleware(ctx *fiber.Ctx) error {
	tokenString := ctx.Get("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return c.jwtS.GetJwtSecretKey(), nil
	})

	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	userID, err := c.jwtS.GetUserID(tokenString)
	if err != nil {
		return errorResponse(ctx, http.StatusBadRequest, jwtservice.StatusInvalidToken)
	}

	log.Println(c.cfg.ADMIN.ID, userID)
	if userID != c.cfg.ADMIN.ID {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "You are not admin",
		})
	}

	return ctx.Next()
}
