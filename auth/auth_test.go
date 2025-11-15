package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vkhangstack/go-zalo-bot/types"
)

func TestNewAuthService(t *testing.T) {
	config := &types.Config{
		BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		BaseURL:     "https://bot-api.zapps.me",
		Environment: types.Production,
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
	}
	
	authService, err := NewAuthService(config)
	if err != nil {
		t.Errorf("NewAuthService() error = %v", err)
	}
	
	if authService == nil {
		t.Error("Expected non-nil auth service")
	}
	
	if authService.GetToken() != config.BotToken {
		t.Errorf("Expected token %s, got %s", config.BotToken, authService.GetToken())
	}
	
	if authService.GetEnvironment() != config.Environment {
		t.Errorf("Expected environment %s, got %s", config.Environment, authService.GetEnvironment())
	}
}

func TestNewAuthService_InvalidConfig(t *testing.T) {
	config := &types.Config{
		BotToken: "", // Invalid empty bot token
	}
	
	_, err := NewAuthService(config)
	if err == nil {
		t.Error("Expected error for invalid config")
	}
}

func TestAuthService_ValidateCredentials(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		responseBody   string
		expectedError  bool
		expectedErrType types.ErrorType
	}{
		{
			name:          "valid credentials",
			statusCode:    http.StatusOK,
			responseBody:  `{"data": {"id": "123", "name": "Test Bot"}}`,
			expectedError: false,
		},
		{
			name:            "unauthorized",
			statusCode:      http.StatusUnauthorized,
			responseBody:    `{"error": {"code": 401, "message": "Unauthorized"}}`,
			expectedError:   true,
			expectedErrType: types.ErrorTypeAuth,
		},
		{
			name:            "forbidden",
			statusCode:      http.StatusForbidden,
			responseBody:    `{"error": {"code": 403, "message": "Forbidden"}}`,
			expectedError:   true,
			expectedErrType: types.ErrorTypeAuth,
		},
		{
			name:            "rate limited",
			statusCode:      http.StatusTooManyRequests,
			responseBody:    `{"error": {"code": 429, "message": "Too Many Requests"}}`,
			expectedError:   true,
			expectedErrType: types.ErrorTypeRateLimit,
		},
		{
			name:            "server error",
			statusCode:      http.StatusInternalServerError,
			responseBody:    `{"error": {"code": 500, "message": "Internal Server Error"}}`,
			expectedError:   true,
			expectedErrType: types.ErrorTypeAPI,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify bot token is in URL path, not in Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader != "" {
					t.Error("Expected no Authorization header for bot token authentication")
				}
				
				// Verify URL contains bot token
				if !contains(r.URL.Path, "bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11") {
					t.Errorf("Expected bot token in URL path, got: %s", r.URL.Path)
				}
				
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()
			
			config := &types.Config{
				BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
				BaseURL:     server.URL,
				Environment: types.Development,
				HTTPClient:  &http.Client{Timeout: 5 * time.Second},
			}
			
			authService, err := NewAuthService(config)
			if err != nil {
				t.Fatalf("NewAuthService() error = %v", err)
			}
			
			ctx := context.Background()
			err = authService.ValidateCredentials(ctx)
			
			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
					return
				}
				
				zaloBotErr, ok := err.(*types.ZaloBotError)
				if !ok {
					t.Errorf("Expected ZaloBotError, got %T", err)
					return
				}
				
				if zaloBotErr.Type != tt.expectedErrType {
					t.Errorf("Expected error type %s, got %s", tt.expectedErrType, zaloBotErr.Type)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateCredentials() error = %v", err)
				}
			}
		})
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestAuthService_SetToken(t *testing.T) {
	config := &types.Config{
		BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		BaseURL:     "https://bot-api.zapps.me",
		Environment: types.Production,
		HTTPClient:  &http.Client{},
	}
	
	authService, err := NewAuthService(config)
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	
	// Test setting valid token with proper format
	newToken := "789012:XYZ-GHI5678jklMn-abc90D3f4g567hi89"
	err = authService.SetToken(newToken)
	if err != nil {
		t.Errorf("SetToken() error = %v", err)
	}
	
	if authService.GetToken() != newToken {
		t.Errorf("Expected token %s, got %s", newToken, authService.GetToken())
	}
	
	// Test invalid token
	err = authService.SetToken("")
	if err == nil {
		t.Error("Expected error for empty token")
	}
}

