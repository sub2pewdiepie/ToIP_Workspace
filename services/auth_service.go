package services

import (
	"errors"
	"space/auth"
	"space/models"
	"space/repositories"
	"space/utils"
)

type AuthService struct {
	UserRepo repositories.UserRepository
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) RegisterUser(input RegisterInput) error {
	_, err := s.UserRepo.GetByUsernameOrEmail(input.Username, input.Email)
	if err == nil {
		return errors.New("user already exists")
	}

	hashed, err := utils.HashPassword(input.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		HashPassword: hashed,
	}

	return s.UserRepo.Create(&user)
}

// services/auth_service.go
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (s *AuthService) LoginUser(input LoginInput) (string, error) {
	user, err := s.UserRepo.GetByUsername(input.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if err := utils.CheckPasswordHash(user.HashPassword, input.Password); err != nil {
		return "", errors.New("invalid username or password")
	}

	token, err := auth.GenerateJWT(user.Username)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
