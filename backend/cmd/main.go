package main

import (
	"life-financial-assistant-backend/internal/api/handler"
	"life-financial-assistant-backend/internal/api/middleware"
	"life-financial-assistant-backend/internal/config"
	"life-financial-assistant-backend/internal/logger"
	"life-financial-assistant-backend/internal/repository"
	"life-financial-assistant-backend/internal/service"
	"life-financial-assistant-backend/pkg/ai"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.Init()

	// Load configuration
	cfg := config.Load()

	// Initialize Database
	repository.InitDB(cfg.Database)

	// 初始化权限服务
	permissionService := service.NewPermissionService(repository.DB)
	middleware.SetPermissionService(permissionService)

	// Initialize Redis
	repository.InitRedis(cfg.Redis.Addr)

	// Initialize Repositories
	userRepo := &repository.UserRepository{}
	spaceRepo := &repository.SpaceRepository{}
	billRepo := &repository.BillRepository{}
	reportRepo := &repository.ReportRepository{}

	// Initialize Services
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	spaceService := &service.SpaceService{SpaceRepo: spaceRepo, UserRepo: userRepo}

	// Set JWT secret for middleware
	middleware.SetJWTSecret(cfg.JWT.Secret)

	// AIAgent
	aiAgent := &ai.AIAgent{APIKey: cfg.AI.APIKey}

	billService := &service.BillService{
		BillRepo:   billRepo,
		SpaceRepo:  spaceRepo,
		ReportRepo: reportRepo,
		AIAgent:    aiAgent,
	}

	// Initialize Handlers
	authHandler := &handler.AuthHandler{AuthService: authService}
	spaceHandler := &handler.SpaceHandler{SpaceService: spaceService}
	billHandler := &handler.BillHandler{BillService: billService}

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Auth routes
	r.POST("/register", authHandler.Register)
	r.POST("/token", authHandler.Login)
	r.POST("/refresh", authHandler.RefreshToken)

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware())
	{
		// Space routes
		auth.POST("/spaces", spaceHandler.CreateSpace)
		auth.GET("/spaces", spaceHandler.GetUserSpaces)
		auth.POST("/spaces/:id/members", middleware.PermissionMiddleware("space", "admin"), spaceHandler.AddMember)
		auth.POST("/spaces/:id/leave", spaceHandler.LeaveSpace)

		// Bill routes
		auth.POST("/bills", middleware.PermissionMiddleware("bill", "write"), billHandler.CreateBill)
		auth.GET("/bills", middleware.PermissionMiddleware("bill", "read"), billHandler.GetBills)

		// AI Analysis routes
		auth.POST("/analyze", billHandler.GenerateAnalysis)
		auth.GET("/reports", billHandler.GetReports)
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	logger.Info.Printf("Server starting on :%s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Error.Fatalf("failed to run server: %v", err)
	}
}
