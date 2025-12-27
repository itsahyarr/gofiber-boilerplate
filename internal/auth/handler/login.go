package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/dto"
	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/validator"
)

// Login godoc
// @Summary      Login user
// @Description  Login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login request"
// @Success      200 {object} response.Response{data=dto.AuthResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return response.BadRequest(c, "invalid request body", validator.FormatValidationErrors(err))
	}

	result, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return response.Unauthorized(c, "invalid email or password")
		}
		if errors.Is(err, service.ErrUserNotActive) {
			return response.Forbidden(c, "user account is not active")
		}
		return response.InternalServerError(c, "failed to login")
	}

	return response.Success(c, fiber.StatusOK, "login successful", result)
}
