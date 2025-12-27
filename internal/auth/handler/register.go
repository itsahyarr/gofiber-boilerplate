package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/dto"
	"github.com/itsahyarr/gofiber-boilerplate/internal/auth/service"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/validator"
)

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "Register request"
// @Success      201 {object} response.Response{data=dto.AuthResponse}
// @Failure      400 {object} response.Response
// @Failure      409 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := validator.ParseAndValidate(c, &req); err != nil {
		return response.BadRequest(c, "invalid request body", validator.FormatValidationErrors(err))
	}

	result, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return response.Conflict(c, "email already exists", "")
		}
		return response.InternalServerError(c, "failed to register user")
	}

	return response.Success(c, fiber.StatusCreated, "user registered successfully", result)
}
