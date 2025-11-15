package types

import (
	"encoding/json"
	"testing"
)

func TestChatType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		ct   ChatType
		want bool
	}{
		{"valid private", ChatTypePrivate, true},
		{"valid group", ChatTypeGroup, true},
		{"invalid type", ChatType("invalid"), false},
		{"empty type", ChatType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.IsValid(); got != tt.want {
				t.Errorf("ChatType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatType_String(t *testing.T) {
	tests := []struct {
		name string
		ct   ChatType
		want string
	}{
		{"private type", ChatTypePrivate, "private"},
		{"group type", ChatTypeGroup, "group"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ct.String(); got != tt.want {
				t.Errorf("ChatType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChatType_JSON(t *testing.T) {
	tests := []struct {
		name string
		ct   ChatType
		want string
	}{
		{"marshal private", ChatTypePrivate, `"private"`},
		{"marshal group", ChatTypeGroup, `"group"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.ct)
			if err != nil {
				t.Errorf("ChatType.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("ChatType.MarshalJSON() = %v, want %v", string(got), tt.want)
			}

			// Test unmarshal
			var ct ChatType
			if err := json.Unmarshal(got, &ct); err != nil {
				t.Errorf("ChatType.UnmarshalJSON() error = %v", err)
				return
			}
			if ct != tt.ct {
				t.Errorf("ChatType.UnmarshalJSON() = %v, want %v", ct, tt.ct)
			}
		})
	}
}

func TestUser_JSON(t *testing.T) {
	user := User{
		ID:     "user123",
		Name:   "Test User",
		Avatar: "https://example.com/avatar.jpg",
		IsBot:  false,
	}

	// Test marshal
	data, err := json.Marshal(user)
	if err != nil {
		t.Errorf("User.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled User
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("User.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.ID != user.ID {
		t.Errorf("User.ID = %v, want %v", unmarshaled.ID, user.ID)
	}
	if unmarshaled.Name != user.Name {
		t.Errorf("User.Name = %v, want %v", unmarshaled.Name, user.Name)
	}
	if unmarshaled.Avatar != user.Avatar {
		t.Errorf("User.Avatar = %v, want %v", unmarshaled.Avatar, user.Avatar)
	}
	if unmarshaled.IsBot != user.IsBot {
		t.Errorf("User.IsBot = %v, want %v", unmarshaled.IsBot, user.IsBot)
	}
}

func TestChat_JSON(t *testing.T) {
	chat := Chat{
		ID:   "chat123",
		Type: ChatTypePrivate,
	}

	// Test marshal
	data, err := json.Marshal(chat)
	if err != nil {
		t.Errorf("Chat.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled Chat
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Chat.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.ID != chat.ID {
		t.Errorf("Chat.ID = %v, want %v", unmarshaled.ID, chat.ID)
	}
	if unmarshaled.Type != chat.Type {
		t.Errorf("Chat.Type = %v, want %v", unmarshaled.Type, chat.Type)
	}
}