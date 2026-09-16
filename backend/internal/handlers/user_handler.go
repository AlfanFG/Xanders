package handlers

import (
	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/models"
	"xanders-gen-video/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	var user models.User
	result := database.DB.WithContext(c.Context()).Where("id = ?", userIDStr).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.APIResponse{
			Success: false,
			Error:   "User not found",
		})
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Data:    user,
	})
}

func (h *UserHandler) GetCreditHistory(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
	}

	var transactions []models.CreditHistory
	result := database.DB.WithContext(c.Context()).
		Where("user_id = ?", userIDStr).
		Order("created_at DESC").
		Limit(50).
		Find(&transactions)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   "Failed to fetch credit history",
		})
	}

	if transactions == nil {
		transactions = []models.CreditHistory{}
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Data:    transactions,
	})
}

