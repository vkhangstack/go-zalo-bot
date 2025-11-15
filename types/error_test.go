package types

import (
	"testing"
	"time"
)

func TestZaloBotError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *ZaloBotError
		want string
	}{
		{
			name: "error with description",
			err: &ZaloBotError{
				Code:        400,
				Message:     "Bad Request",
				Description: "Invalid parameter",
				Type:        ErrorTypeValidation,
			},
			want: "Zalo Bot API Error 400: Bad Request - Invalid parameter",
		},
		{
			name: "error without description",
			err: &ZaloBotError{
				Code:    401,
				Message: "Unauthorized",
				Type:    ErrorTypeAuth,
			},
			want: "Zalo Bot API Error 401: Unauthorized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("ZaloBotError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZaloBotError_IsRetryable(t *testing.T) {
	tests := []struct {
		name string
		err  *ZaloBotError
		want bool
	}{
		{
			name: "network error is retryable",
			err: &ZaloBotError{
				Type: ErrorTypeNetwork,
			},
			want: true,
		},
		{
			name: "rate limit error is retryable",
			err: &ZaloBotError{
				Type: ErrorTypeRateLimit,
			},
			want: true,
		},
		{
			name: "5xx API error is retryable",
			err: &ZaloBotError{
				Code: 500,
				Type: ErrorTypeAPI,
			},
			want: true,
		},
		{
			name: "4xx API error is not retryable",
			err: &ZaloBotError{
				Code: 400,
				Type: ErrorTypeAPI,
			},
			want: false,
		},
		{
			name: "auth error is not retryable",
			err: &ZaloBotError{
				Type: ErrorTypeAuth,
			},
			want: false,
		},
		{
			name: "validation error is not retryable",
			err: &ZaloBotError{
				Type: ErrorTypeValidation,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.IsRetryable(); got != tt.want {
				t.Errorf("ZaloBotError.IsRetryable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryConfig_ShouldRetry(t *testing.T) {
	config := DefaultRetryConfig()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "network error should retry",
			err: &ZaloBotError{
				Type: ErrorTypeNetwork,
			},
			want: true,
		},
		{
			name: "rate limit error should retry",
			err: &ZaloBotError{
				Type: ErrorTypeRateLimit,
			},
			want: true,
		},
		{
			name: "auth error should not retry",
			err: &ZaloBotError{
				Type: ErrorTypeAuth,
			},
			want: false,
		},
		{
			name: "validation error should not retry",
			err: &ZaloBotError{
				Type: ErrorTypeValidation,
			},
			want: false,
		},
		{
			name: "5xx API error should retry",
			err: &ZaloBotError{
				Code: 500,
				Type: ErrorTypeAPI,
			},
			want: true,
		},
		{
			name: "4xx API error should not retry",
			err: &ZaloBotError{
				Code: 400,
				Type: ErrorTypeAPI,
			},
			want: false,
		},
		{
			name: "non-ZaloBotError should retry",
			err:  &testError{},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.ShouldRetry(tt.err); got != tt.want {
				t.Errorf("RetryConfig.ShouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetryConfig_NextDelay(t *testing.T) {
	config := &RetryConfig{
		InitialDelay:  1 * time.Second,
		MaxDelay:      10 * time.Second,
		BackoffFactor: 2.0,
	}

	tests := []struct {
		name    string
		attempt int
		want    time.Duration
	}{
		{"first attempt", 0, 1 * time.Second},
		{"second attempt", 1, 2 * time.Second},
		{"third attempt", 2, 4 * time.Second},
		{"fourth attempt", 3, 8 * time.Second},
		{"max delay reached", 4, 10 * time.Second},
		{"max delay maintained", 5, 10 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.NextDelay(tt.attempt); got != tt.want {
				t.Errorf("RetryConfig.NextDelay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() *ZaloBotError
		wantType ErrorType
		wantCode int
	}{
		{
			name: "NewAPIError",
			fn: func() *ZaloBotError {
				return NewAPIError(400, "Bad Request", "Invalid parameter")
			},
			wantType: ErrorTypeAPI,
			wantCode: 400,
		},
		{
			name: "NewNetworkError",
			fn: func() *ZaloBotError {
				return NewNetworkError("Connection failed")
			},
			wantType: ErrorTypeNetwork,
			wantCode: 0,
		},
		{
			name: "NewAuthError",
			fn: func() *ZaloBotError {
				return NewAuthError("Invalid token")
			},
			wantType: ErrorTypeAuth,
			wantCode: 401,
		},
		{
			name: "NewValidationError",
			fn: func() *ZaloBotError {
				return NewValidationError("Invalid input")
			},
			wantType: ErrorTypeValidation,
			wantCode: 400,
		},
		{
			name: "NewRateLimitError",
			fn: func() *ZaloBotError {
				return NewRateLimitError("Rate limit exceeded")
			},
			wantType: ErrorTypeRateLimit,
			wantCode: 429,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()
			if err.Type != tt.wantType {
				t.Errorf("Error type = %v, want %v", err.Type, tt.wantType)
			}
			if err.Code != tt.wantCode {
				t.Errorf("Error code = %v, want %v", err.Code, tt.wantCode)
			}
		})
	}
}

// testError is a helper type for testing non-ZaloBotError cases
type testError struct{}

func (e *testError) Error() string {
	return "test error"
}