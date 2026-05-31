package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
	"warehouse-management/backend/utils"
)

func GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	query := database.DB.Model(&models.Product{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR sku LIKE ? OR category LIKE ? OR specification LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var products []models.Product
	offset := (page - 1) * pageSize
	query.Limit(pageSize).Offset(offset).Find(&products)

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, models.PaginationResponse{
		Items:      products,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	})
}

func GetProductsCursor(c *gin.Context) {
	cursor := c.Query("cursor")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	direction := c.DefaultQuery("direction", "next")
	keyword := c.Query("keyword")

	params := utils.GetCursorParams(cursor, limit, direction)

	query := database.DB.Model(&models.Product{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR sku LIKE ? OR category LIKE ? OR specification LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	database.DB.Model(&models.Product{}).Count(&total)

	query = utils.ApplyCursorPagination(query, params, "id")

	var products []models.Product
	query.Find(&products)

	hasMore := len(products) == params.Limit
	var lastID string
	if len(products) > 0 {
		lastID = products[len(products)-1].ID
	}

	response := utils.BuildCursorResponse(products, lastID, hasMore, total)
	c.JSON(http.StatusOK, response)
}

func GetProduct(c *gin.Context) {
	var product models.Product
	if err := database.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if product.SKU != "" {
		var existingProduct models.Product
		database.DB.Where("sku = ?", product.SKU).First(&existingProduct)
		if existingProduct.ID != "" {
			c.JSON(http.StatusConflict, gin.H{"error": "SKU 已存在"})
			return
		}
	}

	database.DB.Create(&product)
	c.JSON(http.StatusCreated, product)
}

func UpdateProduct(c *gin.Context) {
	var product models.Product
	if err := database.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Save(&product)
	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	var product models.Product
	if err := database.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		return
	}

	if product.ImageURL != "" {
		imagePath := filepath.Join("uploads", filepath.Base(product.ImageURL))
		os.Remove(imagePath)
	}

	database.DB.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"message": "产品已删除"})
}

func UploadProductImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请选择要上传的图片"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只支持 JPG、PNG、GIF 格式的图片"})
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少产品ID"})
		return
	}

	var product models.Product
	if err := database.DB.Where("id = ?", productID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品不存在"})
		return
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建上传目录失败"})
		return
	}

	filename := "product_" + productID + ext
	filePath := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败"})
		return
	}

	oldImageURL := product.ImageURL
	product.ImageURL = "/uploads/" + filename
	database.DB.Save(&product)

	if oldImageURL != "" {
		oldPath := filepath.Join("uploads", filepath.Base(oldImageURL))
		os.Remove(oldPath)
	}

	c.JSON(http.StatusOK, gin.H{"image_url": product.ImageURL})
}
