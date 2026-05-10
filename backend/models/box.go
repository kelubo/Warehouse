package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Box struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	WarehouseID string     `json:"warehouse_id" gorm:"size:36;not null;index"`
	ShelfID     string     `json:"shelf_id" gorm:"size:36;index"`
	BoxNo       string     `json:"box_no" gorm:"unique;not null;index"`
	Column      int        `json:"column"`
	Row         int        `json:"row"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

func (b *Box) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}
