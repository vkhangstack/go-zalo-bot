package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/vkhangstack/go-zalo-bot/types"
)

func TestWebhookService_ValidateSignature(t *testing.T) {
	secretToken := "test-secret-token"
	service := &WebhookService{
		secretToken: secretToken,
	}

	tests := []struct {
		name      string
		payload   []byte
		signature string
		wantErr   bool
	}{
		{
			name:      "valid signature",
			payload:   []byte(`{"update_id":1,"message":{"message_id":"123","text":"hello"}}`),
			signature: "valid-signature",
			wantErr:   false,
		},
		{
			name:      "invalid signature",
			payload:   []byte(`{"update_id":1,"message":{"message_id":"123","text":"hello"}}`),
			signature: "invalid-signature",
			wantErr:   true,
		},
		{
			name:      "empty signature",
			payload:   []byte(`{"update_id":1}`),
			signature: "",
			wantErr:   true,
		},
		{
			name:      "empty payload",
			payload:   []byte{},
			signature: "some-signature",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For valid signature test, compute the actual signature
			if tt.name == "valid signature" {
				tt.signature = computeSignature(tt.payload, secretToken)
			}

			err := service.ValidateSignature(tt.payload, tt.signature)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSignature() error = %v, wantErr %v", err, tt.wantErr)
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
			name:    "valid update with message",
			payload: `{"update_id":1,"message":{"message_id":"123","text":"hello"}}`,
			wantErr: false,
			check: func(u *types.Update) bool {
				return u.UpdateID == 1 && u.Message != nil && u.Message.Text == "hello"
			},
		},
		{
			name:    "valid update with postback",
			payload: `{"update_id":2,"postback":{"payload":"button_clicked","title":"Click Me"}}`,
			wantErr: false,
			check: func(u *types.Update) bool {
				return u.UpdateID == 2 && u.PostbackEvent != nil && u.PostbackEvent.Payload == "button_clicked"
			},
		},
		{
			name:    "webhook event format - message",
			payload: `{"event_name":"message","app_id":"123","user_id":"user1","oa_id":"oa1","timestamp":1234567890,"data":{"message_id":"msg1","user_id":"user1","text":"test message","timestamp":1234567890}}`,
			wantErr: false,
			check: func(u *types.Update) bool {
				// For webhook event format, just check that update was created
				// The conversion logic will handle the details
				return u != nil && (u.UpdateID > 0 || u.Message != nil)
			},
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

	validPayload := []byte(`{"update_id":1,"message":{"message_id":"123","text":"hello"}}`)
	validSignature := computeSignature(validPayload, secretToken)

	tests := []struct {
		name      string
		payload   []byte
		signature string
		wantErr   bool
	}{
		{
			name:      "valid webhook request",
			payload:   validPayload,
			signature: validSignature,
			wantErr:   false,
		},
		{
			name:      "invalid signature",
			payload:   validPayload,
			signature: "wrong-signature",
			wantErr:   true,
		},
		{
			name:      "empty payload",
			payload:   []byte{},
			signature: validSignature,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update, err := service.ProcessWebhook(tt.payload, tt.signature)
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
		name          string
		payload       []byte
		expectedType  string
		wantErr       bool
	}{
		{
			name:         "message event",
			payload:      []byte(`{"update_id":1,"message":{"message_id":"123","text":"hello"}}`),
			expectedType: "message",
			wantErr:      false,
		},
		{
			name:         "postback event",
			payload:      []byte(`{"update_id":2,"postback":{"payload":"clicked"}}`),
			expectedType: "postback",
			wantErr:      false,
		},
		{
			name:         "user action event",
			payload:      []byte(`{"update_id":3,"user_action":{"type":"join","user_id":"user1"}}`),
			expectedType: "user_action",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signature := computeSignature(tt.payload, secretToken)
			eventType, data, err := service.HandleWebhookEvent(tt.payload, signature)

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
			reason: "invalid signature",
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

func TestWebhookService_ConvertWebhookEventToUpdate(t *testing.T) {
	service := &WebhookService{}

	tests := []struct {
		name    string
		event   *types.WebhookEvent
		wantErr bool
		check   func(*types.Update) bool
	}{
		{
			name: "message event",
			event: &types.WebhookEvent{
				EventName: "message",
				Timestamp: 1234567890,
				Data:      json.RawMessage(`{"message_id":"msg1","user_id":"user1","text":"test"}`),
			},
			wantErr: false,
			check: func(u *types.Update) bool {
				return u.Message != nil && u.Message.Text == "test"
			},
		},
		{
			name: "postback event",
			event: &types.WebhookEvent{
				EventName: "postback",
				Timestamp: 1234567890,
				Data:      json.RawMessage(`{"payload":"clicked","title":"Button"}`),
			},
			wantErr: false,
			check: func(u *types.Update) bool {
				return u.PostbackEvent != nil && u.PostbackEvent.Payload == "clicked"
			},
		},
		{
			name: "user action event",
			event: &types.WebhookEvent{
				EventName: "user_action",
				Timestamp: 1234567890,
				Data:      json.RawMessage(`{"user_id":"user1","action":"join"}`),
			},
			wantErr: false,
			check: func(u *types.Update) bool {
				return u.UserAction != nil && u.UserAction.UserID == "user1"
			},
		},
		{
			name: "unsupported event",
			event: &types.WebhookEvent{
				EventName: "unknown_event",
				Timestamp: 1234567890,
				Data:      json.RawMessage(`{}`),
			},
			wantErr: true,
			check:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			update, err := service.convertWebhookEventToUpdate(tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertWebhookEventToUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.check != nil {
				if !tt.check(update) {
					t.Errorf("convertWebhookEventToUpdate() validation failed")
				}
			}
		})
	}
}

// Helper function to compute HMAC signature
func computeSignature(payload []byte, secret string) string {
	// Import crypto packages at the top if needed
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
