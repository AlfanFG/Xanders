package handlers

import (
	"xanders-gen-video/internal/services"
	"xanders-gen-video/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "Email and password are required",
		})
	}

	user, token, err := h.authService.RegisterUser(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(utils.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	user, token, err := h.authService.LoginUser(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.APIResponse{
			Success: false,
			Error:   "Invalid credentials",
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.APIResponse{
		Success: true,
		Message: "Login successful",
		Data: fiber.Map{
			"user":  user,
			"token": token,
		},
	})
}
