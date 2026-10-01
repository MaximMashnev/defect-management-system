package main

import (
	"defect-management-system/backend/internal/config"
	"defect-management-system/backend/internal/handler"
	"defect-management-system/backend/internal/middleware"
	"defect-management-system/backend/internal/repository"
	"defect-management-system/backend/internal/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting Defect Management API with driver [%s] on port :%s", cfg.DBDriver, cfg.Port)

	db, err := repository.NewDB(cfg.DBDriver, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	buildingRepo := repository.NewBuildingRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	uploadRepo := repository.NewUploadRepository(db)

	// Services
	priorityService := service.NewPriorityService(appRepo, buildingRepo, catRepo, uploadRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(userRepo, cfg.JWTSecret)
	buildingHandler := handler.NewBuildingHandler(buildingRepo)
	categoryHandler := handler.NewCategoryHandler(catRepo)
	appHandler := handler.NewApplicationHandler(appRepo, buildingRepo, priorityService)
	uploadHandler := handler.NewUploadHandler(uploadRepo, cfg.UploadDir, cfg.BaseURL)
	analyticsHandler := handler.NewAnalyticsHandler(appRepo, buildingRepo)

	// Gin Router
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "driver": cfg.DBDriver})
	})

	api := r.Group("/api")
	{
		// Auth routes (supports both Mokky legacy paths and standard REST)
		api.POST("/auth", authHandler.Login)
		api.POST("/auth/login", authHandler.Login)
		api.GET("/auth_me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.AuthMe)
		api.GET("/auth/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.AuthMe)

		// Categories
		api.GET("/category", categoryHandler.GetCategories)
		api.GET("/categories", categoryHandler.GetCategories)

		// Buildings
		api.GET("/buildings", buildingHandler.GetBuildings)
		api.GET("/buildings/:id", buildingHandler.GetBuilding)
		api.POST("/buildings", middleware.AuthMiddleware(cfg.JWTSecret), buildingHandler.CreateBuilding)
		api.PATCH("/buildings/:id", middleware.AuthMiddleware(cfg.JWTSecret), buildingHandler.UpdateBuilding)
		api.DELETE("/buildings/:id", middleware.AuthMiddleware(cfg.JWTSecret), buildingHandler.DeleteBuilding)

		// Applications
		api.GET("/applications", appHandler.GetApplications)
		api.GET("/applications/:id", appHandler.GetApplication)
		api.POST("/applications", appHandler.CreateApplication)
		api.PATCH("/applications/:id", appHandler.UpdateApplication)

		// Uploads
		api.POST("/uploads", uploadHandler.Upload)
		api.GET("/uploads", uploadHandler.GetFiles)
		api.GET("/uploads/files/:filename", uploadHandler.ServeFile)

		// Analytics
		api.GET("/analytics", analyticsHandler.GetAnalytics)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server listening on http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
