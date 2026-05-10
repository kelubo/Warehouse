package models

import "gorm.io/gorm"

type Box struct {
	gorm.Model
	WarehouseID string `json:"warehouse_id" gorm:"size:36;not null;index"`
	ShelfID     string `json:"shelf_id" gorm:"size:36;index"`
	BoxNo       string `json:"box_no" gorm:"unique;not null"`
	Column      int    `json:"column"`
	Row         int    `json:"row"`
}
