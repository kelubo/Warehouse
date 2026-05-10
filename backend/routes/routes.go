package routes

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"warehouse-management/backend/controllers"
	"warehouse-management/backend/middleware"
)

func getStaticDir() string {
	exeDir := filepath.Dir(os.Args[0])
	if exeDir == "." || exeDir == "" {
		exeDir, _ = os.Getwd()
	}
	staticDir := filepath.Join(exeDir, "static")
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		exeDir, _ = os.Getwd()
		staticDir = filepath.Join(exeDir, "static")
	}
	return staticDir
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	staticDir := getStaticDir()

	r.Static("/static", staticDir)
	r.GET("/", func(c *gin.Context) {
		indexPath := filepath.Join(staticDir, "index.html")
		c.File(indexPath)
	})

	RegisterRoutes(r)

	return r
}

func RegisterRoutes(r *gin.Engine) {
	staticDir := getStaticDir()

	r.Static("/static", staticDir)
	r.GET("/", func(c *gin.Context) {
		indexPath := filepath.Join(staticDir, "index.html")
		c.File(indexPath)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 认证路由
	auth := r.Group("/api/v1")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh-token", controllers.RefreshToken)
	}

	// 需要认证的路由
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// 产品/物品管理
		api.GET("/products", controllers.GetProducts)
		api.GET("/products/:id", controllers.GetProduct)
		api.POST("/products", controllers.CreateProduct)
		api.PUT("/products/:id", controllers.UpdateProduct)
		api.DELETE("/products/:id", controllers.DeleteProduct)

		// 仓库管理
		api.GET("/warehouses", controllers.GetWarehouses)
		api.GET("/warehouses/:id", controllers.GetWarehouse)
		api.POST("/warehouses", controllers.CreateWarehouse)
		api.PUT("/warehouses/:id", controllers.UpdateWarehouse)
		api.DELETE("/warehouses/:id", controllers.DeleteWarehouse)

		// 货架管理（通过查询参数过滤仓库）
		api.GET("/shelves", controllers.GetShelves)
		api.GET("/shelves/:id", controllers.GetShelf)
		api.POST("/shelves", controllers.CreateShelf)
		api.PUT("/shelves/:id", controllers.UpdateShelf)
		api.DELETE("/shelves/:id", controllers.DeleteShelf)

		// 箱子管理（通过查询参数过滤仓库）
		api.GET("/boxes", controllers.GetBoxes)
		api.GET("/boxes/:id", controllers.GetBox)
		api.POST("/boxes", controllers.CreateBox)
		api.PUT("/boxes/:id", controllers.UpdateBox)
		api.DELETE("/boxes/:id", controllers.DeleteBox)

		// 库存管理
		api.GET("/inventories", controllers.GetInventories)
		api.GET("/inventories/:id", controllers.GetInventory)
		api.POST("/inventories", controllers.CreateInventory)
		api.PUT("/inventories/:id", controllers.UpdateInventory)
		api.DELETE("/inventories/:id", controllers.DeleteInventory)

		// 数据统计
		api.GET("/stats", controllers.GetStats)
		api.GET("/stats/warehouses", controllers.GetWarehouseStats)
		api.GET("/stats/categories", controllers.GetCategoryStats)
		api.GET("/stats/activities", controllers.GetRecentActivities)
	}

	// 管理员专属路由
	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// 用户管理
		admin.GET("/users", controllers.GetUsers)
		admin.GET("/users/:id", controllers.GetUser)
		admin.POST("/users", controllers.CreateUser)
		admin.PUT("/users/:id", controllers.UpdateUser)
		admin.DELETE("/users/:id", controllers.DeleteUser)
	}
}
