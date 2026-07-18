package services

import (
	"testing"

	"github.com/vkhangstack/go-zalo-bot/types"
)

const realWebhookSample = `{
	"ok": true,
	"result": {
		"message": {
			"from": {"id": "user1", "display_name": "Ted", "is_bot": false},
			"chat": {"id": "chat1", "chat_type": "PRIVATE"},
			"text": "hello",
			"message_id": "msg1",
			"date": 1750316131602
		},
		"event_name": "message.text.received"
	}
}`

func TestNewWebhookService(t *testing.T) {
	service := NewWebhookService(nil, "initial-token")

	if got := service.GetSecretToken(); got != "initial-token" {
		t.Errorf("GetSecretToken() = %v, want initial-token", got)
	}

	service.SetSecretToken("updated-token")
	if got := service.GetSecretToken(); got != "updated-token" {
		t.Errorf("GetSecretToken() after SetSecretToken() = %v, want updated-token", got)
	}
}

func TestWebhookService_ValidateSecretToken(t *testing.T) {
	secretToken := "test-secret-token"
	service := &WebhookService{
		secretToken: secretToken,
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{"valid token", secretToken, false},
		{"invalid token", "wrong-token", true},
		{"empty token", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateSecretToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSecretToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhookService_ParseUpdate(t *testing.T) {
	service := &WebhookService{}

	tests := []struct {
		name    string
		payload string
		wantErr bool
		check   func(*types.Update) bool
	}{
		{
			name:    "real webhook payload - text message",
			payload: realWebhookSample,
			wantErr: false,
			check: func(u *types.Update) bool {
				return u.EventName == types.EventMessageText &&
					u.Message != nil && u.Message.Text == "hello" &&
					u.Message.Chat != nil && u.Message.Chat.Type == types.ChatTypePrivate
			},
		},
		{
			name:    "real webhook payload - unsupported event",
			payload: `{"ok":true,"result":{"event_name":"message.unsupported.received"}}`,
			wantErr: false,
			check: func(u *types.Update) bool {
				return u.EventName == types.EventMessageUnsupported && u.Message == nil
			},
		},
		{
			name:    "missing event_name in result",
			payload: `{"ok":true,"result":{}}`,
			wantErr: true,
			check:   nil,
		},
		{
			name:    "empty payload",
			payload: ``,
			wantErr: true,
			check:   nil,
		},
		{
			name:    "invalid json",
			payload: `{invalid json}`,
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update, err := service.ParseUpdate([]byte(tt.payload))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				if !tt.check(update) {
					t.Errorf("ParseUpdate() validation failed for update: %+v", update)
				}
			}
		})
	}
}

func TestWebhookService_ProcessWebhook(t *testing.T) {
	secretToken := "test-secret"
	service := &WebhookService{
		secretToken: secretToken,
	}

	tests := []struct {
		name        string
		payload     []byte
		secretToken string
		wantErr     bool
	}{
		{
			name:        "valid webhook request",
			payload:     []byte(realWebhookSample),
			secretToken: secretToken,
			wantErr:     false,
		},
		{
			name:        "invalid secret token",
			payload:     []byte(realWebhookSample),
			secretToken: "wrong-secret",
			wantErr:     true,
		},
		{
			name:        "empty payload",
			payload:     []byte{},
			secretToken: secretToken,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update, err := service.ProcessWebhook(tt.payload, tt.secretToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && update == nil {
				t.Error("ProcessWebhook() returned nil update for valid request")
			}
		})
	}
}

func TestWebhookService_HandleWebhookEvent(t *testing.T) {
	secretToken := "test-secret"
	service := &WebhookService{
		secretToken: secretToken,
	}

	tests := []struct {
		name         string
		payload      []byte
		expectedType string
		wantErr      bool
	}{
		{
			name:         "message event",
			payload:      []byte(realWebhookSample),
			expectedType: "message",
			wantErr:      false,
		},
		{
			name:         "unsupported event has no message",
			payload:      []byte(`{"ok":true,"result":{"event_name":"message.unsupported.received"}}`),
			expectedType: "unknown",
			wantErr:      false,
		},
		{
			name:         "missing event_name in result",
			payload:      []byte(`{"ok":true,"result":{}}`),
			expectedType: "",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventType, data, err := service.HandleWebhookEvent(tt.payload, secretToken)

			if (err != nil) != tt.wantErr {
				t.Errorf("HandleWebhookEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if eventType != tt.expectedType {
					t.Errorf("HandleWebhookEvent() eventType = %v, want %v", eventType, tt.expectedType)
				}
				if data == nil {
					t.Error("HandleWebhookEvent() returned nil data")
				}
			}
		})
	}
}

func TestWebhookService_RejectInvalidRequest(t *testing.T) {
	service := &WebhookService{}

	tests := []struct {
		name   string
		reason string
	}{
		{
			name:   "with reason",
			reason: "invalid secret token",
		},
		{
			name:   "empty reason",
			reason: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RejectInvalidRequest(tt.reason)
			if err == nil {
				t.Error("RejectInvalidRequest() should return an error")
			}
		})
	}
}
