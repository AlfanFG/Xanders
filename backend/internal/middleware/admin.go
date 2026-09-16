package middleware

import (
	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/models"
	"xanders-gen-video/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.APIResponse{
				Success: false,
				Error:   "Unauthorized",
			})
		}

		var user models.User
		if err := database.DB.WithContext(c.Context()).Where("id = ?", userID).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.APIResponse{
				Success: false,
				Error:   "User not found",
			})
		}

		if user.Plan != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(utils.APIResponse{
				Success: false,
				Error:   "Admin access required",
			})
		}

		return c.Next()
	}
}
