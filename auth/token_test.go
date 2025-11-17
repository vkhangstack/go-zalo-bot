package auth

import (
	"testing"
)

func TestNewTokenManager(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	tm := NewTokenManager(botToken)

	if tm.GetToken() != botToken {
		t.Errorf("Expected token %s, got %s", botToken, tm.GetToken())
	}

	if !tm.IsValid() {
		t.Error("Expected token to be valid")
	}
}

func TestTokenManager_ValidateToken(t *testing.T) {
	tm := NewTokenManager("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "too short token",
			token:   "short",
			wantErr: true,
		},
		{
			name:    "no colon separator",
			token:   "123456ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr: true,
		},
		{
			name:    "empty bot ID",
			token:   ":ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr: true,
		},
		{
			name:    "non-numeric bot ID",
			token:   "abc123:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr: true,
		},
		{
			name:    "token part too short",
			token:   "123456:ABC",
			wantErr: true,
		},
		{
			name:    "invalid characters in token part",
			token:   "123456:ABC@DEF#GHI",
			wantErr: true,
		},
		{
			name:    "valid bot token",
			token:   "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr: false,
		},
		{
			name:    "valid bot token with underscores",
			token:   "789012:XYZ_GHI5678jklMn_abc90D3f4g567hi89",
			wantErr: false,
		},
		{
			name:    "token with whitespace",
			token:   "  123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tm.ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTokenManager_SetToken(t *testing.T) {
	tm := NewTokenManager("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")

	// Test setting valid token
	newToken := "789012:XYZ-GHI5678jklMn-abc90D3f4g567hi89"
	err := tm.SetToken(newToken)
	if err != nil {
		t.Errorf("SetToken() error = %v", err)
	}

	if tm.GetToken() != newToken {
		t.Errorf("Expected token %s, got %s", newToken, tm.GetToken())
	}

	// Test setting invalid token
	err = tm.SetToken("")
	if err == nil {
		t.Error("Expected error for empty token")
	}

	// Token should remain unchanged after failed set
	if tm.GetToken() != newToken {
		t.Errorf("Token should remain unchanged after failed set")
	}
}

func TestTokenManager_IsValid(t *testing.T) {
	// Test with valid bot token
	tm := NewTokenManager("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
	if !tm.IsValid() {
		t.Error("Expected token to be valid")
	}

	// Test with empty token
	tm = NewTokenManager("")
	if tm.IsValid() {
		t.Error("Expected empty token to be invalid")
	}

	// Test with invalid format token
	tm = NewTokenManager("invalid-token")
	if tm.IsValid() {
		t.Error("Expected invalid format token to be invalid")
	}
}

func TestTokenManager_Clear(t *testing.T) {
	tm := NewTokenManager("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")

	// Verify token is set
	if tm.GetToken() == "" {
		t.Error("Expected token to be set")
	}

	// Clear token
	tm.Clear()

	// Verify token is cleared
	if tm.GetToken() != "" {
		t.Error("Expected token to be cleared")
	}

	if tm.IsValid() {
		t.Error("Expected token to be invalid after clear")
	}
}

func TestTokenManager_GetBotToken(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	tm := NewTokenManager(botToken)

	if tm.GetBotToken() != botToken {
		t.Errorf("Expected bot token %s, got %s", botToken, tm.GetBotToken())
	}
}
