package response

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo represents error details
type ErrorInfo struct {
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated API response (Laravel-style)
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Links   *Links      `json:"links,omitempty"`
}

// Meta holds pagination metadata (Laravel-style)
type Meta struct {
	CurrentPage int   `json:"currentPage"`
	From        int   `json:"from"`
	LastPage    int   `json:"lastPage"`
	PerPage     int   `json:"perPage"`
	To          int   `json:"to"`
	Total       int64 `json:"total"`
}

// Links holds pagination links (Laravel-style)
type Links struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
}

// Pagination is kept for backward compatibility but deprecated
type Pagination = Meta

// Success sends a successful response
func Success(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, statusCode int, message string, errCode string, details string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    errCode,
			Details: details,
		},
	})
}

// Paginated sends a paginated response (Laravel-style)
func Paginated(c *fiber.Ctx, statusCode int, message string, data interface{}, page, perPage int, total int64) error {
	lastPage := int(total) / perPage
	if int(total)%perPage > 0 {
		lastPage++
	}

	from := (page-1)*perPage + 1
	to := page * perPage
	if int64(to) > total {
		to = int(total)
	}
	if total == 0 {
		from = 0
		to = 0
	}

	return c.Status(statusCode).JSON(PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: &Meta{
			CurrentPage: page,
			From:        from,
			LastPage:    lastPage,
			PerPage:     perPage,
			To:          to,
			Total:       total,
		},
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *fiber.Ctx, message string, details string) error {
	return Error(c, fiber.StatusBadRequest, message, "BAD_REQUEST", details)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message, "UNAUTHORIZED", "")
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusForbidden, message, "FORBIDDEN", "")
}

// NotFound sends a 404 Not Found response
func NotFound(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, "NOT_FOUND", "")
}

// InternalServerError sends a 500 Internal Server Error response
func InternalServerError(c *fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message, "INTERNAL_ERROR", "")
}

// Conflict sends a 409 Conflict response
func Conflict(c *fiber.Ctx, message string, details string) error {
	return Error(c, fiber.StatusConflict, message, "CONFLICT", details)
}
