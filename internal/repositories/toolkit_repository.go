package repositories

import (
	"gorm.io/gorm"

	"toolkit-management/internal/models"
	"toolkit-management/pkg/utils"
)

type ToolkitRepository interface {
	Create(toolkit *models.Toolkit) (*models.Toolkit, error)
	GetByID(id int) (*models.Toolkit, error)
	GetAll(filter *models.ToolkitFilterRequest) (*models.ToolkitListResponse, error)
	Update(toolkit *models.Toolkit) (*models.Toolkit, error)
	Delete(id int) error
}

type toolkitRepository struct {
	db *gorm.DB
}

func NewToolkitRepository(db *gorm.DB) ToolkitRepository {
	return &toolkitRepository{db: db}
}

func (r *toolkitRepository) Create(toolkit *models.Toolkit) (*models.Toolkit, error) {
	result := r.db.Create(toolkit)
	if result.Error != nil {
		return nil, result.Error
	}
	return toolkit, nil
}

func (r *toolkitRepository) GetByID(id int) (*models.Toolkit, error) {
	var toolkit models.Toolkit
	result := r.db.First(&toolkit, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &toolkit, nil
}

func (r *toolkitRepository) GetAll(filter *models.ToolkitFilterRequest) (*models.ToolkitListResponse, error) {
	var toolkits []models.Toolkit
	var totalItems int64

	// Build query with filters
	query := r.db.Model(&models.Toolkit{})

	if filter.SearchTerm != "" {
		query = query.Where("name ILIKE ? OR sku ILIKE ? OR description ILIKE ?",
			"%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%")
	}

	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Condition != "" {
		query = query.Where("condition = ?", filter.Condition)
	}

	if filter.Brand != "" {
		query = query.Where("brand = ?", filter.Brand)
	}

	if filter.MinQuantity > 0 {
		query = query.Where("quantity >= ?", filter.MinQuantity)
	}

	if filter.MaxQuantity > 0 {
		query = query.Where("quantity <= ?", filter.MaxQuantity)
	}

	// Get total count
	countResult := query.Count(&totalItems)
	if countResult.Error != nil {
		return nil, countResult.Error
	}

	// Apply pagination using GORM scope
	result := query.Scopes(utils.Paginate(filter.Page, filter.PageSize)).
		Preload("Category").
		Find(&toolkits)

	if result.Error != nil {
		return nil, result.Error
	}

	// Calculate pagination response
	paginationResponse := utils.CalculatePagination(filter.Page, filter.PageSize, totalItems)

	return &models.ToolkitListResponse{
		Data:       toolkits,
		Pagination: paginationResponse,
	}, nil
}

func (r *toolkitRepository) Update(toolkit *models.Toolkit) (*models.Toolkit, error) {
	result := r.db.Save(toolkit)
	if result.Error != nil {
		return nil, result.Error
	}
	return toolkit, nil
}

func (r *toolkitRepository) Delete(id int) error {
	result := r.db.Delete(&models.Toolkit{}, id)
	return result.Error
}