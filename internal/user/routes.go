package user

import (
	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/middleware"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user/handler"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/token"
	"github.com/itsahyarr/gofiber-boilerplate/shared/entity"
)

// RegisterRoutes registers all user routes
func RegisterRoutes(router fiber.Router, h *handler.UserHandler, tokenMaker *token.PasetoMaker) {
	users := router.Group("/users", middleware.AuthMiddleware(tokenMaker))

	// User routes (authenticated)
	users.Get("/me", h.GetCurrentUser)
	users.Put("/me/password", h.ChangePassword)

	// Admin-only routes
	adminUsers := users.Group("", middleware.RequireRoles(entity.RoleAdmin))
	adminUsers.Get("/", h.GetAllUsers)
	adminUsers.Get("/:id", h.GetUserByID)
	adminUsers.Delete("/:id", h.DeleteUser)

	// User update - allow self-update or admin
	users.Put("/:id", h.UpdateUser)
}
