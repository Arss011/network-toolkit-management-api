package services

import (
	"errors"
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
	repo        repositories.LoanRepository
	toolkitRepo repositories.ToolkitRepository
}

func NewLoanService(repo repositories.LoanRepository, toolkitRepo repositories.ToolkitRepository) LoanService {
	return &loanService{repo: repo, toolkitRepo: toolkitRepo}
}

func (s *loanService) Create(req *models.LoanCreateRequest) (*models.Loan, error) {
	// Check toolkit availability
	toolkit, err := s.toolkitRepo.GetByID(req.ToolkitID)
	if err != nil {
		return nil, errors.New("toolkit not found")
	}

	if toolkit.Available < req.Quantity {
		return nil, errors.New("insufficient toolkit quantity available")
	}

	// Create loan
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

	// Create loan
	createdLoan, err := s.repo.Create(loan)
	if err != nil {
		return nil, err
	}

	// Update toolkit availability
	toolkit.Available -= req.Quantity
	if toolkit.Available == 0 {
		toolkit.Status = "borrowed"
	}
	_, err = s.toolkitRepo.Update(toolkit)
	if err != nil {
		//delete loan if toolkit update fails
		_ = s.repo.Delete(createdLoan.ID)
		return nil, errors.New("failed to update toolkit availability")
	}

	return createdLoan, nil
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

	toolkit, err := s.toolkitRepo.GetByID(loan.ToolkitID)
	if err != nil {
		return nil, errors.New("toolkit not found")
	}

	// Handle status changes and qty updates
	oldStatus := loan.Status
	oldQuantity := loan.Quantity

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

	// Handle qty and availability updates
	if oldStatus != "returned" && loan.Status == "returned" {
		// Item is being returned
		toolkit.Available += loan.Quantity
		if toolkit.Available == toolkit.Quantity {
			toolkit.Status = "available"
		} else if toolkit.Available > 0 {
			toolkit.Status = "available"
		}
	} else if oldStatus == "returned" && loan.Status != "returned" {
		if toolkit.Available < loan.Quantity {
			return nil, errors.New("insufficient toolkit quantity available")
		}
		toolkit.Available -= loan.Quantity
		if toolkit.Available == 0 {
			toolkit.Status = "borrowed"
		}
	} else if oldQuantity != loan.Quantity && oldStatus != "returned" {
		// update qty if borowed
		quantityDiff := loan.Quantity - oldQuantity
		if quantityDiff > 0 && toolkit.Available < quantityDiff {
			return nil, errors.New("insufficient toolkit quantity available")
		}
		toolkit.Available -= quantityDiff
		if toolkit.Available == 0 {
			toolkit.Status = "borrowed"
		} else if toolkit.Available > 0 {
			toolkit.Status = "available"
		}
	}

	// Update toolkit
	_, err = s.toolkitRepo.Update(toolkit)
	if err != nil {
		return nil, errors.New("failed to update toolkit availability")
	}

	// Update loan
	return s.repo.Update(loan)
}

func (s *loanService) Delete(id int) error {
	return s.repo.Delete(id)
}
