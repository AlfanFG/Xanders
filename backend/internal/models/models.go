package models

import (
	"encoding/json"
	"time"
)

type User struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email         string    `json:"email" gorm:"unique;not null"`
	PasswordHash  string    `json:"-" gorm:"not null"`
	Plan          string    `json:"plan" gorm:"default:'free'"`
	CreditBalance int       `json:"credit_balance" gorm:"default:0"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type CreditHistory struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null"`
	Amount    int       `json:"amount" gorm:"not null"`
	Action    string    `json:"action" gorm:"column:reason;not null"` // "signup_bonus", "generation_cost", "purchase"
	JobID     *string   `json:"job_id,omitempty" gorm:"type:uuid"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (CreditHistory) TableName() string {
	return "credit_transactions"
}

type Job struct {
	ID              string          `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID          string          `json:"user_id" gorm:"type:uuid;not null"`
	Status          string          `json:"status" gorm:"default:'pending'"` // "pending", "processing", "completed", "failed"
	PromptInput     json.RawMessage `json:"prompt_input" gorm:"type:jsonb"`
	AssembledPrompt string          `json:"assembled_prompt" gorm:"type:text"`
	ImageURL        *string         `json:"image_url,omitempty" gorm:"type:text"` // For image-to-video visual prompting
	ModelName       string          `json:"model_name" gorm:"default:'veo-3.0-fast-generate-001'"`
	DurationSeconds int             `json:"duration_seconds" gorm:"default:5"`
	CreditsCharged  int             `json:"credits_charged" gorm:"not null"`
	GCSVideoURL     *string         `json:"gcs_video_url,omitempty" gorm:"type:text"`
	ErrorMessage    *string         `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt       time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Job) TableName() string {
	return "video_generation_jobs"
}

type VideoGenerationJob struct {
	JobID           string    `json:"job_id" gorm:"primaryKey;type:uuid;uniqueIndex"`
	UserID          uint      `json:"user_id,omitempty"`
	UIProduct       string    `json:"ui_product" gorm:"not null"`
	UIEnvironment   string    `json:"ui_environment" gorm:"not null"`
	UIStyle         string    `json:"ui_style" gorm:"not null"`
	UIMotion        string    `json:"ui_motion" gorm:"not null"`
	OptimizedPrompt string    `json:"optimized_prompt,omitempty" gorm:"type:text"`
	VideoTaskID     string    `json:"video_task_id,omitempty"`
	VideoURL        string    `json:"video_url,omitempty"`
	Status          string    `json:"status" gorm:"not null;default:'PENDING'"`
	ErrorMessage    string    `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type PromptType struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:text;not null"`
	Type        string    `json:"type" gorm:"type:text;not null;default:'template';index"`
	Description     string    `json:"description" gorm:"type:text;not null;default:''"`
	Icon            string    `json:"icon" gorm:"type:text;default:''"`
	BackgroundClass string    `json:"background_class" gorm:"type:text;default:''"`
	IsDeleted       bool      `json:"is_deleted" gorm:"not null;default:false;index"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PromptType) TableName() string {
	return "prompt_type"
}

func (VideoGenerationJob) TableName() string {
	return "cinematic_video_generation_jobs"
}
