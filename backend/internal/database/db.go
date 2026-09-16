package database

import (
	"log"
	"xanders-gen-video/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes the connection pool to PostgreSQL via GORM
func Connect(databaseURL string) {
	if databaseURL == "" {
		log.Println("DATABASE_URL is empty, skipping DB connection")
		return
	}

	var err error
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.CreditHistory{},
		&models.Job{},
		&models.VideoGenerationJob{},
		&models.PromptType{},
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database schema: %v", err)
	}

	log.Println("✅ Connected to PostgreSQL database via GORM")
}

// Close gracefully closes the underlying sql.DB connection
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
			log.Println("Database connection closed")
		}
	}
}
