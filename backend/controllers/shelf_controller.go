package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetShelves(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.Shelf{})

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var shelves []models.Shelf
	offset := (page - 1) * pageSize
	query.Limit(pageSize).Offset(offset).Find(&shelves)

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, models.PaginationResponse{
		Items:      shelves,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func GetShelvesByWarehouse(c *gin.Context) {
	var shelves []models.Shelf
	database.DB.Where("warehouse_id = ?", c.Param("warehouse_id")).Find(&shelves)
	c.JSON(http.StatusOK, shelves)
}

func GetShelf(c *gin.Context) {
	var shelf models.Shelf
	if err := database.DB.Preload("Boxes").Where("id = ?", c.Param("id")).First(&shelf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "货架不存在"})
		return
	}
	c.JSON(http.StatusOK, shelf)
}

func CreateShelf(c *gin.Context) {
	var shelf models.Shelf
	if err := c.ShouldBindJSON(&shelf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&shelf)
	c.JSON(http.StatusCreated, shelf)
}

func UpdateShelf(c *gin.Context) {
	var shelf models.Shelf
	if err := database.DB.Where("id = ?", c.Param("id")).First(&shelf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "货架不存在"})
		return
	}

	if err := c.ShouldBindJSON(&shelf); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&shelf)
	c.JSON(http.StatusOK, shelf)
}

func DeleteShelf(c *gin.Context) {
	var shelf models.Shelf
	if err := database.DB.Where("id = ?", c.Param("id")).First(&shelf).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "货架不存在"})
		return
	}

	var boxes []models.Box
	database.DB.Where("shelf_id = ?", shelf.ID).Find(&boxes)
	if len(boxes) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该置物架下还有收纳盒，请先删除收纳盒"})
		return
	}

	database.DB.Delete(&shelf)
	c.JSON(http.StatusOK, gin.H{"message": "货架已删除"})
}
