package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment    string
	Port           string
	DatabasePath   string
	JWTSecret      string
	JWTExpireHours int

	// 数据库配置
	DBType     string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string

	// 日志配置
	LogLevel string
}

var cfg *Config

func LoadConfig() *Config {
	if cfg != nil {
		return cfg
	}

	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	jwtExpireHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	if err != nil {
		jwtExpireHours = 24
	}

	cfg = &Config{
		Environment:    getEnv("ENV", "development"),
		Port:           getEnv("PORT", "8080"),
		DatabasePath:   getEnv("DATABASE_PATH", "./warehouse.db"),
		JWTSecret:      getEnv("JWT_SECRET", "warehouse-secret-key"),
		JWTExpireHours: jwtExpireHours,
		DBType:         getEnv("DB_TYPE", "mysql"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "warehouse"),
		DBPort:         getEnv("DB_PORT", "3306"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}
