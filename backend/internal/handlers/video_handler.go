package handlers

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/models"
	"xanders-gen-video/internal/services"
	"xanders-gen-video/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	VideoStatusPending          = "PENDING"
	VideoStatusPromptGenerated  = "PROMPT_GENERATED"
	VideoStatusProcessingVideo  = "PROCESSING_VIDEO"
	VideoStatusCompleted        = "COMPLETED"
	VideoStatusFailed           = "FAILED"
	cinematicWorkflowTimeout    = 90 * time.Second
	cinematicStatusQueryTimeout = 5 * time.Second
)

type VideoHandler struct {
	cinematicPromptService *services.CinematicPromptService
}

func NewVideoHandler(cinematicPromptService *services.CinematicPromptService) *VideoHandler {
	return &VideoHandler{cinematicPromptService: cinematicPromptService}
}

func (h *VideoHandler) Generate(c *fiber.Ctx) error {
	var input services.UIInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	input.Product = strings.TrimSpace(input.Product)
	input.Environment = strings.TrimSpace(input.Environment)
	input.Style = strings.TrimSpace(input.Style)
	input.Motion = strings.TrimSpace(input.Motion)

	if input.Product == "" || input.Environment == "" || input.Style == "" || input.Motion == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "product, environment, style, and motion are required",
		})
	}

	job := models.VideoGenerationJob{
		JobID:         uuid.NewString(),
		UIProduct:     input.Product,
		UIEnvironment: input.Environment,
		UIStyle:       input.Style,
		UIMotion:      input.Motion,
		Status:        VideoStatusPending,
	}

	if err := database.DB.WithContext(c.Context()).Create(&job).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   "Failed to create video generation job",
		})
	}

	go h.runCinematicWorkflow(job.JobID, input)

	return c.Status(fiber.StatusAccepted).JSON(utils.APIResponse{
		Success: true,
		Message: "Video generation job accepted",
		Data: fiber.Map{
			"job_id": job.JobID,
			"status": job.Status,
		},
	})
}

func (h *VideoHandler) Status(c *fiber.Ctx) error {
	jobID := strings.TrimSpace(c.Params("job_id"))
	if jobID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.APIResponse{
			Success: false,
			Error:   "job_id is required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), cinematicStatusQueryTimeout)
	defer cancel()

	var job models.VideoGenerationJob
	err := database.DB.WithContext(ctx).Where("job_id = ?", jobID).First(&job).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(utils.APIResponse{
				Success: false,
				Error:   "Job not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(utils.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(utils.APIResponse{
		Success: true,
		Data: fiber.Map{
			"job_id":           job.JobID,
			"status":           job.Status,
			"video_url":        job.VideoURL,
			"optimized_prompt": job.OptimizedPrompt,
			"error_message":    job.ErrorMessage,
		},
	})
}

func (h *VideoHandler) runCinematicWorkflow(jobID string, input services.UIInput) {
	ctx, cancel := context.WithTimeout(context.Background(), cinematicWorkflowTimeout)
	defer cancel()

	optimizedPrompt, err := h.cinematicPromptService.GenerateCinematicPrompt(ctx, input)
	if err != nil {
		h.failJob(jobID, err)
		return
	}

	if err := database.DB.WithContext(ctx).
		Model(&models.VideoGenerationJob{}).
		Where("job_id = ?", jobID).
		Updates(map[string]interface{}{
			"status":           VideoStatusPromptGenerated,
			"optimized_prompt": optimizedPrompt,
			"error_message":    "",
		}).Error; err != nil {
		h.failJob(jobID, err)
		return
	}

	taskID, err := services.CallVeoVideoService(optimizedPrompt)
	if err != nil {
		h.failJob(jobID, err)
		return
	}

	if err := database.DB.WithContext(ctx).
		Model(&models.VideoGenerationJob{}).
		Where("job_id = ?", jobID).
		Updates(map[string]interface{}{
			"status":        VideoStatusProcessingVideo,
			"video_task_id": taskID,
			"error_message": "",
		}).Error; err != nil {
		h.failJob(jobID, err)
	}
}

func (h *VideoHandler) failJob(jobID string, err error) {
	errorMessage := err.Error()
	log.Printf("[CinematicWorkflow] job %s failed: %s", jobID, errorMessage)

	ctx, cancel := context.WithTimeout(context.Background(), cinematicStatusQueryTimeout)
	defer cancel()

	if updateErr := database.DB.WithContext(ctx).
		Model(&models.VideoGenerationJob{}).
		Where("job_id = ?", jobID).
		Updates(map[string]interface{}{
			"status":        VideoStatusFailed,
			"error_message": errorMessage,
		}).Error; updateErr != nil {
		log.Printf("[CinematicWorkflow] failed to persist failure for job %s: %v", jobID, updateErr)
	}
}
