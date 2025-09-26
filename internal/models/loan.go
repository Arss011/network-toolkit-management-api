package models

import (
	"time"
)

type Loan struct {
	ID               int        `json:"id" gorm:"primaryKey"`
	UserID           int        `json:"user_id" gorm:"not null"`
	ToolkitID        int        `json:"toolkit_id" gorm:"not null"`
	Quantity         int        `json:"quantity" gorm:"default:1"`
	Purpose          string     `json:"purpose" binding:"required"`
	BorrowDate       time.Time  `json:"borrow_date"`
	DueDate          time.Time  `json:"due_date" gorm:"not null"`
	ReturnDate       *time.Time `json:"return_date"`
	Status           string     `json:"status" binding:"required" gorm:"default:borrowed"`
	ApprovedBy       string     `json:"approved_by"`
	Notes            string     `json:"notes"`
	ConditionChecked string     `json:"condition_checked"`
	ConditionReturn  string     `json:"condition_return"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Toolkit Toolkit `json:"toolkit,omitempty" gorm:"foreignKey:ToolkitID"`
}

type LoanFilterRequest struct {
	UserID     int        `json:"user_id,omitempty"`
	ToolkitID  int        `json:"toolkit_id,omitempty"`
	Status     string     `json:"status,omitempty"`
	DateFrom   *time.Time `json:"date_from,omitempty"`
	DateTo     *time.Time `json:"date_to,omitempty"`
	Overdue    bool       `json:"overdue,omitempty"`
	SearchTerm string     `json:"search_term,omitempty"`
}

type LoanCreateRequest struct {
	UserID           int       `json:"user_id" binding:"required"`
	ToolkitID        int       `json:"toolkit_id" binding:"required"`
	Quantity         int       `json:"quantity" binding:"required,min=1"`
	Purpose          string    `json:"purpose" binding:"required"`
	BorrowDate       time.Time `json:"borrow_date"`
	DueDate          time.Time `json:"due_date" binding:"required"`
	ApprovedBy       string    `json:"approved_by"`
	Notes            string    `json:"notes"`
	ConditionChecked string    `json:"condition_checked"`
}

type LoanUpdateRequest struct {
	UserID           int        `json:"user_id,omitempty"`
	ToolkitID        int        `json:"toolkit_id,omitempty"`
	Quantity         int        `json:"quantity,omitempty"`
	Purpose          string     `json:"purpose,omitempty"`
	BorrowDate       time.Time  `json:"borrow_date,omitempty"`
	DueDate          time.Time  `json:"due_date,omitempty"`
	ReturnDate       *time.Time `json:"return_date,omitempty"`
	Status           string     `json:"status,omitempty" binding:"omitempty,oneof=borrowed returned overdue damaged"`
	ApprovedBy       string     `json:"approved_by,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	ConditionChecked string     `json:"condition_checked,omitempty"`
	ConditionReturn  string     `json:"condition_return,omitempty"`
}