func TestAuthService_IsAuthenticated(t *testing.T) {
	config := &types.Config{
		BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		BaseURL:     "https://bot-api.zapps.me",
		Environment: types.Production,
		HTTPClient:  &http.Client{},
	}
	
	authService, err := NewAuthService(config)
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	
	if !authService.IsAuthenticated() {
		t.Error("Expected service to be authenticated")
	}
	
	// Clear token
	authService.tokenManager.Clear()
	if authService.IsAuthenticated() {
		t.Error("Expected service not to be authenticated after clearing token")
	}
}

func TestAuthService_ValidateEnvironmentConfig(t *testing.T) {
	tests := []struct {
		name        string
		environment types.Environment
		apiEndpoint string
		token       string
		wantErr     bool
	}{
		{
			name:        "valid development config",
			environment: types.Development,
			apiEndpoint: "https://dev-bot-api.zapps.me",
			token:       "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr:     false,
		},
		{
			name:        "valid production config",
			environment: types.Production,
			apiEndpoint: "https://bot-api.zapps.me",
			token:       "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr:     false,
		},
		{
			name:        "development with empty endpoint",
			environment: types.Development,
			apiEndpoint: "",
			token:       "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr:     true,
		},
		{
			name:        "production with wrong endpoint",
			environment: types.Production,
			apiEndpoint: "https://wrong-endpoint.com",
			token:       "prod-token-123",
			wantErr:     true,
		},
		{
			name:        "production with empty token",
			environment: types.Production,
			apiEndpoint: "https://bot-api.zapps.me",
			token:       "",
			wantErr:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.Config{
				BotToken:    tt.token,
				BaseURL:     tt.apiEndpoint,
				Environment: tt.environment,
				HTTPClient:  &http.Client{},
			}
			
			// Skip validation during creation for invalid configs
			authService := &AuthService{
				tokenManager: NewTokenManager(tt.token),
				httpClient:   config.HTTPClient,
				apiEndpoint:  tt.apiEndpoint,
				environment:  tt.environment,
			}
			
			err := authService.ValidateEnvironmentConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEnvironmentConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthService_CreateAuthenticatedRequest(t *testing.T) {
	config := &types.Config{
		BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		BaseURL:     "https://bot-api.zapps.me",
		Environment: types.Development,
		HTTPClient:  &http.Client{},
	}
	
	authService, err := NewAuthService(config)
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	
	ctx := context.Background()
	req, err := authService.CreateAuthenticatedRequest(ctx, "GET", "sendMessage")
	if err != nil {
		t.Errorf("CreateAuthenticatedRequest() error = %v", err)
	}
	
	// Bot token authentication doesn't use Authorization headers
	// The token is embedded in the URL instead
	authHeader := req.Header.Get("Authorization")
	if authHeader != "" {
		t.Errorf("Expected no auth header for bot token authentication, got %s", authHeader)
	}
	
	// Verify URL contains bot token
	expectedURL := "https://bot-api.zapps.me/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/sendMessage"
	if req.URL.String() != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, req.URL.String())
	}
	
	// Check content type
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", contentType)
	}
	
	// Check User-Agent header
	userAgent := req.Header.Get("User-Agent")
	if userAgent != "Go-Zalo-Bot-SDK/1.0" {
		t.Errorf("Expected User-Agent Go-Zalo-Bot-SDK/1.0, got %s", userAgent)
	}
	
	// Check environment header for development
	envHeader := req.Header.Get("X-Environment")
	if envHeader != "development" {
		t.Errorf("Expected environment header development, got %s", envHeader)
	}
	
	// Test with unauthenticated service
	authService.tokenManager.Clear()
	_, err = authService.CreateAuthenticatedRequest(ctx, "GET", "sendMessage")
	if err == nil {
		t.Error("Expected error for unauthenticated request")
	}
}

