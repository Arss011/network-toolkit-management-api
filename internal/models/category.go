package models

import (
	"time"
)

type Category struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" binding:"required" gorm:"unique;not null"`
	Description string     `json:"description"`
	Color       string     `json:"color" gorm:"default:#6B7280"`
	Icon        string     `json:"icon"`
	ParentID    *int       `json:"parent_id"`
	SortOrder   int        `json:"sort_order" gorm:"default:0"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Toolkits []Toolkit  `json:"toolkits,omitempty" gorm:"foreignKey:CategoryID"`
}

type CategoryFilterRequest struct {
	SearchTerm string `json:"search_term,omitempty"`
	ParentID   *int   `json:"parent_id,omitempty"`
	IsActive   *bool  `json:"is_active,omitempty"`
}

type CategoryCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color" binding:"required"`
	Icon        string `json:"icon"`
	ParentID    *int   `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
}

type CategoryUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
	Icon        string `json:"icon,omitempty"`
	ParentID    *int   `json:"parent_id,omitempty"`
	SortOrder   int    `json:"sort_order,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"`
}
