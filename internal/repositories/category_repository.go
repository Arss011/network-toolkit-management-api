package repositories

import (
	"gorm.io/gorm"

	"toolkit-management/internal/models"
)

type CategoryRepository interface {
	Create(category *models.Category) (*models.Category, error)
	GetByID(id int) (*models.Category, error)
	GetAll(filter *models.CategoryFilterRequest) ([]models.Category, error)
	Update(category *models.Category) (*models.Category, error)
	Delete(id int) error
	GetTree() ([]models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *models.Category) (*models.Category, error) {
	result := r.db.Create(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (r *categoryRepository) GetByID(id int) (*models.Category, error) {
	var category models.Category
	result := r.db.First(&category, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &category, nil
}

func (r *categoryRepository) GetAll(filter *models.CategoryFilterRequest) ([]models.Category, error) {
	var categories []models.Category
	query := r.db.Model(&models.Category{})

	if filter.SearchTerm != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%")
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	result := query.Order("sort_order ASC, name ASC").Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}

func (r *categoryRepository) Update(category *models.Category) (*models.Category, error) {
	result := r.db.Save(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

func (r *categoryRepository) Delete(id int) error {
	result := r.db.Delete(&models.Category{}, id)
	return result.Error
}

func (r *categoryRepository) GetTree() ([]models.Category, error) {
	var categories []models.Category

	result := r.db.Order("sort_order ASC, name ASC").Find(&categories)

	if result.Error != nil {
		return nil, result.Error
	}

	return categories, nil
}