package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Order("created_at desc").Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func GetCategory(c *gin.Context) {
	var category models.Category
	if err := database.DB.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}
	c.JSON(http.StatusOK, category)
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingCategory models.Category
	database.DB.Where("name = ?", category.Name).First(&existingCategory)
	if existingCategory.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "分类已存在"})
		return
	}

	database.DB.Create(&category)
	c.JSON(http.StatusCreated, category)
}

func UpdateCategory(c *gin.Context) {
	var category models.Category
	if err := database.DB.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	var updatedCategory models.Category
	if err := c.ShouldBindJSON(&updatedCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingCategory models.Category
	database.DB.Where("name = ? AND id != ?", updatedCategory.Name, category.ID).First(&existingCategory)
	if existingCategory.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "分类名称已存在"})
		return
	}

	database.DB.Model(&category).Updates(updatedCategory)
	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	var category models.Category
	if err := database.DB.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	var total int64
	database.DB.Model(&models.Product{}).Where("category = ?", category.Name).Count(&total)
	if total > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下有物品，无法删除"})
		return
	}

	database.DB.Delete(&category)
	c.JSON(http.StatusOK, gin.H{"message": "分类已删除"})
}
