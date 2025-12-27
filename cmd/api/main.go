package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"

	"github.com/itsahyarr/gofiber-boilerplate/internal/auth"
	authHandler "github.com/itsahyarr/gofiber-boilerplate/internal/auth/handler"
	authRepo "github.com/itsahyarr/gofiber-boilerplate/internal/auth/repository"
	authService "github.com/itsahyarr/gofiber-boilerplate/internal/auth/service"
	"github.com/itsahyarr/gofiber-boilerplate/internal/config"
	"github.com/itsahyarr/gofiber-boilerplate/internal/database/migration"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user"
	userHandler "github.com/itsahyarr/gofiber-boilerplate/internal/user/handler"
	userRepo "github.com/itsahyarr/gofiber-boilerplate/internal/user/repository"
	userService "github.com/itsahyarr/gofiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/database"
	pkgLogger "github.com/itsahyarr/gofiber-boilerplate/pkg/logger"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/token"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	pkgLogger.Init(cfg.App.LogLevel, cfg.App.Environment)
	defer pkgLogger.Sync()

	pkgLogger.Info("Starting application",
		zap.String("environment", cfg.App.Environment),
		zap.String("port", cfg.Server.Port),
	)

	// Connect to MongoDB
	mongodb, err := database.NewMongoDB(cfg.Database.URI, cfg.Database.Database)
	if err != nil {
		pkgLogger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}

	// Connect to Redis
	redis, err := database.NewRedis(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		pkgLogger.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	// Internal database migrations
	migration.RunMigrations(mongodb)

	// Initialize PASETO token maker
	tokenMaker, err := token.NewPasetoMaker(cfg.Token.SymmetricKey)
	if err != nil {
		pkgLogger.Fatal("Failed to create token maker", zap.Error(err))
	}

	// Initialize repositories
	userRepository := userRepo.NewUserRepository(mongodb)
	tokenRepository := authRepo.NewTokenRepository(redis)

	// Initialize services
	authSvc := authService.NewAuthService(userRepository, tokenRepository, tokenMaker, cfg)
	userSvc := userService.NewUserService(userRepository, mongodb)

	// Initialize handlers
	authHdl := authHandler.NewAuthHandler(authSvc)
	userHdl := userHandler.NewUserHandler(userSvc)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:       "GoFiber Boilerplate",
		ServerHeader:  "Fiber",
		ReadTimeout:   cfg.Server.ReadTimeout,
		WriteTimeout:  cfg.Server.WriteTimeout,
		IdleTimeout:   cfg.Server.IdleTimeout,
		BodyLimit:     10 * 1024 * 1024, // 10 MB
		Prefork:       false,
		StrictRouting: true,
		CaseSensitive: true,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			pkgLogger.Error("Request error", zap.Error(err), zap.Int("status", code))
			return response.Error(c, code, "internal server error", "INTERNAL_ERROR", err.Error())
		},
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: false,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Register feature routes
	auth.RegisterRoutes(api, authHdl, tokenMaker)
	user.RegisterRoutes(api, userHdl, tokenMaker)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := app.Listen(addr); err != nil {
			pkgLogger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	pkgLogger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown Fiber app
	if err := app.ShutdownWithContext(ctx); err != nil {
		pkgLogger.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close database connections
	if err := mongodb.Close(ctx); err != nil {
		pkgLogger.Error("Failed to close MongoDB connection", zap.Error(err))
	}

	if err := redis.Close(); err != nil {
		pkgLogger.Error("Failed to close Redis connection", zap.Error(err))
	}

	pkgLogger.Info("Server shutdown complete")
}
