package initialize

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"warehouse-management/config"
	"warehouse-management/database"
	"warehouse-management/models"
	"warehouse-management/utils"
)

// Run 执行数据库初始化流程
func Run() {
	fmt.Println("==========================================")
	fmt.Println("    仓库管理系统 - 数据库初始化")
	fmt.Println("==========================================")
	fmt.Println()

	// 加载当前配置
	cfg := config.LoadConfig()

	// 询问用户数据库类型
	fmt.Println("请选择数据库类型：")
	fmt.Println("1. MySQL (推荐)")
	fmt.Println("2. PostgreSQL")
	fmt.Println("3. SQLite")
	fmt.Print("请输入选择 (1-3): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	dbTypeChoice := scanner.Text()

	switch dbTypeChoice {
	case "1":
		cfg.DBType = "mysql"
	case "2":
		cfg.DBType = "postgres"
	case "3":
		cfg.DBType = "sqlite"
	default:
		fmt.Println("无效选择，使用默认 MySQL")
		cfg.DBType = "mysql"
	}

	if cfg.DBType != "sqlite" {
		// 获取数据库连接信息
		fmt.Println()
		fmt.Println("请输入数据库连接信息：")

		fmt.Print("数据库主机 (默认: localhost): ")
		scanner.Scan()
		host := scanner.Text()
		if host != "" {
			cfg.DBHost = host
		}

		fmt.Print("数据库端口 (MySQL: 3306, PostgreSQL: 5432): ")
		scanner.Scan()
		port := scanner.Text()
		if port != "" {
			cfg.DBPort = port
		}

		fmt.Print("数据库用户名: ")
		scanner.Scan()
		user := scanner.Text()
		if user != "" {
			cfg.DBUser = user
		}

		fmt.Print("数据库密码 (可省略): ")
		scanner.Scan()
		password := scanner.Text()
		if password != "" {
			cfg.DBPassword = password
		}

		fmt.Print("数据库名称 (默认: warehouse): ")
		scanner.Scan()
		dbName := scanner.Text()
		if dbName != "" {
			cfg.DBName = dbName
		}
	} else {
		fmt.Print("SQLite 数据库文件路径 (默认: ./warehouse.db): ")
		scanner.Scan()
		path := scanner.Text()
		if path != "" {
			cfg.DatabasePath = path
		}
	}

	// 显示配置摘要
	fmt.Println()
	fmt.Println("配置摘要：")
	fmt.Printf("  数据库类型: %s\n", cfg.DBType)
	if cfg.DBType != "sqlite" {
		fmt.Printf("  主机: %s:%s\n", cfg.DBHost, cfg.DBPort)
		fmt.Printf("  用户: %s\n", cfg.DBUser)
		fmt.Printf("  数据库: %s\n", cfg.DBName)
	} else {
		fmt.Printf("  文件路径: %s\n", cfg.DatabasePath)
	}
	fmt.Println()

	// 确认初始化
	fmt.Print("确认初始化数据库? (y/N): ")
	scanner.Scan()
	confirm := strings.ToLower(scanner.Text())
	if confirm != "y" && confirm != "yes" {
		fmt.Println("初始化已取消")
		return
	}

	// 初始化数据库
	fmt.Println()
	fmt.Println("正在初始化数据库...")

	// 初始化日志
	utils.InitLogger("development")

	// 创建数据库连接
	database.InitDB(cfg)

	// 检查数据库是否已存在数据
	if isDatabaseInitialized() {
		fmt.Println()
		fmt.Println("检测到数据库已存在！")
		fmt.Print("是否重新初始化（将删除所有数据）? (y/N): ")
		scanner.Scan()
		reinit := strings.ToLower(scanner.Text())
		if reinit != "y" && reinit != "yes" {
			fmt.Println("跳过数据库初始化")
		} else {
			// 重新初始化 - 删除所有表并重新创建
			fmt.Println("正在重新初始化数据库...")
			resetDatabase()
			fmt.Println("数据库重新初始化成功！")
		}
	} else {
		fmt.Println("数据库初始化成功！")
	}

	// 创建初始管理员账户
	fmt.Println()
	fmt.Println("创建初始管理员账户：")

	fmt.Print("管理员用户名 (默认: admin): ")
	scanner.Scan()
	username := scanner.Text()
	if username == "" {
		username = "admin"
	}

	fmt.Print("管理员邮箱 (默认: admin@example.com): ")
	scanner.Scan()
	email := scanner.Text()
	if email == "" {
		email = "admin@example.com"
	}

	fmt.Print("管理员密码: ")
	scanner.Scan()
	password := scanner.Text()
	if password == "" {
		password = "admin123"
		fmt.Println("未输入密码，使用默认密码: admin123")
	}

	// 检查用户是否已存在
	var existingUser models.User
	result := database.DB.Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		fmt.Println("警告：该邮箱已存在用户，跳过创建")
	} else {
		// 创建管理员用户
		hashedPassword, _ := utils.HashPassword(password)
		admin := models.User{
			Username: username,
			Email:    email,
			Password: hashedPassword,
			Role:     "admin",
		}

		err := database.DB.Create(&admin).Error
		if err != nil {
			fmt.Printf("创建管理员失败: %v\n", err)
		} else {
			fmt.Println("管理员账户创建成功！")
		}
	}

	// 保存配置到 .env 文件
	err := saveConfigToEnv(cfg)
	if err != nil {
		fmt.Printf("保存配置失败: %v\n", err)
	} else {
		fmt.Println("配置已保存到 .env 文件")
	}

	fmt.Println()
	fmt.Println("==========================================")
	fmt.Println("    初始化完成！")
	fmt.Println("==========================================")
	fmt.Printf("启动命令: go run main.go\n")
	fmt.Printf("访问地址: http://localhost:%s\n", cfg.Port)
	fmt.Printf("管理员账户: %s / %s\n", email, password)
}

// isDatabaseInitialized 检查数据库是否已初始化（是否存在用户表数据）
func isDatabaseInitialized() bool {
	var count int64
	// 检查 users 表是否存在数据
	database.DB.Model(&models.User{}).Count(&count)
	return count > 0
}

// resetDatabase 重置数据库（删除所有表数据）
func resetDatabase() {
	// 按依赖顺序删除表数据
	database.DB.Exec("DELETE FROM order_items")
	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM inventories")
	database.DB.Exec("DELETE FROM products")
	database.DB.Exec("DELETE FROM warehouses")
	database.DB.Exec("DELETE FROM users")
}

// saveConfigToEnv 将配置保存到 .env 文件
func saveConfigToEnv(cfg *config.Config) error {
	content := fmt.Sprintf(`# 环境配置
ENV=development
PORT=%s

# 数据库配置
DB_TYPE=%s
DB_HOST=%s
DB_USER=%s
DB_PASSWORD=%s
DB_NAME=%s
DB_PORT=%s

# SQLite 专用配置
DATABASE_PATH=%s

# JWT 配置
JWT_SECRET=%s
JWT_EXPIRE_HOURS=%d

# 日志配置
LOG_LEVEL=%s
`,
		cfg.Port,
		cfg.DBType,
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DatabasePath,
		cfg.JWTSecret,
		cfg.JWTExpireHours,
		cfg.LogLevel,
	)

	return os.WriteFile(".env", []byte(content), 0644)
}
