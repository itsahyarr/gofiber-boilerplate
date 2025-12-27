package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
)

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get a user by their ID (ADMIN only)
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} response.Response{data=dto.UserResponse}
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.GetByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		return response.InternalServerError(c, "failed to get user")
	}

	return response.Success(c, fiber.StatusOK, "user retrieved successfully", user)
}
