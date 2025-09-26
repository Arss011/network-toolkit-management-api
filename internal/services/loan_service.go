package services

import (
	"time"

	"toolkit-management/internal/models"
	"toolkit-management/internal/repositories"
)

type LoanService interface {
	Create(req *models.LoanCreateRequest) (*models.Loan, error)
	GetByID(id int) (*models.Loan, error)
	GetAll(filter *models.LoanFilterRequest) ([]*models.Loan, error)
	Update(id int, req *models.LoanUpdateRequest) (*models.Loan, error)
	Delete(id int) error
}

type loanService struct {
	repo repositories.LoanRepository
}

func NewLoanService(repo repositories.LoanRepository) LoanService {
	return &loanService{repo: repo}
}

func (s *loanService) Create(req *models.LoanCreateRequest) (*models.Loan, error) {
	loan := &models.Loan{
		UserID:           req.UserID,
		ToolkitID:        req.ToolkitID,
		Quantity:         req.Quantity,
		Purpose:          req.Purpose,
		BorrowDate:       time.Now(),
		DueDate:          req.DueDate,
		Status:           "borrowed",
		ApprovedBy:       req.ApprovedBy,
		Notes:            req.Notes,
		ConditionChecked: req.ConditionChecked,
	}

	return s.repo.Create(loan)
}

func (s *loanService) GetByID(id int) (*models.Loan, error) {
	return s.repo.GetByID(id)
}

func (s *loanService) GetAll(filter *models.LoanFilterRequest) ([]*models.Loan, error) {
	return s.repo.GetAll(filter)
}

func (s *loanService) Update(id int, req *models.LoanUpdateRequest) (*models.Loan, error) {
	loan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.UserID != 0 {
		loan.UserID = req.UserID
	}
	if req.ToolkitID != 0 {
		loan.ToolkitID = req.ToolkitID
	}
	if req.Quantity != 0 {
		loan.Quantity = req.Quantity
	}
	if req.Purpose != "" {
		loan.Purpose = req.Purpose
	}
	if !req.BorrowDate.IsZero() {
		loan.BorrowDate = req.BorrowDate
	}
	if !req.DueDate.IsZero() {
		loan.DueDate = req.DueDate
	}
	if req.ReturnDate != nil {
		loan.ReturnDate = req.ReturnDate
	}
	if req.Status != "" {
		loan.Status = req.Status
	}
	if req.ApprovedBy != "" {
		loan.ApprovedBy = req.ApprovedBy
	}
	if req.Notes != "" {
		loan.Notes = req.Notes
	}
	if req.ConditionChecked != "" {
		loan.ConditionChecked = req.ConditionChecked
	}
	if req.ConditionReturn != "" {
		loan.ConditionReturn = req.ConditionReturn
	}

	return s.repo.Update(loan)
}

func (s *loanService) Delete(id int) error {
	return s.repo.Delete(id)
}