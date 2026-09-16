package services

import (
	"context"
	"errors"
	"time"

	"xanders-gen-video/internal/config"
	"xanders-gen-video/internal/database"
	"xanders-gen-video/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) RegisterUser(ctx context.Context, email, password string) (*models.User, string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := models.User{
		Email:         email,
		PasswordHash:  string(hashedPassword),
		Plan:          "free",
		CreditBalance: SignupCreditAllowance,
	}

	if err := database.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		return tx.Create(&models.CreditHistory{
			UserID: user.ID,
			Amount: SignupCreditAllowance,
			Action: "signup_bonus",
		}).Error
	}); err != nil {
		return nil, "", errors.New("email already exists")
	}

	token, err := s.generateToken(&user)
	return &user, token, err
}

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (*models.User, string, error) {
	var user models.User
	result := database.DB.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	token, err := s.generateToken(&user)
	return &user, token, err
}

func (s *AuthService) generateToken(user *models.User) (string, error) {
	cfg := config.LoadConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(cfg.JWTSecret))
}
