package service

import (
	"errors"
	"patungan-server/internal/config"
	"patungan-server/internal/dto/request"
	"patungan-server/internal/dto/response"
	"patungan-server/internal/models"
	"patungan-server/internal/repository"
	"patungan-server/internal/utils"
)

type AuthService interface {
	Register(req *request.RegisterRequest) (*response.AuthResponse, error)
	Login(req *request.LoginRequest) (*response.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg: cfg,
	}
}

func (s *authService) Register(req *request.RegisterRequest) (*response.AuthResponse, error) {
	// Checj if user exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email: req.Email,
		Password: hashedPassword,
		FullName: req.FullName,
		Phone: req.Phone,
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT token
	accessToken, _ := utils.GenerateToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.AccessExpiry)
	refreshToken, _ := utils.GenerateToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.RefreshExpiry)

	return &response.AuthResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: response.UserResponse{
			ID: user.ID,
			Email: user.Email,
			FullName: user.FullName,
			Phone: user.Phone,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

func (s *authService) Login(req *request.LoginRequest) (*response.AuthResponse, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	accessToken, _ := utils.GenerateToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.AccessExpiry)
	refreshToken, _ := utils.GenerateToken(user.ID, user.Email, s.cfg.JWT.Secret, s.cfg.JWT.RefreshExpiry)

	return &response.AuthResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		User: response.UserResponse{
			ID: user.ID,
			Email: user.Email,
			FullName: user.FullName,
			Phone: user.Phone,
			AvatarURL: user.AvatarURL,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}