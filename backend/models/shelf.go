package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shelf struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	WarehouseID string     `json:"warehouse_id" gorm:"size:36;not null;index"`
	Name        string     `json:"name" gorm:"index"`
	Columns     int        `json:"columns"`
	Rows        int        `json:"rows"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

func (s *Shelf) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New().String()
	return
}
