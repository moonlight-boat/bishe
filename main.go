package main

import (
	"job-backend/config"
	"job-backend/database"
	"job-backend/handlers"
	"job-backend/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.LoadConfig()
	
	// 初始化数据库
	if err := database.InitDatabase(cfg.DatabasePath); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.CloseDatabase()
	
	// 启动定时任务
	scheduler := services.NewScheduler(cfg)
	scheduler.Start()
	defer scheduler.Stop()
	
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)
	
	// 创建路由
	router := gin.Default()
	
	// 添加CORS中间件
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})
	
	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"message": "Job Backend Service is running",
		})
	})
	
	// API路由组
	api := router.Group("/api")
	{
		// 招聘信息相关路由
		jobs := api.Group("/jobs")
		{
			jobs.GET("", handlers.GetJobs)           // 获取招聘信息列表
			jobs.GET("/:id", handlers.GetJobByID)    // 获取单个招聘信息
			jobs.GET("/search", handlers.SearchJobs) // 搜索招聘信息
			jobs.GET("/stats", handlers.GetJobStats) // 获取统计信息
		}
		
		// 同步相关路由
		api.POST("/sync", handlers.ManualSync)       // 手动触发同步
		api.GET("/sync/status", handlers.SyncStatus) // 获取同步状态
	}
	
	log.Printf("服务器启动在端口 %s", cfg.ServerPort)
	log.Fatal(router.Run(":" + cfg.ServerPort))
}
