package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/dto"
	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/validator"
)

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Get new access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshTokenRequest true "Refresh token request"
// @Success      200 {object} response.Response{data=dto.TokenResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return response.BadRequest(c, "invalid request body", validator.FormatValidationErrors(err))
	}

	result, err := h.authService.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			return response.Unauthorized(c, "invalid or expired refresh token")
		}
		return response.InternalServerError(c, "failed to refresh token")
	}

	return response.Success(c, fiber.StatusOK, "tokens refreshed successfully", result)
}
