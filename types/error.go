package types

import (
	"fmt"
	"time"
)

// ZaloBotError represents an error from the Zalo Bot API
type ZaloBotError struct {
	Code        int       `json:"code"`
	Message     string    `json:"message"`
	Description string    `json:"description,omitempty"`
	Type        ErrorType `json:"type"`
}

// Error implements the error interface
func (e *ZaloBotError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("Zalo Bot API Error %d: %s - %s", e.Code, e.Message, e.Description)
	}
	return fmt.Sprintf("Zalo Bot API Error %d: %s", e.Code, e.Message)
}

// IsRetryable returns true if the error is retryable
func (e *ZaloBotError) IsRetryable() bool {
	switch e.Type {
	case ErrorTypeNetwork, ErrorTypeRateLimit:
		return true
	case ErrorTypeAPI:
		// Some API errors are retryable (5xx status codes)
		return e.Code >= 500 && e.Code < 600
	default:
		return false
	}
}

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeAPI        ErrorType = "api_error"
	ErrorTypeNetwork    ErrorType = "network_error"
	ErrorTypeAuth       ErrorType = "auth_error"
	ErrorTypeValidation ErrorType = "validation_error"
	ErrorTypeRateLimit  ErrorType = "rate_limit_error"
)

// String returns the string representation of ErrorType
func (et ErrorType) String() string {
	return string(et)
}

// RetryConfig represents configuration for retry logic
type RetryConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []ErrorType
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:    3,
		InitialDelay:  1 * time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		RetryableErrors: []ErrorType{
			ErrorTypeNetwork,
			ErrorTypeRateLimit,
		},
	}
}

// ShouldRetry returns true if the error should be retried
func (rc *RetryConfig) ShouldRetry(err error) bool {
	zaloBotErr, ok := err.(*ZaloBotError)
	if !ok {
		// Non-ZaloBotError, assume it's a network error and retry
		return true
	}
	
	for _, retryableType := range rc.RetryableErrors {
		if zaloBotErr.Type == retryableType {
			return true
		}
	}
	
	return zaloBotErr.IsRetryable()
}

// NextDelay calculates the next delay for retry attempt
func (rc *RetryConfig) NextDelay(attempt int) time.Duration {
	if attempt <= 0 {
		return rc.InitialDelay
	}
	
	delay := rc.InitialDelay
	for i := 0; i < attempt; i++ {
		delay = time.Duration(float64(delay) * rc.BackoffFactor)
		if delay > rc.MaxDelay {
			delay = rc.MaxDelay
			break
		}
	}
	
	return delay
}

// NewAPIError creates a new API error
func NewAPIError(code int, message, description string) *ZaloBotError {
	return &ZaloBotError{
		Code:        code,
		Message:     message,
		Description: description,
		Type:        ErrorTypeAPI,
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(message string) *ZaloBotError {
	return &ZaloBotError{
		Code:    0,
		Message: message,
		Type:    ErrorTypeNetwork,
	}
}

// NewAuthError creates a new authentication error
func NewAuthError(message string) *ZaloBotError {
	return &ZaloBotError{
		Code:    401,
		Message: message,
		Type:    ErrorTypeAuth,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ZaloBotError {
	return &ZaloBotError{
		Code:    400,
		Message: message,
		Type:    ErrorTypeValidation,
	}
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(message string) *ZaloBotError {
	return &ZaloBotError{
		Code:    429,
		Message: message,
		Type:    ErrorTypeRateLimit,
	}
}