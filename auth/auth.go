package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vkhangstack/go-zalo-bot/types"
)

// AuthService handles authentication operations for Zalo Bot API
type AuthService struct {
	tokenManager *TokenManager
	httpClient   *http.Client
	apiEndpoint  string
	environment  types.Environment
}

// NewAuthService creates a new authentication service
func NewAuthService(config *types.Config) (*AuthService, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	tokenManager := NewTokenManager(config.BotToken)

	return &AuthService{
		tokenManager: tokenManager,
		httpClient:   config.HTTPClient,
		apiEndpoint:  config.BaseURL,
		environment:  config.Environment,
	}, nil
}

// GetTokenManager returns the token manager
func (as *AuthService) GetTokenManager() *TokenManager {
	return as.tokenManager
}

// ValidateCredentials validates the provided credentials against Zalo API
func (as *AuthService) ValidateCredentials(ctx context.Context) error {
	if !as.tokenManager.IsValid() {
		return types.NewAuthError("invalid bot token")
	}

	// Make a test API call to validate credentials
	// Using the "getMe" equivalent endpoint for Zalo Bot API with bot token in URL
	url := as.GetAPIEndpoint("getMe")

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return types.NewNetworkError(fmt.Sprintf("failed to create request: %v", err))
	}

	// Bot token authentication doesn't use Authorization headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Go-Zalo-Bot-SDK/1.0")

	// Add environment-specific headers if needed
	if as.environment == types.Development {
		req.Header.Set("X-Environment", "development")
	}

	resp, err := as.httpClient.Do(req)
	if err != nil {
		return types.NewNetworkError(fmt.Sprintf("failed to make request: %v", err))
	}
	defer resp.Body.Close()

	// Handle different response status codes
	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusUnauthorized:
		return types.NewAuthError("invalid credentials or token")
	case http.StatusForbidden:
		return types.NewAuthError("access forbidden - check token permissions")
	case http.StatusTooManyRequests:
		return types.NewRateLimitError("rate limit exceeded")
	default:
		// Try to parse error response
		var apiResp struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&apiResp); err == nil && apiResp.Error.Message != "" {
			return types.NewAPIError(apiResp.Error.Code, apiResp.Error.Message, "credential validation failed")
		}

		return types.NewAPIError(resp.StatusCode, "credential validation failed", fmt.Sprintf("HTTP %d", resp.StatusCode))
	}
}

// SetToken updates the authentication token
func (as *AuthService) SetToken(token string) error {
	return as.tokenManager.SetToken(token)
}

// GetToken returns the current authentication token
func (as *AuthService) GetToken() string {
	return as.tokenManager.GetToken()
}

// IsAuthenticated checks if the service is properly authenticated
func (as *AuthService) IsAuthenticated() bool {
	return as.tokenManager.IsValid()
}

// GetAPIEndpoint constructs the full API endpoint URL with embedded bot token
// Pattern: https://bot-api.zapps.me/bot${BOT_TOKEN}/method
func (as *AuthService) GetAPIEndpoint(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", as.apiEndpoint, as.tokenManager.GetBotToken(), method)
}

// ValidateEnvironmentConfig validates environment-specific configuration
func (as *AuthService) ValidateEnvironmentConfig() error {
	switch as.environment {
	case types.Development:
		// In development, we might have more lenient validation
		if as.apiEndpoint == "" {
			return types.NewValidationError("API endpoint is required for development environment")
		}

		// Check if using development API endpoint
		if as.apiEndpoint == "https://bot-api.zapps.me" {
			// This is the production endpoint, warn but don't fail
			// In a real implementation, you might want to log a warning here
		}

	case types.Production:
		// In production, we have stricter validation
		if as.apiEndpoint == "" {
			return types.NewValidationError("API endpoint is required for production environment")
		}

		// Ensure we're using the correct production endpoint
		if as.apiEndpoint != "https://bot-api.zapps.me" {
			return types.NewValidationError("invalid API endpoint for production environment")
		}

		// Additional production checks
		if !as.tokenManager.IsValid() {
			return types.NewAuthError("valid token is required for production environment")
		}

	default:
		return types.NewValidationError("unsupported environment")
	}

	return nil
}

// CreateAuthenticatedRequest creates an HTTP request with proper authentication
// Note: For bot token authentication, the token is embedded in the URL, not headers
func (as *AuthService) CreateAuthenticatedRequest(ctx context.Context, httpMethod, apiMethod string) (*http.Request, error) {
	if !as.IsAuthenticated() {
		return nil, types.NewAuthError("not authenticated")
	}

	// Construct URL with bot token embedded using GetAPIEndpoint
	url := as.GetAPIEndpoint(apiMethod)

	req, err := http.NewRequestWithContext(ctx, httpMethod, url, nil)
	if err != nil {
		return nil, types.NewNetworkError(fmt.Sprintf("failed to create request: %v", err))
	}

	// Bot token authentication doesn't use Authorization headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Go-Zalo-Bot-SDK/1.0")

	// Add environment-specific headers
	if as.environment == types.Development {
		req.Header.Set("X-Environment", "development")
	}

	return req, nil
}

// HandleAuthError processes authentication-related errors and determines if retry is needed
func (as *AuthService) HandleAuthError(err error) error {
	if err == nil {
		return nil
	}

	zaloBotErr, ok := err.(*types.ZaloBotError)
	if !ok {
		return err
	}

	switch zaloBotErr.Type {
	case types.ErrorTypeAuth:
		// Clear the token if we get an auth error
		as.tokenManager.Clear()
		return zaloBotErr
	case types.ErrorTypeRateLimit:
		// Don't clear token for rate limit errors
		return zaloBotErr
	default:
		return zaloBotErr
	}
}

// GetEnvironment returns the current environment
func (as *AuthService) GetEnvironment() types.Environment {
	return as.environment
}

// SetEnvironment updates the environment configuration
func (as *AuthService) SetEnvironment(env types.Environment) error {
	if !env.IsValid() {
		return types.NewValidationError("invalid environment")
	}

	as.environment = env
	return as.ValidateEnvironmentConfig()
}
