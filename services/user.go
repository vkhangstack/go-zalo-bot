package services

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/vkhangstack/go-zalo-bot/auth"
	"github.com/vkhangstack/go-zalo-bot/types"
)

// UserService handles user-related operations
type UserService struct {
	*BaseService
	authService *auth.AuthService
}

// NewUserService creates a new user service
func NewUserService(authService *auth.AuthService, client *http.Client, config *types.Config) *UserService {
	return &UserService{
		BaseService: NewBaseService(authService, client, config),
		authService: authService,
	}
}

// GetUserProfile retrieves user profile information by user ID
// Implements Requirements 4.1, 4.2
func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*types.UserProfile, error) {
	// Validate user ID format (Requirement 4.1)
	if err := validateUserID(userID); err != nil {
		return nil, err
	}

	// Prepare API request with bot token URL construction
	apiReq := &APIRequest{
		Method:    http.MethodGet,
		APIMethod: "getUserProfile",
		QueryParams: map[string]string{
			"user_id": userID,
		},
	}

	// Execute request with retry logic and rate limiting handling (Requirements 4.3, 4.4)
	resp, err := s.DoRequest(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	// Parse structured user data (Requirement 4.2)
	var userProfile types.UserProfile
	if err := json.Unmarshal(resp.Result, &userProfile); err != nil {
		return nil, types.NewAPIError(0, "failed to parse user profile", err.Error())
	}

	return &userProfile, nil
}

// validateUserID validates the user ID format
func validateUserID(userID string) error {
	if userID == "" {
		return types.NewValidationError("user ID is required")
	}

	// User ID should not be excessively long
	if len(userID) > 100 {
		return types.NewValidationError("user ID is too long")
	}

	return nil
}
