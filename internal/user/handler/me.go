package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/middleware"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
)

// GetCurrentUser godoc
// @Summary      Get current user
// @Description  Get the profile of the currently authenticated user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.Response{data=dto.UserResponse}
// @Failure      401 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /users/me [get]
func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	payload := middleware.GetAuthPayload(c)
	if payload == nil {
		return response.Unauthorized(c, "authentication required")
	}

	user, err := h.userService.GetByID(c.Context(), payload.UserID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		return response.InternalServerError(c, "failed to get user")
	}

	return response.Success(c, fiber.StatusOK, "user retrieved successfully", user)
}
