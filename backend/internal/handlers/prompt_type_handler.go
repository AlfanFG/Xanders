package handlers

import (
	"errors"

	"xanders-gen-video/internal/services"
	"xanders-gen-video/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type PromptTypeHandler struct {
	promptTypeService *services.PromptTypeService
}

type paginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func NewPromptTypeHandler(promptTypeService *services.PromptTypeService) *PromptTypeHandler {
	return &PromptTypeHandler{promptTypeService: promptTypeService}
}

func (h *PromptTypeHandler) List(c *fiber.Ctx) error {
	result, err := h.promptTypeService.List(c.Context(), services.PromptTypeListRequest{
		Page:           c.QueryInt("page", 1),
		Limit:          c.QueryInt("limit", 20),
		Type:           c.Query("type", ""),
		Search:         c.Query("search", ""),
		IncludeDeleted: c.QueryBool("include_deleted", false),
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   "Failed to fetch prompt types",
		})
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Data: fiber.Map{
			"items": result.Items,
			"pagination": paginationResponse{
				Page:       result.Page,
				Limit:      result.Limit,
				Total:      result.Total,
				TotalPages: result.TotalPages,
			},
		},
	})
}

func (h *PromptTypeHandler) GetByID(c *fiber.Ctx) error {
	item, err := h.promptTypeService.GetByID(c.Context(), c.Params("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(utils.APIResponse{
				Success: false,
				Error:   "Prompt type not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   "Failed to fetch prompt type",
		})
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Data:    item,
	})
}

func (h *PromptTypeHandler) Create(c *fiber.Ctx) error {
	var payload services.PromptTypePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	item, err := h.promptTypeService.Create(c.Context(), payload)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(utils.APIResponse{
		Success: true,
		Message: "Prompt type created",
		Data:    item,
	})
}

func (h *PromptTypeHandler) Update(c *fiber.Ctx) error {
	var payload services.PromptTypePayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	item, err := h.promptTypeService.Update(c.Context(), c.Params("id"), payload)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(utils.APIResponse{
				Success: false,
				Error:   "Prompt type not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Message: "Prompt type updated",
		Data:    item,
	})
}

func (h *PromptTypeHandler) Delete(c *fiber.Ctx) error {
	err := h.promptTypeService.Delete(c.Context(), c.Params("id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(utils.APIResponse{
				Success: false,
				Error:   "Prompt type not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   "Failed to delete prompt type",
		})
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Message: "Prompt type deleted",
	})
}
