package server

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"warehouse-management/backend/config"
	"warehouse-management/backend/database"
	"warehouse-management/backend/middleware"
	"warehouse-management/backend/routes"
	"warehouse-management/backend/utils"
)

// @title 仓库管理系统 API
// @version 1.0
// @description 基于 Go + Gin 的仓库管理系统
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// Run 启动服务
func Run() {
	// 加载配置
	cfg := config.LoadConfig()

	// 初始化日志
	utils.InitLogger(cfg.Environment)
	utils.Info("Starting Warehouse Management System...")

	// 初始化数据库
	utils.Info("Initializing database...")
	if err := database.InitDB(cfg); err != nil {
		utils.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 设置 Gin 模式
	utils.Info("Setting Gin mode...")
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		utils.Info("Gin mode: Release")
	} else {
		gin.SetMode(gin.DebugMode)
		utils.Info("Gin mode: Debug")
	}

	// 创建 Gin 引擎
	utils.Info("Creating Gin engine...")
	r := gin.Default()
	
	// 设置文件上传大小限制（32MB）
	r.MaxMultipartMemory = 32 << 20

	// 注册安全中间件
	utils.Info("Registering security middlewares...")
	for _, mid := range middleware.SecurityMiddlewares() {
		r.Use(mid)
	}

	// 注册错误处理中间件
	utils.Info("Registering error handler middleware...")
	r.Use(middleware.ErrorHandler())

	// 注册路由
	utils.Info("Registering routes...")
	routes.RegisterRoutes(r)

	// 注册 Swagger API 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	utils.Info("Starting server on port " + cfg.Port)
	utils.Info("Server is ready at http://localhost:" + cfg.Port)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(":" + cfg.Port); err != nil {
			utils.Fatal("Failed to start server")
		}
	}()

	<-quit
	utils.Info("Shutting down server...")

	// 关闭数据库连接
	if err := database.CloseDB(); err != nil {
		utils.Error("Failed to close database connection")
	}

	utils.Info("Server gracefully stopped")
}
