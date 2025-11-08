package middleware

import (
	"gofiber-template/pkg/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Protected middleware validates JWT tokens and sets user context
func Protected() fiber.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "Missing authorization header")
		}

		// Extract token from header
		token := utils.ExtractTokenFromHeader(authHeader)
		if token == "" {
			return utils.UnauthorizedResponse(c, "Invalid authorization header format")
		}

		// Validate token and get user context
		userCtx, err := utils.ValidateTokenStringToUUID(token, jwtSecret)
		if err != nil {
			log.Printf("❌ Token validation failed: %v", err)
			switch err {
			case utils.ErrExpiredToken:
				return utils.UnauthorizedResponse(c, "Token has expired")
			case utils.ErrInvalidToken:
				return utils.UnauthorizedResponse(c, "Invalid token")
			case utils.ErrMissingToken:
				return utils.UnauthorizedResponse(c, "Missing token")
			default:
				return utils.UnauthorizedResponse(c, "Token validation failed")
			}
		}

		log.Printf("✅ Token validated for user: %s (%s)", userCtx.Email, userCtx.ID)

		// Set user context in fiber locals
		c.Locals("user", userCtx)
		c.Locals("userID", userCtx.ID)

		return c.Next()
	}
}

// RequireRole middleware checks if user has specific role
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			return utils.UnauthorizedResponse(c, "User not authenticated")
		}

		if user.Role != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Insufficient permissions",
				"error":   "Access denied",
			})
		}

		return c.Next()
	}
}

// AdminOnly middleware ensures only admin users can access
func AdminOnly() fiber.Handler {
	return RequireRole("admin")
}

// OwnerOnly middleware checks if user is the owner of the resource
func OwnerOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			return utils.UnauthorizedResponse(c, "User not authenticated")
		}

		c.Locals("requireOwnership", true)
		c.Locals("ownerUserID", user.ID)

		return c.Next()
	}
}

// Optional middleware that doesn't require authentication but sets user context if token is present
func Optional() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		token := utils.ExtractTokenFromHeader(authHeader)
		if token == "" {
			return c.Next()
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		userCtx, err := utils.ValidateTokenStringToUUID(token, jwtSecret)
		if err != nil {
			return c.Next()
		}

		log.Printf("✅ Optional auth: Token validated for user: %s (%s)", userCtx.Email, userCtx.ID)

		// Set user context in fiber locals
		c.Locals("user", userCtx)
		c.Locals("userID", userCtx.ID)
		return c.Next()
	}
}

// WebSocketProtected middleware validates JWT tokens from query parameter or header
// This is specifically for WebSocket connections which can't set custom headers
func WebSocketProtected() fiber.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return func(c *fiber.Ctx) error {
		var token string

		// Try to get token from query parameter first (for WebSocket)
		token = c.Query("token")

		// If not in query, try Authorization header
		if token == "" {
			authHeader := c.Get("Authorization")
			if authHeader != "" {
				token = utils.ExtractTokenFromHeader(authHeader)
			}
		}

		// If still no token, reject
		if token == "" {
			log.Printf("❌ WebSocket auth failed: No token provided")
			return utils.UnauthorizedResponse(c, "Missing authentication token")
		}

		// Validate token and get user context
		userCtx, err := utils.ValidateTokenStringToUUID(token, jwtSecret)
		if err != nil {
			log.Printf("❌ WebSocket token validation failed: %v", err)
			switch err {
			case utils.ErrExpiredToken:
				return utils.UnauthorizedResponse(c, "Token has expired")
			case utils.ErrInvalidToken:
				return utils.UnauthorizedResponse(c, "Invalid token")
			case utils.ErrMissingToken:
				return utils.UnauthorizedResponse(c, "Missing token")
			default:
				return utils.UnauthorizedResponse(c, "Token validation failed")
			}
		}

		log.Printf("✅ WebSocket: Token validated from query param for user: %s (%s)", userCtx.Email, userCtx.ID)

		// Set user context in fiber locals
		c.Locals("user", userCtx)
		c.Locals("userID", userCtx.ID)

		return c.Next()
	}
}
