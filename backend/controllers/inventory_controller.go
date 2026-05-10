package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetInventories(c *gin.Context) {
	var inventories []models.Inventory
	database.DB.Preload("Product").Preload("Warehouse").Find(&inventories)
	c.JSON(http.StatusOK, inventories)
}

func GetInventory(c *gin.Context) {
	var inventory models.Inventory
	if err := database.DB.Preload("Product").Preload("Warehouse").Where("id = ?", c.Param("id")).First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "库存记录不存在"})
		return
	}
	c.JSON(http.StatusOK, inventory)
}

func CreateInventory(c *gin.Context) {
	var inventory models.Inventory
	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查产品是否存在
	var product models.Product
	if err := database.DB.Where("id = ?", inventory.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "产品不存在"})
		return
	}

	// 检查仓库是否存在
	var warehouse models.Warehouse
	if err := database.DB.Where("id = ?", inventory.WarehouseID).First(&warehouse).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "仓库不存在"})
		return
	}

	database.DB.Create(&inventory)
	c.JSON(http.StatusCreated, inventory)
}

func UpdateInventory(c *gin.Context) {
	var inventory models.Inventory
	if err := database.DB.Where("id = ?", c.Param("id")).First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "库存记录不存在"})
		return
	}

	if err := c.ShouldBindJSON(&inventory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&inventory)
	c.JSON(http.StatusOK, inventory)
}

func DeleteInventory(c *gin.Context) {
	var inventory models.Inventory
	if err := database.DB.Where("id = ?", c.Param("id")).First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "库存记录不存在"})
		return
	}

	database.DB.Delete(&inventory)
	c.JSON(http.StatusOK, gin.H{"message": "库存记录已删除"})
}
