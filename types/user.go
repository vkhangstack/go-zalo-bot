package types

import "encoding/json"

// User represents a user in Zalo Bot
type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	IsBot  bool   `json:"is_bot"`
}

// Chat represents a chat in Zalo Bot
type Chat struct {
	ID   string   `json:"id"`
	Type ChatType `json:"type"`
}

// ChatType represents the type of chat
type ChatType string

const (
	ChatTypePrivate ChatType = "private"
	ChatTypeGroup   ChatType = "group"
)

// IsValid validates the chat type
func (ct ChatType) IsValid() bool {
	switch ct {
	case ChatTypePrivate, ChatTypeGroup:
		return true
	default:
		return false
	}
}

// String returns the string representation of ChatType
func (ct ChatType) String() string {
	return string(ct)
}

// MarshalJSON implements json.Marshaler interface
func (ct ChatType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(ct))
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ct *ChatType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*ct = ChatType(s)
	return nil
}

// UserProfile represents extended user profile information
type UserProfile struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	// Additional profile fields based on Zalo API
}