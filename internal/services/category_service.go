package services

import (
	"toolkit-management/internal/models"
	. "toolkit-management/internal/repositories"
)

type CategoryService interface {
	Create(req *models.CategoryCreateRequest) (*models.Category, error)
	GetByID(id int) (*models.Category, error)
	GetAll(filter *models.CategoryFilterRequest) ([]models.Category, error)
	Update(id int, req *models.CategoryUpdateRequest) (*models.Category, error)
	Delete(id int) error
	GetTree() ([]models.Category, error)
}

type categoryService struct {
	categoryRepo CategoryRepository
}

func NewCategoryService(repo CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: repo}
}

func (s *categoryService) Create(req *models.CategoryCreateRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		ParentID:    req.ParentID,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}

	return s.categoryRepo.Create(category)
}

func (s *categoryService) GetByID(id int) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *categoryService) GetAll(filter *models.CategoryFilterRequest) ([]models.Category, error) {
	return s.categoryRepo.GetAll(filter)
}

func (s *categoryService) Update(id int, req *models.CategoryUpdateRequest) (*models.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Color != "" {
		category.Color = req.Color
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.SortOrder != 0 {
		category.SortOrder = req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	return s.categoryRepo.Update(category)
}

func (s *categoryService) Delete(id int) error {
	return s.categoryRepo.Delete(id)
}

func (s *categoryService) GetTree() ([]models.Category, error) {
	return s.categoryRepo.GetTree()
}