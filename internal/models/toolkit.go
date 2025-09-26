package models

import (
	"time"
	"toolkit-management/pkg/utils"
)

type Toolkit struct {
	ID            int        `json:"id" gorm:"primaryKey"`
	Name          string     `json:"name" binding:"required" gorm:"not null"`
	SKU           string     `json:"sku" gorm:"unique;not null"`
	Description   string     `json:"description"`
	CategoryID    int        `json:"category_id" gorm:"not null"`
	Quantity      int        `json:"quantity" gorm:"not null"`
	Available     int        `json:"available" gorm:"not null"`
	Unit          string     `json:"unit" gorm:"default:unit"`
	Brand         string     `json:"brand"`
	Model         string     `json:"model"`
	SerialNumber  string     `json:"serial_number"`
	PurchaseDate  *time.Time `json:"purchase_date"`
	PurchasePrice float64    `json:"purchase_price"`
	Condition     string     `json:"condition" gorm:"default:good"`
	Status        string     `json:"status" gorm:"default:available"`
	ImageURL      string     `json:"image_url"`
	Notes         string     `json:"notes"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Loans    []Loan   `json:"loans,omitempty" gorm:"foreignKey:ToolkitID"`
}

type ToolkitFilterRequest struct {
	SearchTerm  string `json:"search_term,omitempty"`
	CategoryID  int    `json:"category_id,omitempty"`
	Status      string `json:"status,omitempty"`
	Condition   string `json:"condition,omitempty"`
	Brand       string `json:"brand,omitempty"`
	MinQuantity int    `json:"min_quantity,omitempty"`
	MaxQuantity int    `json:"max_quantity,omitempty"`
	Page        int    `json:"page,omitempty" form:"page"`
	PageSize    int    `json:"page_size,omitempty" form:"page_size"`
}

type ToolkitCreateRequest struct {
	Name          string     `json:"name" binding:"required"`
	SKU           string     `json:"sku" binding:"required"`
	Description   string     `json:"description"`
	CategoryID    int        `json:"category_id" binding:"required"`
	Quantity      int        `json:"quantity" binding:"required,min=1"`
	Unit          string     `json:"unit" binding:"required"`
	Brand         string     `json:"brand"`
	Model         string     `json:"model"`
	SerialNumber  string     `json:"serial_number"`
	PurchaseDate  *time.Time `json:"purchase_date"`
	PurchasePrice float64    `json:"purchase_price"`
	Condition     string     `json:"condition" binding:"required,oneof=excellent good fair poor"`
	ImageURL      string     `json:"image_url"`
	Notes         string     `json:"notes"`
}

type ToolkitUpdateRequest struct {
	Name          string     `json:"name,omitempty"`
	SKU           string     `json:"sku,omitempty"`
	Description   string     `json:"description,omitempty"`
	CategoryID    int        `json:"category_id,omitempty"`
	Quantity      int        `json:"quantity,omitempty"`
	Unit          string     `json:"unit,omitempty"`
	Brand         string     `json:"brand,omitempty"`
	Model         string     `json:"model,omitempty"`
	SerialNumber  string     `json:"serial_number,omitempty"`
	PurchaseDate  *time.Time `json:"purchase_date,omitempty"`
	PurchasePrice float64    `json:"purchase_price,omitempty"`
	Condition     string     `json:"condition,omitempty" binding:"omitempty,oneof=excellent good fair poor"`
	Status        string     `json:"status,omitempty" binding:"omitempty,oneof=available borrowed maintenance retired"`
	ImageURL      string     `json:"image_url,omitempty"`
	Notes         string     `json:"notes,omitempty"`
}

type ToolkitStockUpdateRequest struct {
	QuantityChange int    `json:"quantity_change" binding:"required"`
	Reason         string `json:"reason" binding:"required"`
	Notes          string `json:"notes"`
}

type ToolkitListResponse struct {
	Data       []Toolkit                `json:"data"`
	Pagination utils.PaginationResponse `json:"pagination"`
}
