package middleware

import (
	"slices"
	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
	"github.com/itsahyarr/gofiber-boilerplate/shared/entity"
)

// RequireRoles creates a role-based access control middleware
func RequireRoles(allowedRoles ...entity.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		payload := GetAuthPayload(c)
		if payload == nil {
			return response.Unauthorized(c, "authentication required")
		}

		userRole := entity.Role(payload.Role)
		if slices.Contains(allowedRoles, userRole) {
				return c.Next()
			}

		return response.Forbidden(c, "insufficient permissions")
	}
}

// RequireAdmin is a convenience middleware that requires ADMIN role
func RequireAdmin() fiber.Handler {
	return RequireRoles(entity.RoleAdmin)
}

// RequireUser is a convenience middleware that requires USER or ADMIN role
func RequireUser() fiber.Handler {
	return RequireRoles(entity.RoleUser, entity.RoleAdmin)
}
