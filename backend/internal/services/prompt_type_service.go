package services

import (
	"context"
	"errors"
	"math"
	"strings"

	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/models"

	"gorm.io/gorm"
)

var allowedPromptTypeCategories = map[string]bool{
	"template": true,
	"audio":    true,
	"lighting": true,
	"camera":   true,
}

type PromptTypeService struct{}

type PromptTypeListRequest struct {
	Page           int
	Limit          int
	Type           string
	Search         string
	IncludeDeleted bool
}

type PromptTypeListResult struct {
	Items      []models.PromptType `json:"items"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	Total      int64               `json:"total"`
	TotalPages int                 `json:"total_pages"`
}

type PromptTypePayload struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func NewPromptTypeService() *PromptTypeService {
	return &PromptTypeService{}
}

func (s *PromptTypeService) List(ctx context.Context, req PromptTypeListRequest) (*PromptTypeListResult, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}

	limit := req.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	query := database.DB.WithContext(ctx).Model(&models.PromptType{})
	if !req.IncludeDeleted {
		query = query.Where("is_deleted = ?", false)
	}

	category := strings.TrimSpace(strings.ToLower(req.Type))
	if category != "" {
		query = query.Where("type = ?", category)
	}

	search := strings.TrimSpace(req.Search)
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []models.PromptType
	if err := query.
		Order("type ASC, name ASC").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&items).Error; err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.PromptType{}
	}

	return &PromptTypeListResult{
		Items:      items,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}, nil
}

func (s *PromptTypeService) GetByID(ctx context.Context, id string) (*models.PromptType, error) {
	var item models.PromptType
	if err := database.DB.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *PromptTypeService) Create(ctx context.Context, payload PromptTypePayload) (*models.PromptType, error) {
	name, category, description, err := normalizePromptTypePayload(payload)
	if err != nil {
		return nil, err
	}

	item := models.PromptType{
		Name:        name,
		Type:        category,
		Description: description,
	}

	if err := database.DB.WithContext(ctx).Create(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *PromptTypeService) Update(ctx context.Context, id string, payload PromptTypePayload) (*models.PromptType, error) {
	name, category, description, err := normalizePromptTypePayload(payload)
	if err != nil {
		return nil, err
	}

	var item models.PromptType
	err = database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&item).Error; err != nil {
			return err
		}

		item.Name = name
		item.Type = category
		item.Description = description
		return tx.Save(&item).Error
	})
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *PromptTypeService) Delete(ctx context.Context, id string) error {
	result := database.DB.WithContext(ctx).
		Model(&models.PromptType{}).
		Where("id = ? AND is_deleted = ?", id, false).
		Update("is_deleted", true)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func normalizePromptTypePayload(payload PromptTypePayload) (string, string, string, error) {
	name := strings.TrimSpace(payload.Name)
	category := strings.TrimSpace(strings.ToLower(payload.Type))
	description := strings.TrimSpace(payload.Description)

	if name == "" {
		return "", "", "", errors.New("name is required")
	}
	if category == "" {
		category = "template"
	}
	if !allowedPromptTypeCategories[category] {
		return "", "", "", errors.New("type must be one of: template, audio, lighting, camera")
	}

	return name, category, description, nil
}
