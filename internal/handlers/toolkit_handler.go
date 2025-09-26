package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"toolkit-management/internal/models"
	"toolkit-management/internal/services"
)

type ToolkitHandler struct {
	service services.ToolkitService
}

func NewToolkitHandler(service services.ToolkitService) *ToolkitHandler {
	return &ToolkitHandler{service: service}
}

func (h *ToolkitHandler) Create(c *gin.Context) {
	var req models.ToolkitCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	result, err := h.service.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Toolkit created successfully",
		"data":    result,
	})
}

func (h *ToolkitHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID"})
		return
	}

	toolkit, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Toolkit not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Toolkit retrieved successfully",
		"data":    toolkit,
	})
}

func (h *ToolkitHandler) GetAll(c *gin.Context) {
	var filter models.ToolkitFilterRequest

	if err := c.ShouldBindQuery(&filter); err != nil {
		filter = models.ToolkitFilterRequest{}
	}

	var bodyFilter models.ToolkitFilterRequest
	if c.Request.Method == "POST" && c.Request.Header.Get("Content-Type") == "application/json" {
		if err := c.ShouldBindJSON(&bodyFilter); err == nil {
			filter.SearchTerm = bodyFilter.SearchTerm
			filter.CategoryID = bodyFilter.CategoryID
			filter.Status = bodyFilter.Status
			filter.Condition = bodyFilter.Condition
			filter.Brand = bodyFilter.Brand
			filter.MinQuantity = bodyFilter.MinQuantity
			filter.MaxQuantity = bodyFilter.MaxQuantity
		}
	}

	toolkitList, err := h.service.GetAll(&filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "Toolkits retrieved successfully",
		"data":       toolkitList.Data,
		"pagination": toolkitList.Pagination,
	})
}

func (h *ToolkitHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID"})
		return
	}

	var req models.ToolkitUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	result, err := h.service.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Toolkit updated successfully",
		"data":    result,
	})
}

func (h *ToolkitHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID"})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Toolkit deleted successfully",
	})
}

func (h *ToolkitHandler) UpdateStock(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid ID"})
		return
	}

	var req models.ToolkitStockUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	result, err := h.service.UpdateStock(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Stock updated successfully",
		"data":    result,
	})
}
