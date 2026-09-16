package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/models"
	"xanders-gen-video/internal/services"
	"xanders-gen-video/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobHandler struct {
	queueService           *services.QueueService
	creditService          *services.CreditService
	promptService          *services.PromptService
	cinematicPromptService *services.CinematicPromptService
}

func insufficientCreditsResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusPaymentRequired).JSON(utils.APIResponse{
		Success: false,
		Error: fmt.Sprintf(
			"Insufficient credits. Each video costs %d credits. Request more credits at %s.",
			services.VideoGenerationCreditCost,
			services.CreditSupportEmail,
		),
		Data: fiber.Map{
			"credits_required": services.VideoGenerationCreditCost,
			"contact_email":    services.CreditSupportEmail,
		},
	})
}

func NewJobHandler(queue *services.QueueService, credit *services.CreditService, prompt *services.PromptService, cinematicPrompt *services.CinematicPromptService) *JobHandler {
	return &JobHandler{
		queueService:           queue,
		creditService:          credit,
		promptService:          prompt,
		cinematicPromptService: cinematicPrompt,
	}
}

func (h *JobHandler) CreateJob(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)

	// We now accept multipart/form-data for image uploads
	promptInputStr := c.FormValue("prompt_input")
	var promptInput map[string]interface{}
	if err := json.Unmarshal([]byte(promptInputStr), &promptInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "Invalid prompt_input JSON",
		})
	}

	var imageURL *string
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		// Save the uploaded file locally for now (in production, upload to GCS here)
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
				Success: false,
				Error:   "Could not create upload directory",
			})
		}

		filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), filepath.Ext(file.Filename))
		savePath := filepath.Join(uploadDir, filename)

		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
				Success: false,
				Error:   "Failed to save image",
			})
		}

		// Create a local URL reference for now (can be served via static fiber route)
		url := fmt.Sprintf("/uploads/%s", filename)
		imageURL = &url
	}

	// 1. Check the fixed generation price before spending provider resources.
	creditsRequired := services.VideoGenerationCreditCost
	if err := h.creditService.EnsureSufficientCredits(c.Context(), userIDStr, creditsRequired); err != nil {
		if errors.Is(err, services.ErrInsufficientCredits) {
			return insufficientCreditsResponse(c)
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   "Unable to verify credit balance",
		})
	}

	// 2. Assemble prompt
	optimizeCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	assembledPrompt, err := h.cinematicPromptService.GenerateCinematicPrompt(optimizeCtx, services.UIInputFromPromptInput(promptInput))
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(utils.APIResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to optimize prompt with Gemini: %v", err),
		})
	}

	// 3. Create job record & Deduct credits (Transactional)
	jobID := uuid.New().String()

	// Add imageURL to promptInput for reference if needed
	if imageURL != nil {
		promptInput["image_url"] = *imageURL
	}

	err = h.creditService.DeductCreditsAndCreateJob(c.Context(), userIDStr, creditsRequired, jobID, promptInput, assembledPrompt, imageURL)
	if err != nil {
		if errors.Is(err, services.ErrInsufficientCredits) {
			return insufficientCreditsResponse(c)
		}

		return c.Status(fiber.StatusPaymentRequired).JSON(utils.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	var createdJob models.Job
	if err := database.DB.WithContext(c.Context()).Where("id = ? AND user_id = ?", jobID, userIDStr).First(&createdJob).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   "Job creation was not persisted",
		})
	}

	// 4. Enqueue to Cloud Tasks (or local worker)
	_ = h.queueService.EnqueueVideoJob(c.Context(), jobID)

	return c.Status(fiber.StatusAccepted).JSON(utils.APIResponse{
		Success: true,
		Message: "Job accepted and queued",
		Data: fiber.Map{
			"job_id": jobID,
			"status": "in_queue",
		},
	})
}

func (h *JobHandler) GetJob(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	jobID := c.Params("id")

	var job models.Job
	result := database.DB.WithContext(c.Context()).Where("id = ? AND user_id = ?", jobID, userIDStr).First(&job)
	if result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
				Success: false,
				Error:   result.Error.Error(),
			})
		}

		var existingJob models.Job
		if err := database.DB.WithContext(c.Context()).Where("id = ?", jobID).First(&existingJob).Error; err == nil {
			return c.Status(fiber.StatusForbidden).JSON(utils.APIResponse{
				Success: false,
				Error:   "Job belongs to a different user",
			})
		}

		return c.Status(fiber.StatusNotFound).JSON(utils.APIResponse{
			Success: false,
			Error:   "Job not found",
		})
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Data:    job,
	})
}

func (h *JobHandler) ListJobs(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	if limit > 50 {
		limit = 50
	}
	offset := (page - 1) * limit

	var jobs []models.Job
	result := database.DB.WithContext(c.Context()).
		Where("user_id = ?", userIDStr).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   "Failed to fetch jobs",
		})
	}

	if jobs == nil {
		jobs = []models.Job{}
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Data:    jobs,
	})
}
