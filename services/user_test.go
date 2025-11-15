package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vkhangstack/go-zalo-bot/auth"
	"github.com/vkhangstack/go-zalo-bot/types"
)

func setupTestUserService(t *testing.T, botToken string) (*UserService, *types.Config) {
	config := &types.Config{
		BotToken:    botToken,
		BaseURL:     "https://bot-api.zapps.me",
		Timeout:     30 * time.Second,
		Retries:     3,
		Environment: types.Production,
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
	}
	
	if err := config.Validate(); err != nil {
		t.Fatalf("config validation failed: %v", err)
	}
	
	authService, err := auth.NewAuthService(config)
	if err != nil {
		t.Fatalf("failed to create auth service: %v", err)
	}
	
	service := NewUserService(authService, config.HTTPClient, config)
	return service, config
}

func TestUserService_GetUserProfile_Success(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL contains bot token
		expectedPath := "/bot" + botToken + "/getUserProfile"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %v, want %v", r.URL.Path, expectedPath)
		}
		
		// Verify method
		if r.Method != http.MethodGet {
			t.Errorf("Request method = %v, want %v", r.Method, http.MethodGet)
		}
		
		// Verify query parameters
		userID := r.URL.Query().Get("user_id")
		if userID != "user123" {
			t.Errorf("user_id = %v, want user123", userID)
		}
		
		// Return success response
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"id": "user123", "name": "Test User", "avatar": "https://example.com/avatar.jpg"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestUserService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	service.BaseService.authService = authService
	
	ctx := context.Background()
	userProfile, err := service.GetUserProfile(ctx, "user123")
	
	if err != nil {
		t.Errorf("GetUserProfile() error = %v", err)
	}
	
	if userProfile == nil {
		t.Fatal("GetUserProfile() returned nil user profile")
	}
	
	if userProfile.ID != "user123" {
		t.Errorf("UserProfile.ID = %v, want user123", userProfile.ID)
	}
	
	if userProfile.Name != "Test User" {
		t.Errorf("UserProfile.Name = %v, want Test User", userProfile.Name)
	}
	
	if userProfile.Avatar != "https://example.com/avatar.jpg" {
		t.Errorf("UserProfile.Avatar = %v, want https://example.com/avatar.jpg", userProfile.Avatar)
	}
}

func TestUserService_GetUserProfile_ValidationError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	service, _ := setupTestUserService(t, botToken)
	
	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "empty user ID",
			userID:  "",
			wantErr: true,
		},
		{
			name:    "user ID too long",
			userID:  string(make([]byte, 101)),
			wantErr: true,
		},
	}
	
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetUserProfile(ctx, tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_GetUserProfile_RateLimitHandling(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	attemptCount := 0
	
	// Create test server that returns rate limit error first, then success
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		
		if attemptCount == 1 {
			// First attempt: return rate limit error
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			resp := APIResponse{
				OK:          false,
				ErrorCode:   429,
				Description: "rate limit exceeded",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}
		
		// Second attempt: return success
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"id": "user123", "name": "Test User", "avatar": "https://example.com/avatar.jpg"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestUserService(t, botToken)
	config.BaseURL = server.URL
	config.RetryConfig = &types.RetryConfig{
		MaxRetries:    2,
		InitialDelay:  10 * time.Millisecond,
		MaxDelay:      100 * time.Millisecond,
		BackoffFactor: 2.0,
		RetryableErrors: []types.ErrorType{
			types.ErrorTypeRateLimit,
		},
	}
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	service.BaseService = NewBaseService(authService, config.HTTPClient, config)
	
	ctx := context.Background()
	userProfile, err := service.GetUserProfile(ctx, "user123")
	
	if err != nil {
		t.Errorf("GetUserProfile() error = %v", err)
	}
	
	if userProfile == nil {
		t.Fatal("GetUserProfile() returned nil user profile")
	}
	
	if attemptCount != 2 {
		t.Errorf("Expected 2 attempts, got %d", attemptCount)
	}
}

func TestUserService_GetUserProfile_APIError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	// Create test server that returns API error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse{
			OK:          false,
			ErrorCode:   404,
			Description: "user not found",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestUserService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	service.BaseService.authService = authService
	
	ctx := context.Background()
	userProfile, err := service.GetUserProfile(ctx, "user123")
	
	if err == nil {
		t.Error("GetUserProfile() expected error, got nil")
	}
	
	if userProfile != nil {
		t.Error("GetUserProfile() expected nil user profile on error")
	}
	
	// Verify error type
	zaloBotErr, ok := err.(*types.ZaloBotError)
	if !ok {
		t.Errorf("Expected ZaloBotError, got %T", err)
	} else {
		if zaloBotErr.Type != types.ErrorTypeAPI {
			t.Errorf("Expected ErrorTypeAPI, got %v", zaloBotErr.Type)
		}
	}
}

func TestValidateUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "valid user ID",
			userID:  "user123",
			wantErr: false,
		},
		{
			name:    "empty user ID",
			userID:  "",
			wantErr: true,
		},
		{
			name:    "user ID too long",
			userID:  string(make([]byte, 101)),
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateUserID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
