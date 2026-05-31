package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Username  string    `gorm:"unique;not null;size:100" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Email     string    `gorm:"unique;size:100" json:"email"`
	Role      string    `gorm:"not null;default:'user';size:20" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

type Product struct {
	ID          string  `gorm:"primaryKey;size:36" json:"id"`
	Name        string  `gorm:"not null;index" json:"name"`
	SKU         string  `gorm:"unique;size:100" json:"sku"`
	Specification string `json:"specification"` // 规格
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `gorm:"index" json:"category"`
	Unit        string  `gorm:"not null;default:'件';size:20" json:"unit"`
	ImageURL    string  `json:"image_url"`
	// 货架位置信息
	WarehouseID string     `gorm:"size:36;index" json:"warehouse_id"`
	ShelfID     string     `gorm:"size:36;index" json:"shelf_id"`
	BoxID       string     `gorm:"size:36;index" json:"box_id"`
	ShelfColumn int        `json:"shelf_column"`
	ShelfRow    int        `json:"shelf_row"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New().String()
	// 如果SKU为空，自动生成一个唯一SKU
	if p.SKU == "" {
		p.SKU = "SKU-" + strings.ToUpper(uuid.New().String()[:8])
	}
	return
}

type Warehouse struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	Name        string     `gorm:"not null;index" json:"name"`
	Location    string     `json:"location"`
	Description string     `json:"description"`
	Capacity    int        `json:"capacity"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`

	Shelves []Shelf `gorm:"foreignKey:WarehouseID" json:"shelves"`
}

func (w *Warehouse) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New().String()
	return
}

type Inventory struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	ProductID   string     `gorm:"not null;size:36;index" json:"product_id"`
	WarehouseID string     `gorm:"not null;size:36;index" json:"warehouse_id"`
	Quantity    int        `gorm:"not null;default:0" json:"quantity"`
	MinStock    int        `gorm:"default:10" json:"min_stock"`
	MaxStock    int        `gorm:"default:1000" json:"max_stock"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at"`

	Product   Product   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Warehouse Warehouse `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.New().String()
	return
}

type Order struct {
	ID        string     `gorm:"primaryKey;size:36" json:"id"`
	OrderNo   string     `gorm:"unique;not null;size:50" json:"order_no"`
	UserID    string     `gorm:"not null;size:36;index" json:"user_id"`
	Type      string     `gorm:"not null;size:20;index" json:"type"`                     // in: 入库, out: 出库
	Status    string     `gorm:"not null;default:'pending';size:20;index" json:"status"` // pending, completed, cancelled
	Remark    string     `json:"remark"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at"`

	User       User        `gorm:"foreignKey:UserID;references:ID" json:"user"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New().String()
	return
}

type OrderItem struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	OrderID     string    `gorm:"not null;size:36;index" json:"order_id"`
	ProductID   string    `gorm:"not null;size:36;index" json:"product_id"`
	WarehouseID string    `gorm:"not null;size:36;index" json:"warehouse_id"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Product   Product   `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	Warehouse Warehouse `gorm:"foreignKey:WarehouseID;references:ID" json:"warehouse"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	oi.ID = uuid.New().String()
	return
}
