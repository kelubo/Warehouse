package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetBoxes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.Box{})

	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}

	if keyword != "" {
		query = query.Where("box_no LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var boxes []models.Box
	offset := (page - 1) * pageSize
	query.Limit(pageSize).Offset(offset).Find(&boxes)

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, models.PaginationResponse{
		Items:      boxes,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
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

	if box.ShelfID != "" {
		var shelf models.Shelf
		if err := database.DB.Where("id = ?", box.ShelfID).First(&shelf).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "置物架不存在"})
			return
		}
		if box.Column < 1 || box.Column > shelf.Columns {
			c.JSON(http.StatusBadRequest, gin.H{"error": "列号超出范围，该置物架共有 " + fmt.Sprintf("%d", shelf.Columns) + " 列"})
			return
		}
		if box.Row < 1 || box.Row > shelf.Rows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "行号超出范围，该置物架共有 " + fmt.Sprintf("%d", shelf.Rows) + " 行"})
			return
		}
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

	var updatedBox models.Box
	if err := c.ShouldBindJSON(&updatedBox); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedBox.ShelfID != "" {
		var shelf models.Shelf
		if err := database.DB.Where("id = ?", updatedBox.ShelfID).First(&shelf).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "置物架不存在"})
			return
		}
		if updatedBox.Column < 1 || updatedBox.Column > shelf.Columns {
			c.JSON(http.StatusBadRequest, gin.H{"error": "列号超出范围，该置物架共有 " + fmt.Sprintf("%d", shelf.Columns) + " 列"})
			return
		}
		if updatedBox.Row < 1 || updatedBox.Row > shelf.Rows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "行号超出范围，该置物架共有 " + fmt.Sprintf("%d", shelf.Rows) + " 行"})
			return
		}
	}

	database.DB.Model(&box).Updates(updatedBox)
	c.JSON(http.StatusOK, box)
}

func DeleteBox(c *gin.Context) {
	var box models.Box
	if err := database.DB.Where("id = ?", c.Param("id")).First(&box).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "箱子不存在"})
		return
	}

	var products []models.Product
	database.DB.Where("box_id = ?", box.ID).Find(&products)
	if len(products) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该收纳盒中还有物品，请先删除物品"})
		return
	}

	database.DB.Delete(&box)
	c.JSON(http.StatusOK, gin.H{"message": "箱子已删除"})
}
