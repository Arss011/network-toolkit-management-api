package models

import (
	"time"
)

type Category struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" binding:"required" gorm:"unique;not null"`
	Description string     `json:"description"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Toolkits []Toolkit  `json:"toolkits,omitempty" gorm:"foreignKey:CategoryID"`
}

type CategoryFilterRequest struct {
	SearchTerm string `json:"search_term,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

type CategoryCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
}

type CategoryUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}
