package auth

import (
	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/handler"
	"github.com/itsahyarr/gofiber-boilerplate/internal/middleware"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/token"
)

// RegisterRoutes registers all auth routes
func RegisterRoutes(router fiber.Router, h *handler.AuthHandler, tokenMaker *token.PasetoMaker) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.RefreshToken)

	// Protected routes
	authProtected := auth.Group("", middleware.AuthMiddleware(tokenMaker))
	authProtected.Post("/logout", h.Logout)
}
