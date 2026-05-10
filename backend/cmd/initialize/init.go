package initialize

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"warehouse-management/backend/config"
	"warehouse-management/backend/database"
	"warehouse-management/backend/models"
	"warehouse-management/backend/utils"
)

var reader *bufio.Reader

func init() {
	reader = bufio.NewReader(os.Stdin)
}

func readLine(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return input
}

func readPassword(prompt string) string {
	fmt.Print(prompt)
	bytes, _ := reader.ReadBytes('\n')
	password := strings.TrimSpace(string(bytes))
	return password
}

func Run() {
	fmt.Println("==========================================")
	fmt.Println("    仓库管理系统 - 数据库初始化")
	fmt.Println("==========================================")
	fmt.Println()

	cfg := config.LoadConfig()

	fmt.Println("请选择数据库类型：")
	fmt.Println("1. MySQL (推荐)")
	fmt.Println("2. PostgreSQL")
	fmt.Println("3. SQLite")
	fmt.Print("请输入选择 (1-3): ")

	dbTypeChoice, _ := reader.ReadString('\n')
	dbTypeChoice = strings.TrimSpace(dbTypeChoice)

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
		fmt.Println()
		fmt.Println("请输入数据库连接信息：")

		host := readLine("数据库主机 (直接回车使用 localhost): ")
		if host == "" {
			host = "localhost"
		}
		cfg.DBHost = host

		port := readLine("数据库端口 (直接回车使用 3306): ")
		if port == "" {
			port = "3306"
		}
		cfg.DBPort = port

		user := readLine("数据库用户名 (直接回车使用 root): ")
		if user == "" {
			user = "root"
		}
		cfg.DBUser = user

		password := readPassword("数据库密码 (直接回车使用空密码): ")
		cfg.DBPassword = password

		dbName := readLine("数据库名称 (直接回车使用 warehouse): ")
		if dbName == "" {
			dbName = "warehouse"
		}
		cfg.DBName = dbName
	} else {
		path := readLine("SQLite 数据库文件路径 (直接回车使用 ./warehouse.db): ")
		if path == "" {
			path = "./warehouse.db"
		}
		cfg.DatabasePath = path
	}

	fmt.Println()
	fmt.Println("配置摘要：")
	fmt.Printf("  数据库类型: %s\n", cfg.DBType)
	if cfg.DBType != "sqlite" {
		fmt.Printf("  主机: %s:%s\n", cfg.DBHost, cfg.DBPort)
		fmt.Printf("  用户: %s\n", cfg.DBUser)
		if cfg.DBPassword != "" {
			fmt.Printf("  密码: %s\n", maskPassword(cfg.DBPassword))
		} else {
			fmt.Printf("  密码: (空)\n")
		}
		fmt.Printf("  数据库: %s\n", cfg.DBName)
	} else {
		fmt.Printf("  文件路径: %s\n", cfg.DatabasePath)
	}
	fmt.Println()

	confirm := readLine("确认初始化数据库? (y/N): ")
	confirm = strings.ToLower(confirm)
	if confirm != "y" && confirm != "yes" {
		fmt.Println("初始化已取消")
		return
	}

	fmt.Println()
	err := saveConfigToEnv(cfg)
	if err != nil {
		fmt.Printf("警告：保存配置失败: %v\n", err)
	} else {
		fmt.Println("配置已保存到 .env 文件")
	}

	fmt.Println()
	fmt.Println("正在初始化数据库...")

	utils.InitLogger("development")

	if err := database.InitDB(cfg); err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)

		if cfg.DBType == "mysql" || cfg.DBType == "postgres" {
			fmt.Println()
			fmt.Println("尝试创建数据库...")

			if err := database.CreateDatabaseIfNotExists(cfg); err != nil {
				fmt.Printf("创建数据库失败: %v\n", err)
				fmt.Println("配置已保存，您可以稍后检查数据库连接信息后重试")
				os.Exit(1)
			}

			fmt.Println("数据库创建成功，正在重新连接...")

			if err := database.InitDB(cfg); err != nil {
				fmt.Printf("重新连接失败: %v\n", err)
				fmt.Println("配置已保存，您可以稍后检查数据库连接信息后重试")
				os.Exit(1)
			}
		} else {
			fmt.Println("配置已保存，您可以稍后检查数据库连接信息后重试")
			os.Exit(1)
		}
	}

	if isDatabaseInitialized() {
		fmt.Println()
		fmt.Println("检测到数据库已存在！")
		reinit := readLine("是否重新初始化（将删除所有数据）? (y/N): ")
		reinit = strings.ToLower(reinit)
		if reinit != "y" && reinit != "yes" {
			fmt.Println("跳过数据库初始化")
		} else {
			fmt.Println("正在重新初始化数据库...")
			resetDatabase()
			fmt.Println("数据库重新初始化成功！")
		}
	} else {
		fmt.Println("数据库初始化成功！")
	}

	fmt.Println()
	fmt.Println("创建初始管理员账户：")

	username := readLine("管理员用户名 (默认: admin): ")
	if username == "" {
		username = "admin"
	}

	email := readLine("管理员邮箱 (默认: admin@example.com): ")
	if email == "" {
		email = "admin@example.com"
	}

	password := readPassword("管理员密码: ")
	if password == "" {
		password = "admin123"
		fmt.Println("未输入密码，使用默认密码: admin123")
	}

	var existingUser models.User
	result := database.DB.Where("email = ?", email).First(&existingUser)
	if result.Error == nil {
		fmt.Println("警告：该邮箱已存在用户，跳过创建")
	} else {
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

	fmt.Println()
	fmt.Println("==========================================")
	fmt.Println("    初始化完成！")
	fmt.Println("==========================================")
	fmt.Printf("启动命令: ./warehouse-management\n")
	fmt.Printf("访问地址: http://localhost:%s\n", cfg.Port)
	fmt.Printf("管理员账户: %s / %s\n", email, password)
}

func isDatabaseInitialized() bool {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	return count > 0
}

func resetDatabase() {
	database.DB.Exec("DELETE FROM order_items")
	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM inventories")
	database.DB.Exec("DELETE FROM products")
	database.DB.Exec("DELETE FROM warehouses")
	database.DB.Exec("DELETE FROM users")
}

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

func maskPassword(password string) string {
	if len(password) <= 2 {
		return "**"
	}
	return string(password[0]) + strings.Repeat("*", len(password)-2) + string(password[len(password)-1])
}
