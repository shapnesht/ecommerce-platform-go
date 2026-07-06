package services

import (
	"github.com/shapnesht/ecommerce-platform-go/internal/dto"
	"github.com/shapnesht/ecommerce-platform-go/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetProfile(userId uint) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		IsActive:  user.IsActive,
		Role:      string(user.Role),
	}, nil
}

func (s *UserService) UpdateProfile(userId uint, req *dto.ProfileUpdateRequest) (*dto.UserResponse, error) {
	var user models.User
	if err := s.db.Where("user_id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return s.GetProfile(userId)
}
