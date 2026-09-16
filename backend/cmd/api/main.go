package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"xanders-gen-video/internal/config"
	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Connect to database
	database.Connect(cfg.DatabaseURL)
	defer database.Close()

	// 3. Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Xanders AI Video API",
	})

	// 4. Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,https://xanders-web.onrender.com",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// 5. Routes
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "message": "API is running"})
	})

	router.SetupRoutes(app)

	// 6. Graceful shutdown setup
	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}
