package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
)

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete a user (ADMIN only)
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.userService.Delete(c.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		return response.InternalServerError(c, "failed to delete user")
	}

	return response.Success(c, fiber.StatusOK, "user deleted successfully", nil)
}
