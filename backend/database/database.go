package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"warehouse-management/backend/config"
	"warehouse-management/backend/models"
	"warehouse-management/backend/utils"
)

var DB *gorm.DB

// InitDB 初始化数据库连接，支持 MySQL、PostgreSQL 和 SQLite
func InitDB(cfg *config.Config) error {
	utils.Info("Initializing database...")
	utils.Info("Database type: " + cfg.DBType)

	var dsn string
	var dialector gorm.Dialector

	switch cfg.DBType {
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)
		dialector = mysql.Open(dsn)

	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
			cfg.DBPort,
		)
		dialector = postgres.Open(dsn)

	case "sqlite":
		dsn = cfg.DatabasePath
		dialector = sqlite.Open(dsn)

	default:
		return fmt.Errorf("unsupported database type: %s", cfg.DBType)
	}

	logMode := logger.Info
	if cfg.IsProduction() {
		logMode = logger.Error
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if cfg.DBType != "sqlite" {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying DB: %w", err)
		}

		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)
		sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	}

	DB = db

	utils.Info("Database connection established")

	utils.Info("Fixing invalid foreign key references...")
	if cfg.DBType == "mysql" {
		// 修复 products 表中无效的外键引用
		db.Exec(`UPDATE products p LEFT JOIN boxes b ON p.box_id = b.id SET p.box_id = NULL WHERE b.id IS NULL AND p.box_id IS NOT NULL`)
		db.Exec(`UPDATE products p LEFT JOIN shelves s ON p.shelf_id = s.id SET p.shelf_id = NULL WHERE s.id IS NULL AND p.shelf_id IS NOT NULL`)
		db.Exec(`UPDATE products p LEFT JOIN warehouses w ON p.warehouse_id = w.id SET p.warehouse_id = NULL WHERE w.id IS NULL AND p.warehouse_id IS NOT NULL`)
		db.Exec(`UPDATE boxes b LEFT JOIN shelves s ON b.shelf_id = s.id SET b.shelf_id = NULL WHERE s.id IS NULL AND b.shelf_id IS NOT NULL`)
		db.Exec(`UPDATE shelves s LEFT JOIN warehouses w ON s.warehouse_id = w.id SET s.warehouse_id = NULL WHERE w.id IS NULL AND s.warehouse_id IS NOT NULL`)
		// 确保 products 表中有 unit 字段
		db.Exec(`ALTER TABLE products ADD COLUMN IF NOT EXISTS unit VARCHAR(20) DEFAULT '件'`)
	}

	utils.Info("Running auto migration...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Warehouse{},
		&models.Shelf{},
		&models.Box{},
		&models.Unit{},
		&models.Category{},
		&models.Inventory{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// 创建复合索引以优化查询性能
	RunCompositeIndexes()

	utils.Info("Database initialized successfully")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// CreateDatabaseIfNotExists 创建数据库（如果不存在）
func CreateDatabaseIfNotExists(cfg *config.Config) error {
	var dsn string
	var dialector gorm.Dialector

	switch cfg.DBType {
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
		)
		dialector = mysql.Open(dsn)

	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=Asia/Shanghai",
			cfg.DBHost,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBPort,
		)
		dialector = postgres.Open(dsn)

	default:
		return fmt.Errorf("unsupported database type for creating database: %s", cfg.DBType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database server: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying DB: %w", err)
	}
	defer sqlDB.Close()

	var createSQL string
	switch cfg.DBType {
	case "mysql":
		createSQL = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName)
	case "postgres":
		createSQL = fmt.Sprintf("CREATE DATABASE \"%s\"", cfg.DBName)
	}

	if err := db.Exec(createSQL).Error; err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}
