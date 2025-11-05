package serviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
	"gofiber-template/pkg/utils"
)

type OAuthServiceImpl struct {
	userRepo     repositories.UserRepository
	config       *config.Config
	googleConfig *oauth2.Config
}

func NewOAuthService(userRepo repositories.UserRepository, cfg *config.Config) services.OAuthService {
	googleConfig := &oauth2.Config{
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		RedirectURL:  cfg.OAuth.Google.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &OAuthServiceImpl{
		userRepo:     userRepo,
		config:       cfg,
		googleConfig: googleConfig,
	}
}

func (s *OAuthServiceImpl) GetGoogleAuthURL(state string) string {
	return s.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *OAuthServiceImpl) HandleGoogleCallback(ctx context.Context, code string) (*dto.OAuthLoginResponse, error) {
	// Debug logging
	fmt.Printf("\n=== OAuth Service Debug ===\n")
	fmt.Printf("Exchanging code: %s\n", code)
	fmt.Printf("Code length: %d\n", len(code))
	fmt.Printf("Client ID: %s\n", s.googleConfig.ClientID)
	fmt.Printf("Redirect URL: %s\n", s.googleConfig.RedirectURL)
	fmt.Printf("===========================\n\n")

	// Exchange code for token
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		fmt.Printf("ERROR: Failed to exchange code: %v\n", err)
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from Google
	userInfo, err := s.GetUserInfoFromGoogle(ctx, token.AccessToken)
	if err != nil {
		return nil, err
	}

	// Check if user already exists by OAuth ID
	existingUser, err := s.userRepo.GetByOAuth(ctx, "google", userInfo.OAuthID)
	if err == nil && existingUser != nil {
		// User exists - login
		jwtToken, err := utils.GenerateToken(existingUser.ID, s.config.JWT.Secret)
		if err != nil {
			return nil, fmt.Errorf("failed to generate JWT: %w", err)
		}

		return &dto.OAuthLoginResponse{
			Token:        jwtToken,
			User:         *dto.UserToUserResponse(existingUser),
			IsNewUser:    false,
			NeedsProfile: false,
		}, nil
	}

	// Check if email already exists (user registered with email/password)
	existingEmailUser, err := s.userRepo.GetByEmail(ctx, userInfo.Email)
	if err == nil && existingEmailUser != nil {
		// Email exists but not linked to Google
		// Link Google account to existing user
		existingEmailUser.OAuthProvider = "google"
		existingEmailUser.OAuthID = userInfo.OAuthID
		existingEmailUser.IsOAuthUser = true
		existingEmailUser.UpdatedAt = time.Now()

		if err := s.userRepo.Update(ctx, existingEmailUser.ID, existingEmailUser); err != nil {
			return nil, fmt.Errorf("failed to link Google account: %w", err)
		}

		jwtToken, err := utils.GenerateToken(existingEmailUser.ID, s.config.JWT.Secret)
		if err != nil {
			return nil, fmt.Errorf("failed to generate JWT: %w", err)
		}

		return &dto.OAuthLoginResponse{
			Token:        jwtToken,
			User:         *dto.UserToUserResponse(existingEmailUser),
			IsNewUser:    false,
			NeedsProfile: false,
		}, nil
	}

	// Create new user
	username, err := s.generateUniqueUsername(ctx, userInfo.Email, userInfo.Name)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		ID:            uuid.New(),
		Email:         userInfo.Email,
		Username:      username,
		DisplayName:   userInfo.Name,
		Avatar:        userInfo.Picture,
		OAuthProvider: "google",
		OAuthID:       userInfo.OAuthID,
		IsOAuthUser:   true,
		Role:          "user",
		IsActive:      true,
		Karma:         0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	jwtToken, err := utils.GenerateToken(newUser.ID, s.config.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT: %w", err)
	}

	return &dto.OAuthLoginResponse{
		Token:        jwtToken,
		User:         *dto.UserToUserResponse(newUser),
		IsNewUser:    true,
		NeedsProfile: false, // Google provides all necessary info
	}, nil
}

func (s *OAuthServiceImpl) GetUserInfoFromGoogle(ctx context.Context, accessToken string) (*dto.OAuthUserInfo, error) {
	// Get user info from Google
	resp, err := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)).Get("https://www.googleapis.com/oauth2/v2/userinfo")

	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var googleUser struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}

	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	if googleUser.Email == "" {
		return nil, errors.New("email not provided by Google")
	}

	return &dto.OAuthUserInfo{
		Provider:   "google",
		OAuthID:    googleUser.ID,
		Email:      googleUser.Email,
		Name:       googleUser.Name,
		Picture:    googleUser.Picture,
		GivenName:  googleUser.GivenName,
		FamilyName: googleUser.FamilyName,
		Verified:   googleUser.VerifiedEmail,
	}, nil
}

// generateUniqueUsername generates a unique username from email or name
func (s *OAuthServiceImpl) generateUniqueUsername(ctx context.Context, email, name string) (string, error) {
	var baseUsername string

	// Try to use part of email before @
	if email != "" {
		emailParts := []rune(email)
		for i, r := range emailParts {
			if r == '@' {
				baseUsername = string(emailParts[:i])
				break
			}
		}
	}

	// Fallback to name if email didn't work
	if baseUsername == "" && name != "" {
		baseUsername = name
	}

	// Clean username (remove spaces, special chars)
	baseUsername = utils.CleanUsername(baseUsername)

	// Try base username first
	_, err := s.userRepo.GetByUsername(ctx, baseUsername)
	if err != nil {
		// Username available
		return baseUsername, nil
	}

	// Username taken, try with numbers
	for i := 1; i <= 100; i++ {
		username := fmt.Sprintf("%s%d", baseUsername, i)
		_, err := s.userRepo.GetByUsername(ctx, username)
		if err != nil {
			// Username available
			return username, nil
		}
	}

	// Generate random username as last resort
	return fmt.Sprintf("user_%s", uuid.New().String()[:8]), nil
}
