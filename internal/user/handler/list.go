package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
)

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Get a paginated list of all users (ADMIN only)
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        perPage query int false "Items per page" default(10)
// @Success      200 {object} response.PaginatedResponse{data=[]dto.UserResponse}
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("perPage", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	users, total, err := h.userService.GetAll(c.Context(), page, perPage)
	if err != nil {
		return response.InternalServerError(c, "failed to get users")
	}

	return response.Paginated(c, fiber.StatusOK, "users retrieved successfully", users, page, perPage, total)
}

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
