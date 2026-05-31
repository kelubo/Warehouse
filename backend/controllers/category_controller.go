package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category

	database.DB.Where("parent_id IS NULL").Order("created_at desc").Find(&categories)

	for i := range categories {
		loadChildrenRecursive(&categories[i])
	}

	c.JSON(http.StatusOK, categories)
}

func loadChildrenRecursive(category *models.Category) {
	var children []models.Category
	database.DB.Where("parent_id = ?", category.ID).Order("created_at desc").Find(&children)
	category.Children = children
	for i := range children {
		loadChildrenRecursive(&children[i])
	}
}

func GetCategory(c *gin.Context) {
	var category models.Category
	if err := database.DB.Where("id = ?", c.Param("id")).First(&category).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	loadChildrenRecursive(&category)

	c.JSON(http.StatusOK, category)
}

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果设置了 parent_id，检查父分类是否存在
	if category.ParentID != nil && *category.ParentID != "" {
		var parent models.Category
		if err := database.DB.Where("id = ?", category.ParentID).First(&parent).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "父分类不存在"})
			return
		}
	}

	// 检查同一级下分类名称是否重复
	query := database.DB.Where("name = ?", category.Name)
	if category.ParentID != nil && *category.ParentID != "" {
		query = query.Where("parent_id = ?", category.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	
	var existingCategory models.Category
	query.First(&existingCategory)
	if existingCategory.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "该分类名称已存在"})
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

	// 如果设置了 parent_id，检查父分类是否存在，且不能是自己
	if updatedCategory.ParentID != nil && *updatedCategory.ParentID != "" {
		if *updatedCategory.ParentID == category.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能将自己设为父分类"})
			return
		}
		var parent models.Category
		if err := database.DB.Where("id = ?", updatedCategory.ParentID).First(&parent).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "父分类不存在"})
			return
		}
	}

	// 检查同一级下分类名称是否重复（排除自己）
	query := database.DB.Where("name = ? AND id != ?", updatedCategory.Name, category.ID)
	if updatedCategory.ParentID != nil && *updatedCategory.ParentID != "" {
		query = query.Where("parent_id = ?", updatedCategory.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	
	var existingCategory models.Category
	query.First(&existingCategory)
	if existingCategory.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "该分类名称已存在"})
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

	// 检查是否有子分类
	var childCount int64
	database.DB.Model(&models.Category{}).Where("parent_id = ?", category.ID).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下有子分类，无法删除"})
		return
	}

	// 检查是否有产品
	var total int64
	database.DB.Model(&models.Product{}).Where("category = ?", category.Name).Count(&total)
	if total > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该分类下有物品，无法删除"})
		return
	}

	database.DB.Delete(&category)
	c.JSON(http.StatusOK, gin.H{"message": "分类已删除"})
}
