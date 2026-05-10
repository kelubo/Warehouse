package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-management/database"
	"warehouse-management/models"
)

func GetShelves(c *gin.Context) {
	var shelves []models.Shelf
	query := database.DB

	// 通过查询参数过滤仓库
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	query.Find(&shelves)
	c.JSON(http.StatusOK, shelves)
}

func GetShelvesByWarehouse(c *gin.Context) {
	var shelves []models.Shelf
	database.DB.Where("warehouse_id = ?", c.Param("warehouse_id")).Find(&shelves)
	c.JSON(http.StatusOK, shelves)
}

func GetShelf(c *gin.Context) {
	var shelf models.Shelf
	if err := database.DB.Where("id = ?", c.Param("id")).First(&shelf).Error; err != nil {
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

	database.DB.Delete(&shelf)
	c.JSON(http.StatusOK, gin.H{"message": "货架已删除"})
}
