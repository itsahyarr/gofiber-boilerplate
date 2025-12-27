package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/middleware"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
)

// Logout godoc
// @Summary      Logout user
// @Description  Logout and invalidate refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	payload := middleware.GetAuthPayload(c)
	if payload == nil {
		return response.Unauthorized(c, "authentication required")
	}

	if err := h.authService.Logout(c.Context(), payload.UserID); err != nil {
		return response.InternalServerError(c, "failed to logout")
	}

	return response.Success(c, fiber.StatusOK, "logout successful", nil)
}
