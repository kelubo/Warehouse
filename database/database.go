package database

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"warehouse-management/config"
	"warehouse-management/models"
	"warehouse-management/utils"
)

var DB *gorm.DB

// InitDB 初始化数据库连接，支持 MySQL、PostgreSQL 和 SQLite
func InitDB(cfg *config.Config) {
	utils.Info("Initializing database...")
	utils.Info("Database type: " + cfg.DBType)

	var dsn string
	var dialector gorm.Dialector

	switch cfg.DBType {
	case "mysql":
		// MySQL DSN: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
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
		// PostgreSQL DSN
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
		// SQLite DSN
		dsn = cfg.DatabasePath
		dialector = sqlite.Open(dsn)

	default:
		utils.Fatal("Unsupported database type: " + cfg.DBType)
	}

	// 根据环境设置日志级别
	logMode := logger.Info
	if cfg.IsProduction() {
		logMode = logger.Error
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		utils.Fatal("Failed to connect to database")
	}

	// 获取底层 sql.DB 并配置连接池（SQLite 不需要连接池）
	if cfg.DBType != "sqlite" {
		sqlDB, err := db.DB()
		if err != nil {
			utils.Fatal("Failed to get underlying DB")
		}

		// 设置连接池参数
		sqlDB.SetMaxOpenConns(100)                 // 最大打开连接数
		sqlDB.SetMaxIdleConns(10)                  // 最大空闲连接数
		sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大生命周期
		sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 连接最大空闲时间
	}

	DB = db

	utils.Info("Database connection established")

	// 自动迁移
	utils.Info("Running auto migration...")
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Warehouse{},
		&models.Inventory{},
		&models.Order{},
		&models.OrderItem{},
	)
	if err != nil {
		utils.Fatal("Failed to migrate database")
	}

	utils.Info("Database initialized successfully")
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
