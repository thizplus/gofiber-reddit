package services

import (
	"context"
	"gofiber-template/domain/dto"
)

type OAuthService interface {
	// GetGoogleAuthURL generates Google OAuth authorization URL
	GetGoogleAuthURL(state string) string

	// HandleGoogleCallback processes Google OAuth callback
	HandleGoogleCallback(ctx context.Context, code string) (*dto.OAuthLoginResponse, error)

	// GetUserInfoFromGoogle retrieves user info from Google OAuth token
	GetUserInfoFromGoogle(ctx context.Context, accessToken string) (*dto.OAuthUserInfo, error)
}
