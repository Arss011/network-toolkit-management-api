package services

import (
	"toolkit-management/internal/models"
	. "toolkit-management/internal/repositories"
)

type ToolkitService interface {
	Create(req *models.ToolkitCreateRequest) (*models.Toolkit, error)
	GetByID(id int) (*models.Toolkit, error)
	GetAll(filter *models.ToolkitFilterRequest) (*models.ToolkitListResponse, error)
	Update(id int, req *models.ToolkitUpdateRequest) (*models.Toolkit, error)
	Delete(id int) error
	UpdateStock(id int, req *models.ToolkitStockUpdateRequest) (*models.Toolkit, error)
}

type toolkitService struct {
	toolkitRepo ToolkitRepository
}

func NewToolkitService(repo ToolkitRepository) ToolkitService {
	return &toolkitService{toolkitRepo: repo}
}

func (s *toolkitService) Create(req *models.ToolkitCreateRequest) (*models.Toolkit, error) {
	toolkit := &models.Toolkit{
		Name:          req.Name,
		SKU:           req.SKU,
		Description:   req.Description,
		CategoryID:    req.CategoryID,
		Quantity:      req.Quantity,
		Available:     req.Quantity,
		Unit:          req.Unit,
		Brand:         req.Brand,
		Model:         req.Model,
		SerialNumber:  req.SerialNumber,
		PurchaseDate:  req.PurchaseDate,
		PurchasePrice: req.PurchasePrice,
		Condition:     req.Condition,
		Status:        "available",
		ImageURL:      req.ImageURL,
		Notes:         req.Notes,
	}

	return s.toolkitRepo.Create(toolkit)
}

func (s *toolkitService) GetByID(id int) (*models.Toolkit, error) {
	return s.toolkitRepo.GetByID(id)
}

func (s *toolkitService) GetAll(filter *models.ToolkitFilterRequest) (*models.ToolkitListResponse, error) {
	return s.toolkitRepo.GetAll(filter)
}

func (s *toolkitService) Update(id int, req *models.ToolkitUpdateRequest) (*models.Toolkit, error) {
	toolkit, err := s.toolkitRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		toolkit.Name = req.Name
	}
	if req.SKU != "" {
		toolkit.SKU = req.SKU
	}
	if req.Description != "" {
		toolkit.Description = req.Description
	}
	if req.CategoryID != 0 {
		toolkit.CategoryID = req.CategoryID
	}
	if req.Quantity != 0 {
		toolkit.Quantity = req.Quantity
	}
	if req.Unit != "" {
		toolkit.Unit = req.Unit
	}
	if req.Brand != "" {
		toolkit.Brand = req.Brand
	}
	if req.Model != "" {
		toolkit.Model = req.Model
	}
	if req.SerialNumber != "" {
		toolkit.SerialNumber = req.SerialNumber
	}
	if req.PurchaseDate != nil {
		toolkit.PurchaseDate = req.PurchaseDate
	}
	if req.PurchasePrice != 0 {
		toolkit.PurchasePrice = req.PurchasePrice
	}
	if req.Condition != "" {
		toolkit.Condition = req.Condition
	}
	if req.Status != "" {
		toolkit.Status = req.Status
	}
	if req.ImageURL != "" {
		toolkit.ImageURL = req.ImageURL
	}
	if req.Notes != "" {
		toolkit.Notes = req.Notes
	}

	return s.toolkitRepo.Update(toolkit)
}

func (s *toolkitService) Delete(id int) error {
	return s.toolkitRepo.Delete(id)
}

func (s *toolkitService) UpdateStock(id int, req *models.ToolkitStockUpdateRequest) (*models.Toolkit, error) {
	toolkit, err := s.toolkitRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	toolkit.Quantity += req.QuantityChange
	toolkit.Available += req.QuantityChange

	if toolkit.Quantity < 0 {
		toolkit.Quantity = 0
	}
	if toolkit.Available < 0 {
		toolkit.Available = 0
	}

	return s.toolkitRepo.Update(toolkit)
}