package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	SignupCreditAllowance     = 15
	VideoGenerationCreditCost = 5
	CreditSupportEmail        = "alfanfaturahman10@gmail.com"
)

var ErrInsufficientCredits = errors.New("insufficient credits")

type CreditService struct{}

func NewCreditService() *CreditService {
	return &CreditService{}
}

// EnsureSufficientCredits avoids starting costly provider work when the user
// cannot afford a generation. The transactional deduction below remains the
// authoritative check because a balance can change concurrently.
func (s *CreditService) EnsureSufficientCredits(ctx context.Context, userID string, amount int) error {
	var user models.User
	if err := database.DB.WithContext(ctx).Select("credit_balance").Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	if user.CreditBalance < amount {
		return fmt.Errorf("%w: %d available, %d required", ErrInsufficientCredits, user.CreditBalance, amount)
	}

	return nil
}

// DeductCreditsAndCreateJob ensures credit deduction and job creation happen in one atomic transaction.
func (s *CreditService) DeductCreditsAndCreateJob(ctx context.Context, userID string, amount int, jobID string, promptInput map[string]interface{}, assembledPrompt string, imageURL *string) error {
	return database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var user models.User

		// Lock the row to prevent concurrent race conditions (double-spend)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", userID).First(&user).Error; err != nil {
			return errors.New("user not found")
		}

		if user.CreditBalance < amount {
			return fmt.Errorf("%w: %d available, %d required", ErrInsufficientCredits, user.CreditBalance, amount)
		}

		// 1. Deduct balance
		if err := tx.Model(&user).Update("credit_balance", gorm.Expr("credit_balance - ?", amount)).Error; err != nil {
			return err
		}

		// 2. Insert transaction log
		history := models.CreditHistory{
			UserID: userID,
			Amount: -amount,
			Action: "generation_cost",
			JobID:  &jobID,
		}
		if err := tx.Create(&history).Error; err != nil {
			return err
		}

		promptJSON, err := json.Marshal(promptInput)
		if err != nil {
			return err
		}

		// 3. Insert job record
		job := models.Job{
			ID:              jobID,
			UserID:          userID,
			Status:          "pending",
			PromptInput:     promptJSON,
			AssembledPrompt: assembledPrompt,
			CreditsCharged:  amount,
			ImageURL:        imageURL,
		}
		if err := tx.Create(&job).Error; err != nil {
			return err
		}

		return nil
	})
}
