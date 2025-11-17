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

func setupTestService(t *testing.T, botToken string) (*BaseService, *types.Config) {
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

	service := NewBaseService(authService, config.HTTPClient, config)
	return service, config
}

func TestBaseService_DoRequest_Success(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL contains bot token
		expectedPath := "/bot" + botToken + "/testMethod"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %v, want %v", r.URL.Path, expectedPath)
		}

		// Verify headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %v, want application/json", r.Header.Get("Content-Type"))
		}

		// Return success response
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "123"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	service, config := setupTestService(t, botToken)
	config.BaseURL = server.URL

	// Update auth service with new base URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService

	req := &APIRequest{
		Method:    "POST",
		APIMethod: "testMethod",
		Body:      map[string]string{"test": "data"},
	}

	ctx := context.Background()
	resp, err := service.DoRequest(ctx, req)

	if err != nil {
		t.Errorf("DoRequest() error = %v", err)
	}

	if resp == nil {
		t.Fatal("DoRequest() returned nil response")
	}

	if !resp.OK {
		t.Error("Response OK = false, want true")
	}
}

func TestBaseService_DoRequest_RateLimitError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"

	// Create test server that returns rate limit error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)

		resp := APIResponse{
			OK:          false,
			ErrorCode:   429,
			Description: "rate limit exceeded",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	service, config := setupTestService(t, botToken)
	config.BaseURL = server.URL
	config.Retries = 0 // Disable retries for this test

	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	service.config = config

	req := &APIRequest{
		Method:    "POST",
		APIMethod: "testMethod",
	}

	ctx := context.Background()
	_, err := service.DoRequest(ctx, req)

	if err == nil {
		t.Fatal("DoRequest() expected error but got none")
	}

	zaloBotErr, ok := err.(*types.ZaloBotError)
	if !ok {
		t.Fatalf("Error type = %T, want *types.ZaloBotError", err)
	}

	if zaloBotErr.Type != types.ErrorTypeRateLimit {
		t.Errorf("Error type = %v, want %v", zaloBotErr.Type, types.ErrorTypeRateLimit)
	}
}

func TestBaseService_DoRequest_AuthError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"

	// Create test server that returns auth error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		resp := APIResponse{
			OK:          false,
			ErrorCode:   401,
			Description: "invalid token",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	service, config := setupTestService(t, botToken)
	config.BaseURL = server.URL
	config.Retries = 0

	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	service.config = config

	req := &APIRequest{
		Method:    "POST",
		APIMethod: "testMethod",
	}

	ctx := context.Background()
	_, err := service.DoRequest(ctx, req)

	if err == nil {
		t.Fatal("DoRequest() expected error but got none")
	}

	zaloBotErr, ok := err.(*types.ZaloBotError)
	if !ok {
		t.Fatalf("Error type = %T, want *types.ZaloBotError", err)
	}

	if zaloBotErr.Type != types.ErrorTypeAuth {
		t.Errorf("Error type = %v, want %v", zaloBotErr.Type, types.ErrorTypeAuth)
	}
}

func TestBaseService_DoRequest_RetryMechanism(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"

	attemptCount := 0
	maxAttempts := 3

	// Create test server that fails first attempts then succeeds
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++

		if attemptCount < maxAttempts {
			// Return server error for first attempts
			w.WriteHeader(http.StatusInternalServerError)
			resp := APIResponse{
				OK:          false,
				ErrorCode:   500,
				Description: "server error",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Success on final attempt
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"success": true}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	service, config := setupTestService(t, botToken)
	config.BaseURL = server.URL
	config.Retries = 3
	config.RetryConfig = &types.RetryConfig{
		MaxRetries:    3,
		InitialDelay:  10 * time.Millisecond,
		MaxDelay:      100 * time.Millisecond,
		BackoffFactor: 2.0,
		RetryableErrors: []types.ErrorType{
			types.ErrorTypeNetwork,
			types.ErrorTypeAPI,
		},
	}

	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	service.config = config

	req := &APIRequest{
		Method:    "POST",
		APIMethod: "testMethod",
	}

	ctx := context.Background()
	resp, err := service.DoRequest(ctx, req)

	if err != nil {
		t.Errorf("DoRequest() error = %v", err)
	}

	if resp == nil {
		t.Fatal("DoRequest() returned nil response")
	}

	if !resp.OK {
		t.Error("Response OK = false, want true")
	}

	if attemptCount != maxAttempts {
		t.Errorf("Attempt count = %d, want %d", attemptCount, maxAttempts)
	}
}

func TestBaseService_DoRequest_Timeout(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"

	// Create test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"success": true}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	service, config := setupTestService(t, botToken)
	config.BaseURL = server.URL
	config.Timeout = 50 * time.Millisecond
	config.HTTPClient = &http.Client{Timeout: 50 * time.Millisecond}
	config.Retries = 0

	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	service.config = config
	service.client = config.HTTPClient

	req := &APIRequest{
		Method:    "POST",
		APIMethod: "testMethod",
	}

	ctx := context.Background()
	_, err := service.DoRequest(ctx, req)

	if err == nil {
		t.Fatal("DoRequest() expected timeout error but got none")
	}

	zaloBotErr, ok := err.(*types.ZaloBotError)
	if !ok {
		t.Fatalf("Error type = %T, want *types.ZaloBotError", err)
	}

	if zaloBotErr.Type != types.ErrorTypeNetwork {
		t.Errorf("Error type = %v, want %v", zaloBotErr.Type, types.ErrorTypeNetwork)
	}
}

func TestBaseService_CalculateBackoffDelay(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	service, _ := setupTestService(t, botToken)

	retryConfig := &types.RetryConfig{
		InitialDelay:  1 * time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
	}

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{
			name:    "first retry",
			attempt: 1,
			want:    2 * time.Second, // 1 * 2^1
		},
		{
			name:    "second retry",
			attempt: 2,
			want:    4 * time.Second, // 1 * 2^2
		},
		{
			name:    "third retry",
			attempt: 3,
			want:    8 * time.Second, // 1 * 2^3
		},
		{
			name:    "exceeds max delay",
			attempt: 10,
			want:    30 * time.Second, // capped at MaxDelay
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.calculateBackoffDelay(tt.attempt, retryConfig)
			if got != tt.want {
				t.Errorf("calculateBackoffDelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseService_DoRequest_WithQueryParams(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		if r.URL.Query().Get("offset") != "10" {
			t.Errorf("Query param offset = %v, want 10", r.URL.Query().Get("offset"))
		}
		if r.URL.Query().Get("limit") != "100" {
			t.Errorf("Query param limit = %v, want 100", r.URL.Query().Get("limit"))
		}

		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"success": true}`),
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	service, config := setupTestService(t, botToken)
	config.BaseURL = server.URL

	authService, _ := auth.NewAuthService(config)
	service.authService = authService

	req := &APIRequest{
		Method:    "GET",
		APIMethod: "getUpdates",
		QueryParams: map[string]string{
			"offset": "10",
			"limit":  "100",
		},
	}

	ctx := context.Background()
	resp, err := service.DoRequest(ctx, req)

	if err != nil {
		t.Errorf("DoRequest() error = %v", err)
	}

	if resp == nil || !resp.OK {
		t.Error("DoRequest() failed")
	}
}
