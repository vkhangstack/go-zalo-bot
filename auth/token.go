package auth

import (
	"regexp"
	"strings"
	"sync"

	"github.com/vkhangstack/go-zalo-bot/types"
)

// TokenManager manages bot tokens for Zalo Bot API
// Bot tokens are long-lived and don't expire automatically
type TokenManager struct {
	botToken string
	mutex    sync.RWMutex
}

// NewTokenManager creates a new token manager for bot tokens
func NewTokenManager(botToken string) *TokenManager {
	return &TokenManager{
		botToken: botToken,
	}
}

// GetToken returns the current bot token
func (tm *TokenManager) GetToken() string {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	return tm.botToken
}

// SetToken sets a new bot token
func (tm *TokenManager) SetToken(botToken string) error {
	if err := tm.ValidateToken(botToken); err != nil {
		return err
	}
	
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	tm.botToken = botToken
	
	return nil
}

// ValidateToken validates the format and content of a bot token
func (tm *TokenManager) ValidateToken(token string) error {
	if token == "" {
		return types.NewValidationError("bot token cannot be empty")
	}
	
	// Remove any whitespace
	token = strings.TrimSpace(token)
	
	// Bot token should be in format: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	// At minimum, it should contain a colon separating bot ID and token
	if len(token) < 10 {
		return types.NewValidationError("bot token is too short")
	}
	
	colonIndex := -1
	for i, char := range token {
		if char == ':' {
			colonIndex = i
			break
		}
	}
	
	// Must have a colon and both parts should be non-empty
	if colonIndex <= 0 || colonIndex >= len(token)-1 {
		return types.NewValidationError("invalid bot token format - must contain ':' separator")
	}
	
	// Bot ID part should be numeric
	botIDPart := token[:colonIndex]
	for _, char := range botIDPart {
		if char < '0' || char > '9' {
			return types.NewValidationError("invalid bot token format - bot ID must be numeric")
		}
	}
	
	// Token part should be at least 10 characters
	tokenPart := token[colonIndex+1:]
	if len(tokenPart) < 10 {
		return types.NewValidationError("invalid bot token format - token part too short")
	}
	
	// Check for valid characters in token part (alphanumeric, hyphens, underscores)
	validTokenRegex := regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
	if !validTokenRegex.MatchString(tokenPart) {
		return types.NewValidationError("invalid bot token format - token part contains invalid characters")
	}
	
	return nil
}

// IsValid checks if the current bot token is valid
// Bot tokens don't expire, so we only check format validity
func (tm *TokenManager) IsValid() bool {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()
	
	if tm.botToken == "" {
		return false
	}
	
	return tm.ValidateToken(tm.botToken) == nil
}

// GetBotToken returns the bot token for URL embedding
func (tm *TokenManager) GetBotToken() string {
	return tm.GetToken()
}

// Clear clears the stored bot token
func (tm *TokenManager) Clear() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	
	tm.botToken = ""
}