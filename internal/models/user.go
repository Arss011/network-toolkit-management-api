package models

import (
	"time"
	"toolkit-management/pkg/utils"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	Username    string     `json:"username" binding:"required" gorm:"unique;not null"`
	Email       string     `json:"email" binding:"required,email" gorm:"unique;not null"`
	FullName    string     `json:"full_name" binding:"required"`
	Password    string     `json:"-" gorm:"not null"`
	Role        string     `json:"role" binding:"required" gorm:"default:user"`
	Department  string     `json:"department"`
	PhoneNumber string     `json:"phone_number"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	LastLogin   *time.Time `json:"last_login"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Loans       []Loan     `json:"loans,omitempty" gorm:"foreignKey:UserID"`
}

type UserFilterRequest struct {
	SearchTerm string `json:"search_term,omitempty"`
	Role       string `json:"role,omitempty"`
	Department string `json:"department,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
	Page       int    `json:"page,omitempty" form:"page"`
	PageSize   int    `json:"page_size,omitempty" form:"page_size"`
}

type UserListResponse struct {
	Data       []User                   `json:"data"`
	Pagination utils.PaginationResponse `json:"pagination"`
}

type UserCreateRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	FullName    string `json:"full_name" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Role        string `json:"role" binding:"required,oneof=admin user technician"`
	Department  string `json:"department"`
	PhoneNumber string `json:"phone_number"`
}

type UserUpdateRequest struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	FullName    string `json:"full_name,omitempty"`
	Password    string `json:"password,omitempty"`
	Role        string `json:"role,omitempty" binding:"omitempty,oneof=admin user technician"`
	Department  string `json:"department,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	User      User      `json:"user"`
	ExpiresAt time.Time `json:"expires_at"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
