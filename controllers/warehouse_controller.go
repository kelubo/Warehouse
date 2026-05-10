package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-management/database"
	"warehouse-management/models"
)

func GetWarehouses(c *gin.Context) {
	var warehouses []models.Warehouse
	database.DB.Find(&warehouses)
	c.JSON(http.StatusOK, warehouses)
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

	database.DB.Delete(&warehouse)
	c.JSON(http.StatusOK, gin.H{"message": "仓库已删除"})
}
