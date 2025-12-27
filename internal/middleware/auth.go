package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/itsahyarr/gofiber-boilerplate/pkg/response"
	"github.com/itsahyarr/gofiber-boilerplate/pkg/token"
)

const (
	AuthorizationHeader     = "Authorization"
	AuthorizationTypeBearer = "Bearer"
	AuthPayloadKey          = "auth_payload"
)

// AuthMiddleware creates an authentication middleware
func AuthMiddleware(tokenMaker *token.PasetoMaker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(AuthorizationHeader)
		if authHeader == "" {
			return response.Unauthorized(c, "authorization header is required")
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			return response.Unauthorized(c, "invalid authorization header format")
		}

		authType := fields[0]
		if !strings.EqualFold(authType, AuthorizationTypeBearer) {
			return response.Unauthorized(c, "unsupported authorization type")
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			if err == token.ErrExpiredToken {
				return response.Unauthorized(c, "access token has expired")
			}
			return response.Unauthorized(c, "invalid access token")
		}

		// Check if it's an access token
		if payload.TokenType != "access" {
			return response.Unauthorized(c, "invalid token type")
		}

		// Store payload in context
		c.Locals(AuthPayloadKey, payload)
		return c.Next()
	}
}

// GetAuthPayload retrieves the auth payload from context
func GetAuthPayload(c *fiber.Ctx) *token.Payload {
	payload, ok := c.Locals(AuthPayloadKey).(*token.Payload)
	if !ok {
		return nil
	}
	return payload
}
