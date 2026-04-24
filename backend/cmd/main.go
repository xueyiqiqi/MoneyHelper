package main

import (
	"path/filepath"

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
	logger.Init()
	_ = config.LoadEnvFileIfPresent(filepath.Join(".", ".env"))
	cfg := config.Load()

	repository.InitDB(cfg.Database)
	permissionService := service.NewPermissionService(repository.DB)
	middleware.SetPermissionService(permissionService)
	repository.InitRedis(cfg.Redis.Addr)

	userRepo := &repository.UserRepository{}
	spaceRepo := &repository.SpaceRepository{}
	billRepo := &repository.BillRepository{}
	reportRepo := &repository.ReportRepository{}

	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	spaceService := &service.SpaceService{SpaceRepo: spaceRepo, UserRepo: userRepo}
	middleware.SetJWTSecret(cfg.JWT.Secret)

	aiAgent := ai.NewAIAgent(ai.Config{
		BaseURL:        cfg.AI.BaseURL,
		APIKey:         cfg.AI.APIKey,
		Model:          cfg.AI.Model,
		TimeoutSeconds: cfg.AI.TimeoutSeconds,
	})
	billService := &service.BillService{
		BillRepo:   billRepo,
		SpaceRepo:  spaceRepo,
		ReportRepo: reportRepo,
		AIAgent:    aiAgent,
	}

	authHandler := &handler.AuthHandler{AuthService: authService}
	spaceHandler := &handler.SpaceHandler{SpaceService: spaceService}
	billHandler := &handler.BillHandler{BillService: billService}

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", authHandler.Register)
		v1.POST("/token", authHandler.Login)
		v1.POST("/refresh", authHandler.RefreshToken)

		auth := v1.Group("/")
		auth.Use(middleware.JWTMiddleware())
		{
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/spaces", spaceHandler.CreateSpace)
			auth.GET("/spaces", spaceHandler.GetUserSpaces)
			auth.POST("/spaces/:id/members", middleware.PermissionMiddleware("space", "admin"), spaceHandler.AddMember)
			auth.GET("/spaces/:id/members", middleware.PermissionMiddleware("space", "read"), spaceHandler.GetMembers)
			auth.PUT("/spaces/:id/members/:userId", middleware.PermissionMiddleware("space", "admin"), spaceHandler.UpdateMemberRole)
			auth.DELETE("/spaces/:id/members/:userId", middleware.PermissionMiddleware("space", "admin"), spaceHandler.RemoveMember)
			auth.POST("/spaces/:id/leave", spaceHandler.LeaveSpace)

			auth.POST("/bills", middleware.PermissionMiddleware("bill", "write"), billHandler.CreateBill)
			auth.GET("/bills", middleware.PermissionMiddleware("bill", "read"), billHandler.GetBills)
			auth.GET("/bills/:id", middleware.BillPermissionMiddleware("read"), billHandler.GetBill)
			auth.PUT("/bills/:id", middleware.BillPermissionMiddleware("write"), billHandler.UpdateBill)

			auth.POST("/analyze", middleware.PermissionMiddleware("bill", "read"), billHandler.GenerateAnalysis)
			auth.GET("/reports", middleware.PermissionMiddleware("bill", "read"), billHandler.GetReports)
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	logger.Info.Printf("Server starting on :%s...", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Error.Fatalf("failed to run server: %v", err)
	}
}
