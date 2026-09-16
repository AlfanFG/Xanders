package router

import (
	"xanders-gen-video/internal/config"
	"xanders-gen-video/internal/handlers"
	"xanders-gen-video/internal/middleware"
	"xanders-gen-video/internal/services"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	// Initialize Config
	cfg := config.LoadConfig()

	// Initialize Services
	authService := services.NewAuthService()
	creditService := services.NewCreditService()
	promptService := services.NewPromptService()
	videoService := services.NewVideoService(cfg)
	queueService := services.NewQueueService(videoService)
	cinematicPromptService := services.NewCinematicPromptService()
	promptTypeService := services.NewPromptTypeService()

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler()
	videoHandler := handlers.NewVideoHandler(cinematicPromptService)
	jobHandler := handlers.NewJobHandler(queueService, creditService, promptService, cinematicPromptService)
	promptTypeHandler := handlers.NewPromptTypeHandler(promptTypeService)

	video := app.Group("/api/video")
	video.Post("/generate", videoHandler.Generate)
	video.Get("/status/:job_id", videoHandler.Status)

	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected routes
	me := api.Group("/me", middleware.Protected())
	me.Get("/", userHandler.GetProfile)
	me.Get("/credits", userHandler.GetCreditHistory)

	jobs := api.Group("/jobs", middleware.Protected())
	jobs.Post("/", jobHandler.CreateJob)
	jobs.Get("/", jobHandler.ListJobs)
	jobs.Get("/:id", jobHandler.GetJob)

	promptTypes := api.Group("/prompt-types", middleware.Protected())
	promptTypes.Get("/", promptTypeHandler.List)
	promptTypes.Get("/:id", promptTypeHandler.GetByID)
	promptTypes.Post("/", middleware.RequireAdmin(), promptTypeHandler.Create)
	promptTypes.Put("/:id", middleware.RequireAdmin(), promptTypeHandler.Update)
	promptTypes.Delete("/:id", middleware.RequireAdmin(), promptTypeHandler.Delete)
}
