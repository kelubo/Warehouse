package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetWarehouses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.Warehouse{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR location LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var warehouses []models.Warehouse
	offset := (page - 1) * pageSize
	query.Limit(pageSize).Offset(offset).Find(&warehouses)

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, models.PaginationResponse{
		Items:      warehouses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func GetWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := database.DB.Where("id = ?", c.Param("id")).First(&warehouse).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}
	c.JSON(http.StatusOK, warehouse)
}

func CreateWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&warehouse)
	c.JSON(http.StatusCreated, warehouse)
}

func UpdateWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := database.DB.Where("id = ?", c.Param("id")).First(&warehouse).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&warehouse)
	c.JSON(http.StatusOK, warehouse)
}

func DeleteWarehouse(c *gin.Context) {
	var warehouse models.Warehouse
	if err := database.DB.Where("id = ?", c.Param("id")).First(&warehouse).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "仓库不存在"})
		return
	}

	var shelves []models.Shelf
	database.DB.Where("warehouse_id = ?", warehouse.ID).Find(&shelves)
	if len(shelves) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该仓库下还有置物架，请先删除置物架"})
		return
	}

	var boxes []models.Box
	database.DB.Where("warehouse_id = ?", warehouse.ID).Find(&boxes)
	if len(boxes) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该仓库下还有收纳盒，请先删除收纳盒"})
		return
	}

	database.DB.Delete(&warehouse)
	c.JSON(http.StatusOK, gin.H{"message": "仓库已删除"})
}
