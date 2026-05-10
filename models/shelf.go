package models

import "gorm.io/gorm"

type Shelf struct {
	gorm.Model
	WarehouseID string `json:"warehouse_id" gorm:"size:36;not null;index"`
	Name        string `json:"name"`
	Columns     int    `json:"columns"`
	Rows        int    `json:"rows"`
}