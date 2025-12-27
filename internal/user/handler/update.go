package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/middleware"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user/dto"
	"github.com/itsahyarr/gofiber-boilerplate/internal/user/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/validator"
	"github.com/itsahyarr/gofiber-boilerplate/shared/entity"
)

// UpdateUser godoc
// @Summary      Update user
// @Description  Update a user's profile (ADMIN can update any user, users can update themselves)
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Param        request body dto.UpdateUserRequest true "Update request"
// @Success      200 {object} response.Response{data=dto.UserResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	payload := middleware.GetAuthPayload(c)
	if payload == nil {
		return response.Unauthorized(c, "authentication required")
	}

	// Check if user is updating themselves or is an admin
	if payload.UserID != id && entity.Role(payload.Role) != entity.RoleAdmin {
		return response.Forbidden(c, "you can only update your own profile")
	}

	var req dto.UpdateUserRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return response.BadRequest(c, "invalid request body", validator.FormatValidationErrors(err))
	}

	// Only ADMIN can update role
	if req.Role != nil && entity.Role(payload.Role) != entity.RoleAdmin {
		return response.Forbidden(c, "only admins can update roles")
	}

	user, err := h.userService.Update(c.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		return response.InternalServerError(c, "failed to update user")
	}

	return response.Success(c, fiber.StatusOK, "user updated successfully", user)
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change user's password
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.ChangePasswordRequest true "Change password request"
// @Success      200 {object} response.Response
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /users/me/password [put]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	payload := middleware.GetAuthPayload(c)
	if payload == nil {
		return response.Unauthorized(c, "authentication required")
	}

	var req dto.ChangePasswordRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return response.BadRequest(c, "invalid request body", validator.FormatValidationErrors(err))
	}

	if err := h.userService.ChangePassword(c.Context(), payload.UserID, &req); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return response.NotFound(c, "user not found")
		}
		if errors.Is(err, service.ErrInvalidOldPassword) {
			return response.BadRequest(c, "invalid old password", "")
		}
		return response.InternalServerError(c, "failed to change password")
	}

	return response.Success(c, fiber.StatusOK, "password changed successfully", nil)
}
