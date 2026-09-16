package middleware

import (
	"strings"

	"xanders-gen-video/internal/config"
	"xanders-gen-video/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.APIResponse{
				Success: false,
				Error:   "Missing or invalid authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		cfg := config.LoadConfig()

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.APIResponse{
				Success: false,
				Error:   "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.APIResponse{
				Success: false,
				Error:   "Invalid token claims",
			})
		}

		c.Locals("user_id", claims["sub"])
		return c.Next()
	}
}
