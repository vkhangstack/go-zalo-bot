package types

import (
	"encoding/json"
	"time"
)

// Update represents an incoming update from Zalo Bot, either from a webhook
// delivery (Message/EventName populated) or from the getUpdates polling API
// (UpdateID/Message/PostbackEvent/UserAction populated).
type Update struct {
	UpdateID      int            `json:"update_id,omitempty"`
	EventName     string         `json:"event_name,omitempty"`
	Message       *Message       `json:"message,omitempty"`
	PostbackEvent *PostbackEvent `json:"postback,omitempty"`
	UserAction    *UserAction    `json:"user_action,omitempty"`
}

// PostbackEvent represents a postback event from button interactions
type PostbackEvent struct {
	Payload string `json:"payload"`
	Title   string `json:"title,omitempty"`
}

// UserAction represents a user action event
type UserAction struct {
	Type   UserActionType `json:"type"`
	UserID string         `json:"user_id"`
	Data   interface{}    `json:"data,omitempty"`
}

// UserActionType represents the type of user action
type UserActionType string

const (
	UserActionTypeJoin  UserActionType = "join"
	UserActionTypeLeave UserActionType = "leave"
	UserActionTypeBlock UserActionType = "block"
)

// IsValid validates the user action type
func (uat UserActionType) IsValid() bool {
	switch uat {
	case UserActionTypeJoin, UserActionTypeLeave, UserActionTypeBlock:
		return true
	default:
		return false
	}
}

// String returns the string representation of UserActionType
func (uat UserActionType) String() string {
	return string(uat)
}

// MarshalJSON implements json.Marshaler interface
func (uat UserActionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(uat))
}

// UnmarshalJSON implements json.Unmarshaler interface
func (uat *UserActionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*uat = UserActionType(s)
	return nil
}

// WebhookInfo represents information about the current webhook
type WebhookInfo struct {
	URL                  string    `json:"url"`
	HasCustomCertificate bool      `json:"has_custom_certificate"`
	PendingUpdateCount   int       `json:"pending_update_count"`
	LastErrorDate        time.Time `json:"last_error_date,omitempty"`
	LastErrorMessage     string    `json:"last_error_message,omitempty"`
	MaxConnections       int       `json:"max_connections,omitempty"`
	AllowedUpdates       []string  `json:"allowed_updates,omitempty"`
}

// Event names sent by Zalo in the webhook payload's result.event_name field.
// See https://bot.zapps.me/docs/webhook/
const (
	EventMessageText        = "message.text.received"
	EventMessageImage       = "message.image.received"
	EventMessageSticker     = "message.sticker.received"
	EventMessageVoice       = "message.voice.received"
	EventMessageUnsupported = "message.unsupported.received"
)

// WebhookPayload is the envelope Zalo POSTs to a configured webhook URL:
//
//	{"ok": true, "result": {"event_name": "message.text.received", "message": {...}}}
type WebhookPayload struct {
	OK     bool          `json:"ok"`
	Result WebhookResult `json:"result"`
}

// WebhookResult carries the event name and, for message events, the message content.
type WebhookResult struct {
	EventName string   `json:"event_name"`
	Message   *Message `json:"message,omitempty"`
}

// ParseWebhookPayload parses a raw webhook request body into a WebhookPayload.
func ParseWebhookPayload(data []byte) (*WebhookPayload, error) {
	var payload WebhookPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
