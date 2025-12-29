package response

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Response represents a standard API response
type Response struct {
	Success bool       `json:"success"`
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
}

// ErrorInfo represents error details
type ErrorInfo struct {
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// PaginatedResponse represents a paginated API response (Laravel-style)
type PaginatedResponse struct {
	Success bool       `json:"success"`
	Code    int        `json:"code"`
	Status  string     `json:"status"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
	Links   []LinkItem `json:"links,omitempty"`
}

// LinkItem represents a single pagination link
type LinkItem struct {
	URL    *string `json:"url"`
	Label  string  `json:"label"`
	Active bool    `json:"active"`
}

// Meta holds pagination metadata (Laravel-style)
type Meta struct {
	CurrentPage  int        `json:"currentPage"`
	From         *int       `json:"from"`
	LastPage     int        `json:"lastPage"`
	Links        []LinkItem `json:"links,omitempty"`
	Path         string     `json:"path"`
	PerPage      int        `json:"perPage"`
	To           *int       `json:"to"`
	Total        int64      `json:"total"`
	FirstPageURL string     `json:"firstPageUrl"`
	LastPageURL  string     `json:"lastPageUrl"`
	NextPageURL  *string    `json:"nextPageUrl"`
	PrevPageURL  *string    `json:"prevPageUrl"`
}

// Pagination is kept for backward compatibility but deprecated
type Pagination = Meta

// Success sends a successful response
func Success(c *fiber.Ctx, statusCode int, message string, data any) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Code:    statusCode,
		Status:  getStatus(statusCode),
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, statusCode int, message string, errCode string, details string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Code:    statusCode,
		Status:  getStatus(statusCode),
		Message: message,
		Error: &ErrorInfo{
			Code:    errCode,
			Details: details,
		},
	})
}

func getStatus(code int) string {
	switch code {
	case fiber.StatusOK:
		return "OK"
	case fiber.StatusCreated:
		return "CREATED"
	case fiber.StatusNoContent:
		return "NO_CONTENT"
	case fiber.StatusBadRequest:
		return "BAD_REQUEST"
	case fiber.StatusUnauthorized:
		return "UNAUTHORIZED"
	case fiber.StatusForbidden:
		return "FORBIDDEN"
	case fiber.StatusNotFound:
		return "NOT_FOUND"
	case fiber.StatusConflict:
		return "CONFLICT"
	case fiber.StatusUnprocessableEntity:
		return "UNPROCESSABLE_ENTITY"
	case fiber.StatusInternalServerError:
		return "INTERNAL_SERVER_ERROR"
	default:
		return "UNKNOWN"
	}
}

// Paginated sends a paginated response (Laravel-style)
func Paginated(c *fiber.Ctx, statusCode int, message string, data any, page, perPage int, total int64) error {
	lastPage := int(total) / perPage
	if int(total)%perPage > 0 {
		lastPage++
	}

	fromVal := (page-1)*perPage + 1
	toVal := page * perPage
	if int64(toVal) > total {
		toVal = int(total)
	}

	var from, to *int
	if total > 0 {
		from = &fromVal
		to = &toVal
	}

	baseUrl := c.BaseURL() + c.Path()

	// Create query string without page
	queries := c.Queries()
	delete(queries, "page")

	queryString := ""
	for k, v := range queries {
		if queryString == "" {
			queryString = "?"
		} else {
			queryString += "&"
		}
		queryString += k + "=" + v
	}

	glue := "?"
	if queryString != "" {
		glue = "&"
	}

	generateUrl := func(p int) string {
		return baseUrl + queryString + glue + "page=" + fmt.Sprintf("%d", p)
	}

	firstPageUrl := generateUrl(1)
	lastPageUrl := generateUrl(lastPage)

	var nextPageUrl, prevPageUrl *string
	if page < lastPage {
		u := generateUrl(page + 1)
		nextPageUrl = &u
	}
	if page > 1 {
		u := generateUrl(page - 1)
		prevPageUrl = &u
	}

	links := []LinkItem{}
	// Previous link
	links = append(links, LinkItem{
		URL:    prevPageUrl,
		Label:  "&laquo; Previous",
		Active: false,
	})

	// Page links (simplified: just show all for now, or could be optimized)
	for i := 1; i <= lastPage; i++ {
		u := generateUrl(i)
		links = append(links, LinkItem{
			URL:    &u,
			Label:  fmt.Sprintf("%d", i),
			Active: i == page,
		})
	}

	// Next link
	links = append(links, LinkItem{
		URL:    nextPageUrl,
		Label:  "Next &raquo;",
		Active: false,
	})

	meta := &Meta{
		CurrentPage:  page,
		From:         from,
		LastPage:     lastPage,
		Links:        links,
		Path:         baseUrl,
		PerPage:      perPage,
		To:           to,
		Total:        total,
		FirstPageURL: firstPageUrl,
		LastPageURL:  lastPageUrl,
		NextPageURL:  nextPageUrl,
		PrevPageURL:  prevPageUrl,
	}

	return c.Status(statusCode).JSON(PaginatedResponse{
		Success: true,
		Code:    statusCode,
		Status:  getStatus(statusCode),
		Message: message,
		Data:    data,
		Meta:    meta,
		Links:   links,
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
