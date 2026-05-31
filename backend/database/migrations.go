package database

import (
	"warehouse-management/backend/utils"
)

// RunCompositeIndexes 创建复合索引以优化查询性能
func RunCompositeIndexes() error {
	utils.Info("Creating composite indexes...")

	// 产品表复合索引：按仓库+货架查询
	if err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_products_warehouse_shelf ON products(warehouse_id, shelf_id)`).Error; err != nil {
		utils.Warn("Index idx_products_warehouse_shelf may already exist: " + err.Error())
	}

	// 产品表复合索引：按分类+仓库查询
	if err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_products_category_warehouse ON products(category, warehouse_id)`).Error; err != nil {
		utils.Warn("Index idx_products_category_warehouse may already exist: " + err.Error())
	}

	// 库存表复合索引：按仓库+产品查询
	if err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_inventory_warehouse_product ON inventory(warehouse_id, product_id)`).Error; err != nil {
		utils.Warn("Index idx_inventory_warehouse_product may already exist: " + err.Error())
	}

	// 货架表复合索引：按仓库查询
	if err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_shelves_warehouse ON shelves(warehouse_id)`).Error; err != nil {
		utils.Warn("Index idx_shelves_warehouse may already exist: " + err.Error())
	}

	// 箱子表复合索引：按货架查询
	if err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_boxes_shelf ON boxes(shelf_id)`).Error; err != nil {
		utils.Warn("Index idx_boxes_shelf may already exist: " + err.Error())
	}

	// 订单表复合索引：按用户+状态查询
	if err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status)`).Error; err != nil {
		utils.Warn("Index idx_orders_user_status may already exist: " + err.Error())
	}

	// 订单表复合索引：按类型+状态+创建时间查询
	if err := DB.Exec(`CREATE INDEX IF NOT EXISTS idx_orders_type_status_created ON orders(type, status, created_at)`).Error; err != nil {
		utils.Warn("Index idx_orders_type_status_created may already exist: " + err.Error())
	}

	utils.Info("Composite indexes created successfully")
	return nil
}