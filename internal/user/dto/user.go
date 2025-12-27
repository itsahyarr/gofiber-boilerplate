package dto

import (
	"github.com/itsahyarr/gofiber-boilerplate/pkg/utils"
	"github.com/itsahyarr/gofiber-boilerplate/shared/entity"
)

// UserResponse represents the user response
type UserResponse struct {
	ID        string      `json:"id"`
	Email     string      `json:"email"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Role      entity.Role `json:"role"`
	IsActive  bool        `json:"isActive"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt"`
}

// ToUserResponse converts a User entity to UserResponse DTO
func ToUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: utils.FormatIndonesian(user.CreatedAt),
		UpdatedAt: utils.FormatIndonesian(user.UpdatedAt),
	}
}

// ToUserResponses converts a slice of User entities to UserResponse DTOs
func ToUserResponses(users []*entity.User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = ToUserResponse(user)
	}
	return responses
}

// UpdateUserRequest represents the update user request body
type UpdateUserRequest struct {
	FirstName *string      `json:"firstName,omitempty" validate:"omitempty,min=2"`
	LastName  *string      `json:"lastName,omitempty" validate:"omitempty,min=2"`
	Role      *entity.Role `json:"role,omitempty"`
	IsActive  *bool        `json:"isActive,omitempty"`
}

// ChangePasswordRequest represents the change password request body
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}
