package services

import (
	"time"
	"toolkit-management/internal/models"
	. "toolkit-management/internal/repositories"
	"toolkit-management/pkg/auth"
)

type UserService interface {
	Create(req *models.UserCreateRequest) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetAll(filter *models.UserFilterRequest) (*models.UserListResponse, error)
	Update(id int, req *models.UserUpdateRequest) (*models.User, error)
	Delete(id int) error
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
}

type userService struct {
	userRepo UserRepository
	authSvc  *auth.AuthService
}

func NewUserService(repo UserRepository) UserService {
	authConfig := auth.AuthConfig{
		SecretKey:     "your-secret-key-change-in-production",
		TokenDuration: 24 * time.Hour,
	}
	return &userService{
		userRepo: repo,
		authSvc:  auth.NewAuthService(authConfig),
	}
}

func (s *userService) Create(req *models.UserCreateRequest) (*models.User, error) {
	// Hash password
	hashedPassword, err := models.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		FullName:    req.FullName,
		Password:    hashedPassword,
		Role:        req.Role,
		Department:  req.Department,
		PhoneNumber: req.PhoneNumber,
		IsActive:    true,
	}

	return s.userRepo.Create(user)
}

func (s *userService) GetByID(id int) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetAll(filter *models.UserFilterRequest) (*models.UserListResponse, error) {
	return s.userRepo.GetAll(filter)
}

func (s *userService) Update(id int, req *models.UserUpdateRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Password != "" {
		user.Password = req.Password
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Department != "" {
		user.Department = req.Department
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	return s.userRepo.Update(user)
}

func (s *userService) Delete(id int) error {
	return s.userRepo.Delete(id)
}

func (s *userService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	// Check password
	if err := models.CheckPassword(user.Password, req.Password); err != nil {
		return nil, err
	}

	// Generate JWT token
	token, err := s.authSvc.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token:     token,
		User:      *user,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}, nil
}