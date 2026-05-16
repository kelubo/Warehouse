package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetUnits(c *gin.Context) {
	var units []models.Unit
	database.DB.Order("created_at desc").Find(&units)
	c.JSON(http.StatusOK, units)
}

func GetUnit(c *gin.Context) {
	var unit models.Unit
	if err := database.DB.Where("id = ?", c.Param("id")).First(&unit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "单位不存在"})
		return
	}
	c.JSON(http.StatusOK, unit)
}

func CreateUnit(c *gin.Context) {
	var unit models.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUnit models.Unit
	database.DB.Where("name = ?", unit.Name).First(&existingUnit)
	if existingUnit.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "单位已存在"})
		return
	}

	database.DB.Create(&unit)
	c.JSON(http.StatusCreated, unit)
}

func UpdateUnit(c *gin.Context) {
	var unit models.Unit
	if err := database.DB.Where("id = ?", c.Param("id")).First(&unit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "单位不存在"})
		return
	}

	var updatedUnit models.Unit
	if err := c.ShouldBindJSON(&updatedUnit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUnit models.Unit
	database.DB.Where("name = ? AND id != ?", updatedUnit.Name, unit.ID).First(&existingUnit)
	if existingUnit.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "单位名称已存在"})
		return
	}

	database.DB.Model(&unit).Updates(updatedUnit)
	c.JSON(http.StatusOK, unit)
}

func DeleteUnit(c *gin.Context) {
	var unit models.Unit
	if err := database.DB.Where("id = ?", c.Param("id")).First(&unit).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "单位不存在"})
		return
	}

	var total int64
	database.DB.Model(&models.Product{}).Where("unit = ?", unit.Name).Count(&total)
	if total > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该单位下有物品，无法删除"})
		return
	}

	database.DB.Delete(&unit)
	c.JSON(http.StatusOK, gin.H{"message": "单位已删除"})
}
