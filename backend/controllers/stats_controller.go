package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetStats(c *gin.Context) {
	var productCount int64
	database.DB.Model(&models.Product{}).Count(&productCount)

	var warehouseCount int64
	database.DB.Model(&models.Warehouse{}).Count(&warehouseCount)

	var shelfCount int64
	database.DB.Model(&models.Shelf{}).Count(&shelfCount)

	var boxCount int64
	database.DB.Model(&models.Box{}).Count(&boxCount)

	var inventoryCount int64
	database.DB.Model(&models.Inventory{}).Count(&inventoryCount)

	var totalQuantity int64
	database.DB.Model(&models.Inventory{}).Select("SUM(quantity)").Scan(&totalQuantity)

	c.JSON(http.StatusOK, gin.H{
		"product_count":   productCount,
		"warehouse_count": warehouseCount,
		"shelf_count":     shelfCount,
		"box_count":       boxCount,
		"inventory_count": inventoryCount,
		"total_quantity":  totalQuantity,
	})
}

func GetWarehouseStats(c *gin.Context) {
	var warehouses []models.Warehouse
	database.DB.Find(&warehouses)

	var stats []gin.H
	for _, warehouse := range warehouses {
		var shelfCount int64
		database.DB.Model(&models.Shelf{}).Where("warehouse_id = ?", warehouse.ID).Count(&shelfCount)

		var boxCount int64
		database.DB.Model(&models.Box{}).Where("warehouse_id = ?", warehouse.ID).Count(&boxCount)

		var inventoryCount int64
		database.DB.Model(&models.Inventory{}).Where("warehouse_id = ?", warehouse.ID).Count(&inventoryCount)

		var totalQuantity int64
		database.DB.Model(&models.Inventory{}).Where("warehouse_id = ?", warehouse.ID).Select("SUM(quantity)").Scan(&totalQuantity)

		stats = append(stats, gin.H{
			"warehouse_id":   warehouse.ID,
			"warehouse_name": warehouse.Name,
			"shelf_count":    shelfCount,
			"box_count":      boxCount,
			"inventory_count": inventoryCount,
			"total_quantity":  totalQuantity,
		})
	}

	c.JSON(http.StatusOK, stats)
}

func GetCategoryStats(c *gin.Context) {
	var categories []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}

	database.DB.Model(&models.Product{}).Select("category, COUNT(*) as count").Group("category").Scan(&categories)

	c.JSON(http.StatusOK, categories)
}

func GetRecentActivities(c *gin.Context) {
	var activities []gin.H

	database.DB.Model(&models.Inventory{}).Select("created_at, 'inventory' as type, product_id").Order("created_at DESC").Limit(10).Scan(&activities)

	c.JSON(http.StatusOK, activities)
}
