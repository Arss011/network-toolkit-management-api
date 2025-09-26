package repositories

import (
	"gorm.io/gorm"

	"toolkit-management/internal/models"
	"toolkit-management/pkg/utils"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetAll(filter *models.UserFilterRequest) (*models.UserListResponse, error)
	Update(user *models.User) (*models.User, error)
	Delete(id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	result := r.db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *userRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) GetAll(filter *models.UserFilterRequest) (*models.UserListResponse, error) {
	var users []models.User
	var totalItems int64

	// Build query with filters
	query := r.db.Model(&models.User{})

	if filter.SearchTerm != "" {
		query = query.Where("username ILIKE ? OR email ILIKE ? OR full_name ILIKE ?",
			"%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%", "%"+filter.SearchTerm+"%")
	}

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	if filter.Department != "" {
		query = query.Where("department = ?", filter.Department)
	}

	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	// Get total count
	countResult := query.Count(&totalItems)
	if countResult.Error != nil {
		return nil, countResult.Error
	}

	// Apply pagination using GORM scope
	result := query.Scopes(utils.Paginate(filter.Page, filter.PageSize)).
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	// Calculate pagination response
	paginationResponse := utils.CalculatePagination(filter.Page, filter.PageSize, totalItems)

	return &models.UserListResponse{
		Data:       users,
		Pagination: paginationResponse,
	}, nil
}

func (r *userRepository) Update(user *models.User) (*models.User, error) {
	result := r.db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (r *userRepository) Delete(id int) error {
	result := r.db.Delete(&models.User{}, id)
	return result.Error
}