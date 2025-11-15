package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// APIResponse represents a generic API response wrapper
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
		return NewRateLimitError(r.Description)
	}

	// Check if it's an auth error
	if r.ErrorCode == 401 || r.ErrorCode == 403 {
		return NewAuthError(r.Description)
	}

	// Generic API error
	return NewAPIError(r.ErrorCode, r.Description, "API request failed")
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       json.RawMessage `json:"data"`
	TotalCount int             `json:"total_count"`
	HasMore    bool            `json:"has_more"`
}

// RateLimitInfo represents rate limiting information from API response headers
type RateLimitInfo struct {
	Limit     int       // Maximum number of requests allowed
	Remaining int       // Number of requests remaining
	Reset     time.Time // Time when the rate limit resets
	RetryAfter time.Duration // Duration to wait before retrying (for 429 responses)
}

// ParseRateLimitHeaders parses rate limit information from HTTP headers
func ParseRateLimitHeaders(headers map[string][]string) *RateLimitInfo {
	info := &RateLimitInfo{}

	// Parse X-RateLimit-Limit
	if limit := headers["X-Ratelimit-Limit"]; len(limit) > 0 {
		if val, err := strconv.Atoi(limit[0]); err == nil {
			info.Limit = val
		}
	}

	// Parse X-RateLimit-Remaining
	if remaining := headers["X-Ratelimit-Remaining"]; len(remaining) > 0 {
		if val, err := strconv.Atoi(remaining[0]); err == nil {
			info.Remaining = val
		}
	}

	// Parse X-RateLimit-Reset (Unix timestamp)
	if reset := headers["X-Ratelimit-Reset"]; len(reset) > 0 {
		if val, err := strconv.ParseInt(reset[0], 10, 64); err == nil {
			info.Reset = time.Unix(val, 0)
		}
	}

	// Parse Retry-After (seconds)
	if retryAfter := headers["Retry-After"]; len(retryAfter) > 0 {
		if val, err := strconv.Atoi(retryAfter[0]); err == nil {
			info.RetryAfter = time.Duration(val) * time.Second
		}
	}

	return info
}

// ShouldBackoff returns true if we're approaching rate limits
func (r *RateLimitInfo) ShouldBackoff() bool {
	if r.Limit == 0 {
		return false
	}
	// Back off if we have less than 10% of requests remaining
	threshold := float64(r.Limit) * 0.1
	return float64(r.Remaining) < threshold
}

// String returns a string representation of rate limit info
func (r *RateLimitInfo) String() string {
	return fmt.Sprintf("RateLimit{Limit: %d, Remaining: %d, Reset: %s, RetryAfter: %s}",
		r.Limit, r.Remaining, r.Reset.Format(time.RFC3339), r.RetryAfter)
}

// LogLevel represents the logging level
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// IsValid validates the log level
func (ll LogLevel) IsValid() bool {
	switch ll {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return true
	default:
		return false
	}
}

// String returns the string representation of LogLevel
func (ll LogLevel) String() string {
	return string(ll)
}