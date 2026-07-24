package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/shapnesht/ecommerce-platform-go/internal/config"
	"github.com/shapnesht/ecommerce-platform-go/internal/dto"
	"github.com/shapnesht/ecommerce-platform-go/internal/events"
	"github.com/shapnesht/ecommerce-platform-go/internal/models"
	"github.com/shapnesht/ecommerce-platform-go/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db             *gorm.DB
	config         *config.Config
	eventPublisher events.EventPublisher
}

func NewAuthService(db *gorm.DB, config *config.Config, eventPublisher events.EventPublisher) *AuthService {
	return &AuthService{
		db:             db,
		config:         config,
		eventPublisher: eventPublisher,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// check if user exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists with same email")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// create a cart
	cart := models.Cart{
		UserID: user.ID,
	}

	if err := s.db.Create(&cart).Error; err != nil {
		fmt.Printf("Unable to create cart.")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	var user models.User
	if err := s.db.Where("user_id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	s.db.Delete(&refreshToken)

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(refreshToken string) error {
	return s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{}).Error
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateTokenPair(&s.config.JWT, user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpiry),
	}

	if err := s.db.Create(&refreshTokenModel).Error; err != nil {
		return nil, err
	}

	err = s.eventPublisher.Publish("USER_LOGGED_IN", user, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("unable to publish login message: %w", err)
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      string(user.Role),
			Phone:     user.Phone,
			IsActive:  user.IsActive,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
