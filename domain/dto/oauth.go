package dto

// GoogleOAuthRequest - Request for Google OAuth callback
type GoogleOAuthRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state"`
}

// OAuthUserInfo - User information from OAuth provider
type OAuthUserInfo struct {
	Provider    string `json:"provider"`
	OAuthID     string `json:"oauthId"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	GivenName   string `json:"givenName"`
	FamilyName  string `json:"familyName"`
	Verified    bool   `json:"verified"`
}

// OAuthLoginResponse - Response after successful OAuth login
type OAuthLoginResponse struct {
	Token        string       `json:"token"`
	User         UserResponse `json:"user"`
	IsNewUser    bool         `json:"isNewUser"`
	NeedsProfile bool         `json:"needsProfile"` // True if user needs to complete profile (e.g., choose username)
}

// OAuthURLResponse - Response containing OAuth authorization URL
type OAuthURLResponse struct {
	URL string `json:"url"`
}

// ExchangeCodeRequest - Request to exchange authorization code for token
type ExchangeCodeRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state"`
}

// ExchangeCodeResponse - Response after exchanging code for token
type ExchangeCodeResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	IsNewUser bool         `json:"isNewUser"`
}