func TestAuthService_GetAPIEndpoint(t *testing.T) {
	config := &types.Config{
		BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		BaseURL:     "https://bot-api.zapps.me",
		Environment: types.Production,
		HTTPClient:  &http.Client{},
	}
	
	authService, err := NewAuthService(config)
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	
	tests := []struct {
		name     string
		method   string
		expected string
	}{
		{
			name:     "sendMessage method",
			method:   "sendMessage",
			expected: "https://bot-api.zapps.me/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/sendMessage",
		},
		{
			name:     "getMe method",
			method:   "getMe",
			expected: "https://bot-api.zapps.me/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getMe",
		},
		{
			name:     "getUserProfile method",
			method:   "getUserProfile",
			expected: "https://bot-api.zapps.me/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getUserProfile",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := authService.GetAPIEndpoint(tt.method)
			if endpoint != tt.expected {
				t.Errorf("Expected endpoint %s, got %s", tt.expected, endpoint)
			}
		})
	}
}

func TestAuthService_GetAPIEndpoint_DifferentEnvironments(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		botToken    string
		method      string
		expected    string
	}{
		{
			name:     "production environment",
			baseURL:  "https://bot-api.zapps.me",
			botToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			method:   "sendMessage",
			expected: "https://bot-api.zapps.me/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/sendMessage",
		},
		{
			name:     "development environment",
			baseURL:  "https://dev-bot-api.zapps.me",
			botToken: "789012:XYZ-GHI5678jklMn-abc90D3f4g567hi89",
			method:   "getMe",
			expected: "https://dev-bot-api.zapps.me/bot789012:XYZ-GHI5678jklMn-abc90D3f4g567hi89/getMe",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.Config{
				BotToken:    tt.botToken,
				BaseURL:     tt.baseURL,
				Environment: types.Development,
				HTTPClient:  &http.Client{},
			}
			
			authService, err := NewAuthService(config)
			if err != nil {
				t.Fatalf("NewAuthService() error = %v", err)
			}
			
			endpoint := authService.GetAPIEndpoint(tt.method)
			if endpoint != tt.expected {
				t.Errorf("Expected endpoint %s, got %s", tt.expected, endpoint)
			}
		})
	}
}

func TestAuthService_HandleAuthError(t *testing.T) {
	config := &types.Config{
		BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		BaseURL:     "https://bot-api.zapps.me",
		Environment: types.Production,
		HTTPClient:  &http.Client{},
	}
	
	authService, err := NewAuthService(config)
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	
	// Test with auth error - should clear token
	authErr := types.NewAuthError("invalid token")
	handledErr := authService.HandleAuthError(authErr)
	
	if handledErr != authErr {
		t.Error("Expected same error to be returned")
	}
	
	if authService.GetToken() != "" {
		t.Error("Expected token to be cleared after auth error")
	}
	
	// Reset token for next test with valid bot token format
	authService.SetToken("789012:XYZ-GHI5678jklMn-abc90D3f4g567hi89")
	
	// Test with rate limit error - should not clear token
	rateLimitErr := types.NewRateLimitError("rate limited")
	handledErr = authService.HandleAuthError(rateLimitErr)
	
	if handledErr != rateLimitErr {
		t.Error("Expected same error to be returned")
	}
	
	if authService.GetToken() == "" {
		t.Error("Expected token not to be cleared after rate limit error")
	}
}

func TestAuthService_SetEnvironment(t *testing.T) {
	config := &types.Config{
		BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
		BaseURL:     "https://bot-api.zapps.me",
		Environment: types.Production,
		HTTPClient:  &http.Client{},
	}
	
	authService, err := NewAuthService(config)
	if err != nil {
		t.Fatalf("NewAuthService() error = %v", err)
	}
	
	// Test setting valid environment
	err = authService.SetEnvironment(types.Development)
	if err != nil {
		t.Errorf("SetEnvironment() error = %v", err)
	}
	
	if authService.GetEnvironment() != types.Development {
		t.Errorf("Expected environment %s, got %s", types.Development, authService.GetEnvironment())
	}
	
	// Test setting invalid environment
	err = authService.SetEnvironment(types.Environment("invalid"))
	if err == nil {
		t.Error("Expected error for invalid environment")
	}
}