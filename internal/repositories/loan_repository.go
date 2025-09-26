package repositories

import (
	"gorm.io/gorm"
	"time"

	"toolkit-management/internal/models"
)

type LoanRepository interface {
	Create(loan *models.Loan) (*models.Loan, error)
	GetByID(id int) (*models.Loan, error)
	GetAll(filter *models.LoanFilterRequest) ([]*models.Loan, error)
	Update(loan *models.Loan) (*models.Loan, error)
	Delete(id int) error
}

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db: db}
}

func (r *loanRepository) Create(loan *models.Loan) (*models.Loan, error) {
	result := r.db.Create(loan)
	if result.Error != nil {
		return nil, result.Error
	}
	return loan, nil
}

func (r *loanRepository) GetByID(id int) (*models.Loan, error) {
	var loan models.Loan
	result := r.db.First(&loan, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &loan, nil
}

func (r *loanRepository) GetAll(filter *models.LoanFilterRequest) ([]*models.Loan, error) {
	var loans []*models.Loan
	query := r.db.Model(&models.Loan{})

	if filter.UserID != 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}

	if filter.ToolkitID != 0 {
		query = query.Where("toolkit_id = ?", filter.ToolkitID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.DateFrom != nil {
		query = query.Where("borrow_date >= ?", *filter.DateFrom)
	}

	if filter.DateTo != nil {
		query = query.Where("borrow_date <= ?", *filter.DateTo)
	}

	if filter.Overdue {
		query = query.Where("due_date < ? AND status != 'returned'", time.Now())
	}

	if filter.SearchTerm != "" {
		query = query.Joins("JOIN users ON loans.user_id = users.id").
			Joins("JOIN toolkits ON loans.toolkit_id = toolkits.id").
			Where("users.username ILIKE ? OR users.full_name ILIKE ? OR toolkits.name ILIKE ? OR loans.purpose ILIKE ?",
				"%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%")
	}

	result := query.Preload("User").Preload("Toolkit").Order("created_at DESC").Find(&loans)
	if result.Error != nil {
		return nil, result.Error
	}

	return loans, nil
}

func (r *loanRepository) Update(loan *models.Loan) (*models.Loan, error) {
	result := r.db.Save(loan)
	if result.Error != nil {
		return nil, result.Error
	}
	return loan, nil
}

func (r *loanRepository) Delete(id int) error {
	result := r.db.Delete(&models.Loan{}, id)
	return result.Error
}