package types

import (
	"encoding/json"
	"time"
)

// Update represents an incoming update from Zalo Bot
type Update struct {
	UpdateID      int            `json:"update_id"`
	Message       *Message       `json:"message,omitempty"`
	PostbackEvent *PostbackEvent `json:"postback,omitempty"`
	UserAction    *UserAction    `json:"user_action,omitempty"`
	// Other event types as supported by Zalo (Requirement 3.5)
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

// WebhookEvent represents a webhook event payload
type WebhookEvent struct {
	EventName string          `json:"event_name"`
	AppID     string          `json:"app_id"`
	UserID    string          `json:"user_id"`
	OAID      string          `json:"oa_id"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// MessageEvent represents a message event in webhook
type MessageEvent struct {
	MessageID   string       `json:"message_id"`
	UserID      string       `json:"user_id"`
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Timestamp   int64        `json:"timestamp"`
}

// UserActionEvent represents a user action event in webhook
type UserActionEvent struct {
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}

// ParseWebhookEvent parses a webhook event from JSON data
func ParseWebhookEvent(data []byte) (*WebhookEvent, error) {
	var event WebhookEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// ParseMessageEvent parses message event data
func (we *WebhookEvent) ParseMessageEvent() (*MessageEvent, error) {
	var msgEvent MessageEvent
	if err := json.Unmarshal(we.Data, &msgEvent); err != nil {
		return nil, err
	}
	return &msgEvent, nil
}

// ParseUserActionEvent parses user action event data
func (we *WebhookEvent) ParseUserActionEvent() (*UserActionEvent, error) {
	var actionEvent UserActionEvent
	if err := json.Unmarshal(we.Data, &actionEvent); err != nil {
		return nil, err
	}
	return &actionEvent, nil
}