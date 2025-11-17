package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/vkhangstack/go-zalo-bot/auth"
	"github.com/vkhangstack/go-zalo-bot/types"
)

// BaseService provides common functionality for all services
type BaseService struct {
	authService *auth.AuthService
	client      *http.Client
	config      *types.Config
}

// NewBaseService creates a new base service
func NewBaseService(authService *auth.AuthService, client *http.Client, config *types.Config) *BaseService {
	return &BaseService{
		authService: authService,
		client:      client,
		config:      config,
	}
}

// APIRequest represents a request to the Zalo Bot API
type APIRequest struct {
	Method      string
	APIMethod   string
	Body        interface{}
	QueryParams map[string]string
}

// APIResponse represents a response from the Zalo Bot API
type APIResponse struct {
	OK          bool            `json:"ok"`
	Result      json.RawMessage `json:"result,omitempty"`
	ErrorCode   int             `json:"error_code,omitempty"`
	Description string          `json:"description,omitempty"`
}

// IsError returns true if the response indicates an error
func (r *APIResponse) IsError() bool {
	return !r.OK
}

// GetError extracts error information from the response
func (r *APIResponse) GetError() error {
	if !r.IsError() {
		return nil
	}

	// Check if it's a rate limit error
	if r.ErrorCode == 429 {
		return types.NewRateLimitError(r.Description)
	}

	// Check if it's an auth error
	if r.ErrorCode == 401 || r.ErrorCode == 403 {
		return types.NewAuthError(r.Description)
	}

	// Generic API error
	return types.NewAPIError(r.ErrorCode, r.Description, "API request failed")
}

// DoRequest executes an HTTP request with retry logic, connection pooling, and timeout handling
// URL pattern: https://bot-api.zapps.me/bot${BOT_TOKEN}/method
func (s *BaseService) DoRequest(ctx context.Context, apiReq *APIRequest) (*APIResponse, error) {
	var lastErr error
	retryConfig := s.config.RetryConfig
	if retryConfig == nil {
		retryConfig = types.DefaultRetryConfig()
	}

	// Retry loop with exponential backoff
	for attempt := 0; attempt <= retryConfig.MaxRetries; attempt++ {
		// Add delay for retry attempts
		if attempt > 0 {
			delay := s.calculateBackoffDelay(attempt, retryConfig)
			select {
			case <-ctx.Done():
				return nil, types.NewNetworkError(fmt.Sprintf("request cancelled: %v", ctx.Err()))
			case <-time.After(delay):
				// Continue with retry
			}
		}

		// Execute the request
		resp, err := s.executeRequest(ctx, apiReq)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Check if error is retryable
		if !retryConfig.ShouldRetry(err) {
			return nil, err
		}

		// Check if we've exhausted retries
		if attempt >= retryConfig.MaxRetries {
			break
		}
	}

	return nil, lastErr
}

// executeRequest performs a single HTTP request
func (s *BaseService) executeRequest(ctx context.Context, apiReq *APIRequest) (*APIResponse, error) {
	// Construct URL with bot token embedded
	url := s.authService.GetAPIEndpoint(apiReq.APIMethod)

	// Add query parameters if any
	if len(apiReq.QueryParams) > 0 {
		url += "?"
		first := true
		for key, value := range apiReq.QueryParams {
			if !first {
				url += "&"
			}
			url += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	// Prepare request body
	var bodyReader io.Reader
	if apiReq.Body != nil {
		bodyBytes, err := json.Marshal(apiReq.Body)
		if err != nil {
			return nil, types.NewValidationError(fmt.Sprintf("failed to marshal request body: %v", err))
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, apiReq.Method, url, bodyReader)
	if err != nil {
		return nil, types.NewNetworkError(fmt.Sprintf("failed to create request: %v", err))
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent())

	// Add environment-specific headers
	if s.config.Environment == types.Development {
		req.Header.Set("X-Environment", "development")
	}

	// Execute request with connection pooling (handled by http.Client)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, types.NewNetworkError(fmt.Sprintf("request failed: %v", err))
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.NewNetworkError(fmt.Sprintf("failed to read response body: %v", err))
	}

	// Parse rate limit information from headers
	rateLimitInfo := types.ParseRateLimitHeaders(resp.Header)

	// Handle HTTP status codes
	if resp.StatusCode == http.StatusTooManyRequests {
		errMsg := "rate limit exceeded"
		if rateLimitInfo.RetryAfter > 0 {
			errMsg = fmt.Sprintf("rate limit exceeded, retry after %s", rateLimitInfo.RetryAfter)
		}
		return nil, types.NewRateLimitError(errMsg)
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, types.NewAuthError(fmt.Sprintf("authentication failed: HTTP %d", resp.StatusCode))
	}

	// Parse API response
	var apiResp APIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		// If we can't parse the response, check if it's a server error
		if resp.StatusCode >= 500 {
			return nil, types.NewAPIError(resp.StatusCode, "server error", string(bodyBytes))
		}
		return nil, types.NewAPIError(resp.StatusCode, "failed to parse response", err.Error())
	}

	// Check if API returned an error
	if !apiResp.OK {
		// Use GetError to properly categorize the error
		if err := apiResp.GetError(); err != nil {
			return nil, err
		}
		return nil, types.NewAPIError(apiResp.ErrorCode, apiResp.Description, fmt.Sprintf("API method: %s", apiReq.APIMethod))
	}

	return &apiResp, nil
}

// calculateBackoffDelay calculates the delay for exponential backoff
func (s *BaseService) calculateBackoffDelay(attempt int, config *types.RetryConfig) time.Duration {
	if attempt <= 0 {
		return config.InitialDelay
	}

	// Calculate exponential backoff: initialDelay * (backoffFactor ^ attempt)
	delay := time.Duration(float64(config.InitialDelay) * math.Pow(config.BackoffFactor, float64(attempt)))

	// Cap at max delay
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	return delay
}

// GetAuthService returns the authentication service
func (s *BaseService) GetAuthService() *auth.AuthService {
	return s.authService
}

// GetHTTPClient returns the HTTP client
func (s *BaseService) GetHTTPClient() *http.Client {
	return s.client
}

// GetConfig returns the configuration
func (s *BaseService) GetConfig() *types.Config {
	return s.config
}

// GetFieldSignature returns the header field name for webhook signature
func (s *BaseService) GetFieldSignature() string {
	return "X-Zalo-Signature"
}

// parseResult parses the result field from API response into the target struct
func parseResult(result json.RawMessage, target interface{}) error {
	if len(result) == 0 {
		return fmt.Errorf("empty result")
	}
	return json.Unmarshal(result, target)
}
