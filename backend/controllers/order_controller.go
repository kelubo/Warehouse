package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetOrders(c *gin.Context) {
	var orders []models.Order
	database.DB.Preload("User").Preload("OrderItems").Find(&orders)
	c.JSON(http.StatusOK, orders)
}

func GetOrder(c *gin.Context) {
	var order models.Order
	if err := database.DB.Preload("User").Preload("OrderItems").Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成订单号
	order.OrderNo = "ORD" + time.Now().Format("20060102150405")

	// 开启事务
	tx := database.DB.Begin()

	for _, item := range order.OrderItems {
		// 检查库存
		var inventory models.Inventory
		if err := tx.Where("product_id = ? AND warehouse_id = ?", item.ProductID, item.WarehouseID).First(&inventory).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "库存记录不存在"})
			return
		}

		if order.Type == "out" && inventory.Quantity < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "库存不足"})
			return
		}

		// 更新库存
		if order.Type == "in" {
			tx.Model(&inventory).Update("quantity", inventory.Quantity+item.Quantity)
		} else {
			tx.Model(&inventory).Update("quantity", inventory.Quantity-item.Quantity)
		}
	}

	tx.Create(&order)
	tx.Commit()

	c.JSON(http.StatusCreated, order)
}

func UpdateOrder(c *gin.Context) {
	var order models.Order
	if err := database.DB.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&order)
	c.JSON(http.StatusOK, order)
}

func DeleteOrder(c *gin.Context) {
	var order models.Order
	if err := database.DB.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	database.DB.Delete(&order)
	c.JSON(http.StatusOK, gin.H{"message": "订单已删除"})
}
