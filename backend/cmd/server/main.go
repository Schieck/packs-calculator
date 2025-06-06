package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Schieck/packs-calculator/internal/adapter/config"
	httpAdapter "github.com/Schieck/packs-calculator/internal/adapter/http"
	"github.com/Schieck/packs-calculator/internal/adapter/repository"

	authService "github.com/Schieck/packs-calculator/internal/service/auth"
	healthService "github.com/Schieck/packs-calculator/internal/service/health"
	packCalculatorService "github.com/Schieck/packs-calculator/internal/service/pack_calculator"
	packConfigurationService "github.com/Schieck/packs-calculator/internal/service/pack_configuration"

	authUseCase "github.com/Schieck/packs-calculator/internal/usecase/auth"
	healthUseCase "github.com/Schieck/packs-calculator/internal/usecase/health"
	packCalculatorUseCase "github.com/Schieck/packs-calculator/internal/usecase/pack_calculator"
	packConfigurationUseCase "github.com/Schieck/packs-calculator/internal/usecase/pack_configuration"

	"github.com/Schieck/packs-calculator/pkg/db"
	"github.com/Schieck/packs-calculator/pkg/middleware"

	"github.com/Schieck/packs-calculator/docs"
)

// @title Packs Calculator API
// @version 1.0
// @description API for calculating packs and managing pack configurations with simple JWT authentication
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Packs Calculator API server")

	// Load configuration
	cfg := config.LoadConfig()

	// Configure Swagger host dynamically
	configureSwaggerHost(cfg.Server.Port)

	// Initialize database
	database, err := db.NewConnection(cfg.Database.DSN)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Initialize repositories
	packConfigRepo := repository.NewPackConfigurationRepository(database.DB)

	// Initialize services
	authSvc := authService.NewAuthServiceWithDefaults(cfg.Auth.JWTSecret, cfg.Auth.AuthSecret)
	healthSvc := healthService.NewHealthService(database, "1.0.0")
	packCalculatorSvc := packCalculatorService.NewPackCalculatorService()
	packSizeProcessorSvc := packCalculatorService.NewPackSizeProcessorService()
	packConfigSvc := packConfigurationService.NewPackConfigurationService(packConfigRepo)

	// Initialize use cases
	authenticateUseCase := authUseCase.NewAuthenticateUseCase(authSvc, logger)
	validateTokenUseCase := authUseCase.NewValidateTokenUseCase(authSvc, logger)
	healthCheckUseCase := healthUseCase.NewHealthUseCase(healthSvc, logger)
	calculatePacksUseCase := packCalculatorUseCase.NewCalculatePacksUseCase(packCalculatorSvc, packSizeProcessorSvc, logger)

	// Pack configuration use cases
	getAllConfigurationsUseCase := packConfigurationUseCase.NewGetAllConfigurationsUseCase(packConfigSvc, logger)
	getConfigurationByIDUseCase := packConfigurationUseCase.NewGetConfigurationByIDUseCase(packConfigSvc, logger)
	getDefaultConfigurationUseCase := packConfigurationUseCase.NewGetDefaultConfigurationUseCase(packConfigSvc, logger)
	createConfigurationUseCase := packConfigurationUseCase.NewCreateConfigurationUseCase(packConfigSvc, logger)
	updateConfigurationUseCase := packConfigurationUseCase.NewUpdateConfigurationUseCase(packConfigSvc, logger)
	deleteConfigurationUseCase := packConfigurationUseCase.NewDeleteConfigurationUseCase(packConfigSvc, logger)
	setDefaultConfigurationUseCase := packConfigurationUseCase.NewSetDefaultConfigurationUseCase(packConfigSvc, logger)

	// Initialize HTTP handlers
	authHandler := httpAdapter.NewAuthHandler(authenticateUseCase, logger)
	healthHandler := httpAdapter.NewHealthHandler(healthCheckUseCase, logger)
	packCalculatorHandler := httpAdapter.NewCalculatorHandler(calculatePacksUseCase, logger)
	packConfigHandler := httpAdapter.NewPackConfigurationHandler(
		getAllConfigurationsUseCase,
		getConfigurationByIDUseCase,
		getDefaultConfigurationUseCase,
		createConfigurationUseCase,
		updateConfigurationUseCase,
		deleteConfigurationUseCase,
		setDefaultConfigurationUseCase,
		logger,
	)

	// Setup Gin
	if gin.Mode() == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := router.Group("/api/v1")

	v1.GET("/health", healthHandler.Health)

	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/token", authHandler.Authenticate)
	}

	protected := v1.Group("/")
	protected.Use(middleware.JWT(validateTokenUseCase))
	{
		protected.POST("/calculate", packCalculatorHandler.Calculate)

		protected.GET("/pack-configurations", packConfigHandler.GetAllConfigurations)
		protected.GET("/pack-configurations/default", packConfigHandler.GetDefaultConfiguration)
		protected.GET("/pack-configurations/:id", packConfigHandler.GetConfigurationByID)
		protected.POST("/pack-configurations", packConfigHandler.CreateConfiguration)
		protected.PUT("/pack-configurations/:id", packConfigHandler.UpdateConfiguration)
		protected.DELETE("/pack-configurations/:id", packConfigHandler.DeleteConfiguration)
		protected.PATCH("/pack-configurations/:id/default", packConfigHandler.SetDefaultConfiguration)
	}

	// Setup HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		slog.Info("Server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server exited")
}

// configureSwaggerHost dynamically sets the Swagger host based on the environment
func configureSwaggerHost(port string) {
	// Check if we're running on fly.io by looking for FLY_APP_NAME environment variable
	flyAppName := os.Getenv("FLY_APP_NAME")

	if flyAppName != "" {
		// Running on fly.io - use the fly.io hostname
		docs.SwaggerInfo.Host = flyAppName + ".fly.dev"
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		// Running locally or in other environments
		hostname := os.Getenv("SWAGGER_HOST")
		if hostname == "" {
			// Default to localhost with the configured port
			hostname = "localhost:" + port
		}

		// Determine scheme based on hostname
		if strings.Contains(hostname, "localhost") || strings.Contains(hostname, "127.0.0.1") {
			docs.SwaggerInfo.Host = hostname
			docs.SwaggerInfo.Schemes = []string{"http"}
		} else {
			docs.SwaggerInfo.Host = hostname
			docs.SwaggerInfo.Schemes = []string{"https", "http"}
		}
	}

	slog.Info("Swagger configured", "host", docs.SwaggerInfo.Host, "schemes", docs.SwaggerInfo.Schemes)
}
