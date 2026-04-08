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
	repository.InitDB(cfg.Database.Path)

	// Initialize Repositories
	userRepo := &repository.UserRepository{}
	spaceRepo := &repository.SpaceRepository{}
	billRepo := &repository.BillRepository{}
	reportRepo := &repository.ReportRepository{}

	// Initialize Services
	authService := &service.AuthService{UserRepo: userRepo}
	spaceService := &service.SpaceService{SpaceRepo: spaceRepo}

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

	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware())
	{
		// Space routes
		auth.POST("/spaces", spaceHandler.CreateSpace)
		auth.GET("/spaces", spaceHandler.GetUserSpaces)

		// Bill routes
		auth.POST("/bills", billHandler.CreateBill)
		auth.GET("/bills", billHandler.GetBills)

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
