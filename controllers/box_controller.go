package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-management/database"
	"warehouse-management/models"
)

func GetBoxes(c *gin.Context) {
	var boxes []models.Box
	query := database.DB

	// 通过查询参数过滤仓库
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	query.Find(&boxes)
	c.JSON(http.StatusOK, boxes)
}

func GetBoxesByWarehouse(c *gin.Context) {
	var boxes []models.Box
	database.DB.Where("warehouse_id = ?", c.Param("warehouse_id")).Find(&boxes)
	c.JSON(http.StatusOK, boxes)
}

func GetBox(c *gin.Context) {
	var box models.Box
	if err := database.DB.Where("id = ?", c.Param("id")).First(&box).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "箱子不存在"})
		return
	}
	c.JSON(http.StatusOK, box)
}

func CreateBox(c *gin.Context) {
	var box models.Box
	if err := c.ShouldBindJSON(&box); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&box)
	c.JSON(http.StatusCreated, box)
}

func UpdateBox(c *gin.Context) {
	var box models.Box
	if err := database.DB.Where("id = ?", c.Param("id")).First(&box).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "箱子不存在"})
		return
	}

	if err := c.ShouldBindJSON(&box); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&box)
	c.JSON(http.StatusOK, box)
}

func DeleteBox(c *gin.Context) {
	var box models.Box
	if err := database.DB.Where("id = ?", c.Param("id")).First(&box).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "箱子不存在"})
		return
	}

	database.DB.Delete(&box)
	c.JSON(http.StatusOK, gin.H{"message": "箱子已删除"})
}
