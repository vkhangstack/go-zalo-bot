package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUpdate_JSON(t *testing.T) {
	update := Update{
		UpdateID: 123,
		Message: &Message{
			MessageID: "msg123",
			From: &User{
				ID:   "user123",
				Name: "Test User",
			},
			Chat: &Chat{
				ID:   "chat123",
				Type: ChatTypePrivate,
			},
			Date: time.Now(),
			Text: "Hello",
		},
	}

	// Test marshal
	data, err := json.Marshal(update)
	if err != nil {
		t.Errorf("Update.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled Update
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Update.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.UpdateID != update.UpdateID {
		t.Errorf("Update.UpdateID = %v, want %v", unmarshaled.UpdateID, update.UpdateID)
	}
	if unmarshaled.Message.MessageID != update.Message.MessageID {
		t.Errorf("Update.Message.MessageID = %v, want %v", unmarshaled.Message.MessageID, update.Message.MessageID)
	}
}

func TestWebhookInfo_JSON(t *testing.T) {
	webhookInfo := WebhookInfo{
		URL:                  "https://example.com/webhook",
		HasCustomCertificate: true,
		PendingUpdateCount:   5,
		LastErrorDate:        time.Now(),
		LastErrorMessage:     "Connection timeout",
		MaxConnections:       40,
		AllowedUpdates:       []string{"message", "callback_query"},
	}

	// Test marshal
	data, err := json.Marshal(webhookInfo)
	if err != nil {
		t.Errorf("WebhookInfo.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled WebhookInfo
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("WebhookInfo.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.URL != webhookInfo.URL {
		t.Errorf("WebhookInfo.URL = %v, want %v", unmarshaled.URL, webhookInfo.URL)
	}
	if unmarshaled.HasCustomCertificate != webhookInfo.HasCustomCertificate {
		t.Errorf("WebhookInfo.HasCustomCertificate = %v, want %v", unmarshaled.HasCustomCertificate, webhookInfo.HasCustomCertificate)
	}
	if unmarshaled.PendingUpdateCount != webhookInfo.PendingUpdateCount {
		t.Errorf("WebhookInfo.PendingUpdateCount = %v, want %v", unmarshaled.PendingUpdateCount, webhookInfo.PendingUpdateCount)
	}
}

// realWebhookSample mirrors the sample payload from https://bot.zapps.me/docs/webhook/
const realWebhookSample = `{
	"ok": true,
	"result": {
		"message": {
			"from": {"id": "user1", "display_name": "Ted", "is_bot": false},
			"chat": {"id": "chat1", "chat_type": "PRIVATE"},
			"text": "Xin chao",
			"message_id": "msg1",
			"date": 1750316131602
		},
		"event_name": "message.text.received"
	}
}`

func TestParseWebhookPayload(t *testing.T) {
	payload, err := ParseWebhookPayload([]byte(realWebhookSample))
	if err != nil {
		t.Fatalf("ParseWebhookPayload() error = %v", err)
	}

	if !payload.OK {
		t.Errorf("WebhookPayload.OK = %v, want true", payload.OK)
	}
	if payload.Result.EventName != EventMessageText {
		t.Errorf("WebhookPayload.Result.EventName = %v, want %v", payload.Result.EventName, EventMessageText)
	}

	msg := payload.Result.Message
	if msg == nil {
		t.Fatal("WebhookPayload.Result.Message is nil, want non-nil")
	}
	if msg.MessageID != "msg1" {
		t.Errorf("Message.MessageID = %v, want %v", msg.MessageID, "msg1")
	}
	if msg.Text != "Xin chao" {
		t.Errorf("Message.Text = %v, want %v", msg.Text, "Xin chao")
	}
	if msg.From == nil || msg.From.ID != "user1" {
		t.Errorf("Message.From = %+v, want ID user1", msg.From)
	}
	if msg.Chat == nil || msg.Chat.Type != ChatTypePrivate {
		t.Errorf("Message.Chat = %+v, want Type %v", msg.Chat, ChatTypePrivate)
	}
	wantDate := time.UnixMilli(1750316131602)
	if !msg.Date.Equal(wantDate) {
		t.Errorf("Message.Date = %v, want %v", msg.Date, wantDate)
	}
}

func TestParseWebhookPayload_Unsupported(t *testing.T) {
	sample := `{"ok":true,"result":{"event_name":"message.unsupported.received"}}`

	payload, err := ParseWebhookPayload([]byte(sample))
	if err != nil {
		t.Fatalf("ParseWebhookPayload() error = %v", err)
	}
	if payload.Result.EventName != EventMessageUnsupported {
		t.Errorf("WebhookPayload.Result.EventName = %v, want %v", payload.Result.EventName, EventMessageUnsupported)
	}
	if payload.Result.Message != nil {
		t.Errorf("WebhookPayload.Result.Message = %+v, want nil", payload.Result.Message)
	}
}

func TestPostbackEvent_JSON(t *testing.T) {
	postback := PostbackEvent{
		Payload: "button_clicked",
		Title:   "Click Me",
	}

	// Test marshal
	data, err := json.Marshal(postback)
	if err != nil {
		t.Errorf("PostbackEvent.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled PostbackEvent
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("PostbackEvent.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.Payload != postback.Payload {
		t.Errorf("PostbackEvent.Payload = %v, want %v", unmarshaled.Payload, postback.Payload)
	}
	if unmarshaled.Title != postback.Title {
		t.Errorf("PostbackEvent.Title = %v, want %v", unmarshaled.Title, postback.Title)
	}
}

func TestUserAction_JSON(t *testing.T) {
	userAction := UserAction{
		Type:   UserActionTypeJoin,
		UserID: "user123",
		Data:   map[string]interface{}{"channel": "general"},
	}

	// Test marshal
	data, err := json.Marshal(userAction)
	if err != nil {
		t.Errorf("UserAction.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled UserAction
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("UserAction.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.Type != userAction.Type {
		t.Errorf("UserAction.Type = %v, want %v", unmarshaled.Type, userAction.Type)
	}
	if unmarshaled.UserID != userAction.UserID {
		t.Errorf("UserAction.UserID = %v, want %v", unmarshaled.UserID, userAction.UserID)
	}
}

func TestUserActionType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		uat  UserActionType
		want bool
	}{
		{"valid join", UserActionTypeJoin, true},
		{"valid leave", UserActionTypeLeave, true},
		{"valid block", UserActionTypeBlock, true},
		{"invalid type", UserActionType("invalid"), false},
		{"empty type", UserActionType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uat.IsValid(); got != tt.want {
				t.Errorf("UserActionType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserActionType_String(t *testing.T) {
	tests := []struct {
		name string
		uat  UserActionType
		want string
	}{
		{"join", UserActionTypeJoin, "join"},
		{"leave", UserActionTypeLeave, "leave"},
		{"block", UserActionTypeBlock, "block"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.uat.String(); got != tt.want {
				t.Errorf("UserActionType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserActionType_JSON(t *testing.T) {
	tests := []struct {
		name string
		uat  UserActionType
	}{
		{"join", UserActionTypeJoin},
		{"leave", UserActionTypeLeave},
		{"block", UserActionTypeBlock},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshal
			data, err := json.Marshal(tt.uat)
			if err != nil {
				t.Errorf("UserActionType.MarshalJSON() error = %v", err)
				return
			}

			// Test unmarshal
			var unmarshaled UserActionType
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Errorf("UserActionType.UnmarshalJSON() error = %v", err)
				return
			}

			if unmarshaled != tt.uat {
				t.Errorf("UserActionType after JSON round-trip = %v, want %v", unmarshaled, tt.uat)
			}
		})
	}
}

func TestUpdate_WithPostbackEvent(t *testing.T) {
	update := Update{
		UpdateID: 123,
		PostbackEvent: &PostbackEvent{
			Payload: "button_clicked",
			Title:   "Click Me",
		},
	}

	// Test marshal
	data, err := json.Marshal(update)
	if err != nil {
		t.Errorf("Update.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled Update
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Update.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.UpdateID != update.UpdateID {
		t.Errorf("Update.UpdateID = %v, want %v", unmarshaled.UpdateID, update.UpdateID)
	}
	if unmarshaled.PostbackEvent == nil {
		t.Error("Update.PostbackEvent is nil, want non-nil")
		return
	}
	if unmarshaled.PostbackEvent.Payload != update.PostbackEvent.Payload {
		t.Errorf("Update.PostbackEvent.Payload = %v, want %v", unmarshaled.PostbackEvent.Payload, update.PostbackEvent.Payload)
	}
}

func TestUpdate_WithUserAction(t *testing.T) {
	update := Update{
		UpdateID: 123,
		UserAction: &UserAction{
			Type:   UserActionTypeJoin,
			UserID: "user123",
			Data:   map[string]interface{}{"channel": "general"},
		},
	}

	// Test marshal
	data, err := json.Marshal(update)
	if err != nil {
		t.Errorf("Update.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled Update
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Update.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.UpdateID != update.UpdateID {
		t.Errorf("Update.UpdateID = %v, want %v", unmarshaled.UpdateID, update.UpdateID)
	}
	if unmarshaled.UserAction == nil {
		t.Error("Update.UserAction is nil, want non-nil")
		return
	}
	if unmarshaled.UserAction.Type != update.UserAction.Type {
		t.Errorf("Update.UserAction.Type = %v, want %v", unmarshaled.UserAction.Type, update.UserAction.Type)
	}
	if unmarshaled.UserAction.UserID != update.UserAction.UserID {
		t.Errorf("Update.UserAction.UserID = %v, want %v", unmarshaled.UserAction.UserID, update.UserAction.UserID)
	}
}
