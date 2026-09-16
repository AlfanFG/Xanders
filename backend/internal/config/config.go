package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	DatabaseURL            string
	GoogleAPIKey           string
	GoogleCloudProject     string
	GoogleCloudLocation    string
	GoogleGenaiUseVertexAI string
	GoogleCloudGCSBucket   string
	JWTSecret              string
}

func LoadConfig() *Config {
	// Try to load .env file, ignore if it doesn't exist
	_ = godotenv.Load()

	cfg := &Config{
		Port:                   getEnv("PORT", "8080"),
		DatabaseURL:            getEnv("DATABASE_URL", ""),
		GoogleAPIKey:           getEnv("GOOGLE_API_KEY", ""),
		GoogleCloudProject:     getEnv("GOOGLE_CLOUD_PROJECT", ""),
		GoogleCloudLocation:    getEnv("GOOGLE_CLOUD_LOCATION", "asia-southeast1"),
		GoogleGenaiUseVertexAI: getEnv("GOOGLE_GENAI_USE_VERTEXAI", "true"),
		GoogleCloudGCSBucket:   getEnv("GOOGLE_CLOUD_GCS_BUCKET", "gs://xanders-gen-v1/gen-video/"),
		JWTSecret:              getEnv("JWT_SECRET", "super-secret-key-change-me"),
	}

	if cfg.DatabaseURL == "" {
		log.Println("WARNING: DATABASE_URL is not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
