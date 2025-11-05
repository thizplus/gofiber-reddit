package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/auth_code_store"
	"gofiber-template/pkg/config"
	"gofiber-template/pkg/utils"
)

type OAuthHandler struct {
	oauthService services.OAuthService
	config       *config.Config
}

func NewOAuthHandler(oauthService services.OAuthService, cfg *config.Config) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		config:       cfg,
	}
}

// GetGoogleAuthURL generates Google OAuth authorization URL
// @Summary Get Google OAuth URL
// @Description Get Google OAuth authorization URL to start OAuth flow
// @Tags OAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.OAuthURLResponse
// @Router /auth/google [get]
func (h *OAuthHandler) GetGoogleAuthURL(c *fiber.Ctx) error {
	// Generate random state for CSRF protection
	state := utils.GenerateRandomString(32)

	// Store state in session or cookie for validation in callback
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
		Path:     "/",
	})

	url := h.oauthService.GetGoogleAuthURL(state)

	return utils.SuccessResponse(c, "Google OAuth URL generated", dto.OAuthURLResponse{
		URL: url,
	})
}

// GoogleCallback handles Google OAuth callback
// @Summary Handle Google OAuth Callback
// @Description Process Google OAuth callback and login/register user
// @Tags OAuth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State parameter for CSRF protection"
// @Success 200 {object} dto.OAuthLoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /auth/google/callback [get]
func (h *OAuthHandler) GoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	// Debug logging
	fmt.Printf("\n=== Google OAuth Callback Debug ===\n")
	fmt.Printf("Received code: %s\n", code)
	fmt.Printf("Code length: %d\n", len(code))
	fmt.Printf("Received state: %s\n", state)
	fmt.Printf("Redirect URL from config: %s\n", h.config.OAuth.Google.RedirectURL)
	fmt.Printf("===================================\n\n")

	if code == "" {
		// Redirect to frontend with error if redirect URL is provided
		redirectURL := c.Query("redirect_url")
		if redirectURL != "" {
			return c.Redirect(redirectURL + "?error=missing_code")
		}
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Authorization code is required", nil)
	}

	// Validate state parameter for CSRF protection
	storedState := c.Cookies("oauth_state")

	// Debug logging
	c.Locals("debug_received_state", state)
	c.Locals("debug_stored_state", storedState)

	// Skip state validation if cookie is not present (direct Google redirect)
	// In production, this should be strict
	if storedState == "" {
		// No cookie found - user likely came directly from Google without calling /auth/google first
		// For now, we'll allow this in development
		c.Locals("debug_state_validation", "skipped_no_cookie")
	} else if storedState != state {
		// Cookie exists but doesn't match
		redirectURL := c.Query("redirect_url")
		if redirectURL != "" {
			return c.Redirect(redirectURL + "?error=invalid_state")
		}
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid state parameter", nil)
	} else {
		// Clear state cookie only if validation passed
		c.ClearCookie("oauth_state")
		c.Locals("debug_state_validation", "passed")
	}

	// Handle OAuth callback
	response, err := h.oauthService.HandleGoogleCallback(c.Context(), code)
	if err != nil {
		// Redirect to frontend with error if redirect URL is provided
		redirectURL := c.Query("redirect_url")
		if redirectURL != "" {
			return c.Redirect(redirectURL + "?error=oauth_failed")
		}
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "OAuth authentication failed", err)
	}

	// Generate authorization code
	store := auth_code_store.GetInstance()
	authCode, err := store.GenerateCode(response.Token, response.User, response.IsNewUser, state)
	if err != nil {
		redirectURL := c.Query("redirect_url")
		if redirectURL != "" {
			return c.Redirect(redirectURL + "?error=code_generation_failed")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate authorization code", err)
	}

	// Get redirect URL from query or use default frontend URL
	redirectURL := c.Query("redirect_url")
	if redirectURL == "" {
		// Use default frontend URL from config
		redirectURL = h.config.App.FrontendURL + "/auth/callback"
	}

	// Redirect to frontend with authorization code and state
	return c.Redirect(redirectURL + "?code=" + authCode + "&state=" + state)
}

// ExchangeCodeForToken exchanges authorization code for JWT token
// @Summary Exchange authorization code for token
// @Description Exchange authorization code for JWT token
// @Tags OAuth
// @Accept json
// @Produce json
// @Param request body dto.ExchangeCodeRequest true "Exchange code request"
// @Success 200 {object} dto.ExchangeCodeResponse
// @Failure 400 {object} map[string]interface{}
// @Router /auth/exchange [post]
func (h *OAuthHandler) ExchangeCodeForToken(c *fiber.Ctx) error {
	var req dto.ExchangeCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Exchange code for token
	store := auth_code_store.GetInstance()
	data, ok := store.ExchangeCode(req.Code, req.State)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid or expired authorization code", nil)
	}

	// Return token and user info
	return utils.SuccessResponse(c, "Token exchanged successfully", dto.ExchangeCodeResponse{
		Token:     data.Token,
		IsNewUser: data.IsNewUser,
		User:      data.User,
	})
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
