package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vkhangstack/go-zalo-bot/auth"
	"github.com/vkhangstack/go-zalo-bot/types"
)

func setupTestMessageService(t *testing.T, botToken string) (*MessageService, *types.Config) {
	config := &types.Config{
		BotToken:    botToken,
		BaseURL:     "https://bot-api.zapps.me",
		Timeout:     30 * time.Second,
		Retries:     3,
		Environment: types.Production,
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
	}
	
	if err := config.Validate(); err != nil {
		t.Fatalf("config validation failed: %v", err)
	}
	
	authService, err := auth.NewAuthService(config)
	if err != nil {
		t.Fatalf("failed to create auth service: %v", err)
	}
	
	service := NewMessageService(authService, config.HTTPClient, config)
	return service, config
}

func TestMessageService_Send_Success(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL contains bot token
		expectedPath := "/bot" + botToken + "/sendMessage"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %v, want %v", r.URL.Path, expectedPath)
		}
		
		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("Request method = %v, want %v", r.Method, http.MethodPost)
		}
		
		// Parse request body
		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		
		// Verify payload
		if payload["chat_id"] != "user123" {
			t.Errorf("chat_id = %v, want user123", payload["chat_id"])
		}
		if payload["text"] != "Hello, World!" {
			t.Errorf("text = %v, want Hello, World!", payload["text"])
		}
		
		// Return success response
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "text": "Hello, World!", "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	messageConfig := types.MessageConfig{
		ChatID: "user123",
		Text:   "Hello, World!",
	}
	
	ctx := context.Background()
	message, err := service.Send(ctx, messageConfig)
	
	if err != nil {
		t.Errorf("Send() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("Send() returned nil message")
	}
	
	if message.MessageID != "msg123" {
		t.Errorf("Message.MessageID = %v, want msg123", message.MessageID)
	}
}

func TestMessageService_Send_UnicodeSupport(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		
		// Verify Vietnamese text is properly handled
		expectedText := "Xin chào! Đây là tin nhắn tiếng Việt"
		if payload["text"] != expectedText {
			t.Errorf("text = %v, want %v", payload["text"], expectedText)
		}
		
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "text": "` + expectedText + `", "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	messageConfig := types.MessageConfig{
		ChatID: "user123",
		Text:   "Xin chào! Đây là tin nhắn tiếng Việt",
	}
	
	ctx := context.Background()
	message, err := service.Send(ctx, messageConfig)
	
	if err != nil {
		t.Errorf("Send() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("Send() returned nil message")
	}
}

func TestMessageService_Send_ValidationError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	service, _ := setupTestMessageService(t, botToken)
	
	tests := []struct {
		name    string
		config  types.MessageConfig
		wantErr bool
	}{
		{
			name: "missing chat ID",
			config: types.MessageConfig{
				Text: "Hello",
			},
			wantErr: true,
		},
		{
			name: "empty text for text message",
			config: types.MessageConfig{
				ChatID:      "user123",
				MessageType: types.MessageTypeText,
			},
			wantErr: true,
		},
		{
			name: "invalid attachment type",
			config: types.MessageConfig{
				ChatID: "user123",
				Attachments: []types.Attachment{
					{
						Type: types.AttachmentType("invalid"),
					},
				},
			},
			wantErr: true,
		},
	}
	
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Send(ctx, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageService_SendImage_Success(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		
		// Verify attachments
		attachments, ok := payload["attachments"].([]interface{})
		if !ok || len(attachments) == 0 {
			t.Error("Expected attachments in payload")
		}
		
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	imageConfig := types.ImageMessageConfig{
		ChatID:   "user123",
		ImageURL: "https://example.com/image.jpg",
		Caption:  "Test image",
		MimeType: "image/jpeg",
	}
	
	ctx := context.Background()
	message, err := service.SendImage(ctx, imageConfig)
	
	if err != nil {
		t.Errorf("SendImage() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("SendImage() returned nil message")
	}
}

func TestMessageService_SendImage_ValidationError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	service, _ := setupTestMessageService(t, botToken)
	
	tests := []struct {
		name    string
		config  types.ImageMessageConfig
		wantErr bool
	}{
		{
			name: "missing chat ID",
			config: types.ImageMessageConfig{
				ImageURL: "https://example.com/image.jpg",
			},
			wantErr: true,
		},
		{
			name: "missing image URL",
			config: types.ImageMessageConfig{
				ChatID: "user123",
			},
			wantErr: true,
		},
		{
			name: "invalid MIME type",
			config: types.ImageMessageConfig{
				ChatID:   "user123",
				ImageURL: "https://example.com/image.jpg",
				MimeType: "application/pdf",
			},
			wantErr: true,
		},
	}
	
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.SendImage(ctx, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageService_SendFile_Success(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	fileConfig := types.FileMessageConfig{
		ChatID:   "user123",
		FileURL:  "https://example.com/document.pdf",
		FileName: "document.pdf",
		MimeType: "application/pdf",
		Size:     1024,
	}
	
	ctx := context.Background()
	message, err := service.SendFile(ctx, fileConfig)
	
	if err != nil {
		t.Errorf("SendFile() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("SendFile() returned nil message")
	}
}

func TestMessageService_SendFile_ValidationError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	service, _ := setupTestMessageService(t, botToken)
	
	tests := []struct {
		name    string
		config  types.FileMessageConfig
		wantErr bool
	}{
		{
			name: "missing chat ID",
			config: types.FileMessageConfig{
				FileURL:  "https://example.com/file.pdf",
				FileName: "file.pdf",
			},
			wantErr: true,
		},
		{
			name: "missing file URL",
			config: types.FileMessageConfig{
				ChatID:   "user123",
				FileName: "file.pdf",
			},
			wantErr: true,
		},
		{
			name: "missing file name",
			config: types.FileMessageConfig{
				ChatID:  "user123",
				FileURL: "https://example.com/file.pdf",
			},
			wantErr: true,
		},
		{
			name: "file size exceeds limit",
			config: types.FileMessageConfig{
				ChatID:   "user123",
				FileURL:  "https://example.com/file.pdf",
				FileName: "file.pdf",
				Size:     51 * 1024 * 1024, // 51MB
			},
			wantErr: true,
		},
	}
	
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.SendFile(ctx, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageService_SendVideo_Success(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	ctx := context.Background()
	message, err := service.SendVideo(ctx, "user123", "https://example.com/video.mp4", "video/mp4")
	
	if err != nil {
		t.Errorf("SendVideo() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("SendVideo() returned nil message")
	}
}

func TestMessageService_SendVideo_ValidationError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	service, _ := setupTestMessageService(t, botToken)
	
	tests := []struct {
		name     string
		chatID   string
		videoURL string
		mimeType string
		wantErr  bool
	}{
		{
			name:     "missing chat ID",
			chatID:   "",
			videoURL: "https://example.com/video.mp4",
			mimeType: "video/mp4",
			wantErr:  true,
		},
		{
			name:     "missing video URL",
			chatID:   "user123",
			videoURL: "",
			mimeType: "video/mp4",
			wantErr:  true,
		},
		{
			name:     "invalid MIME type",
			chatID:   "user123",
			videoURL: "https://example.com/video.mp4",
			mimeType: "application/pdf",
			wantErr:  true,
		},
	}
	
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.SendVideo(ctx, tt.chatID, tt.videoURL, tt.mimeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendVideo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageService_SendAudio_Success(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	ctx := context.Background()
	message, err := service.SendAudio(ctx, "user123", "https://example.com/audio.mp3", "audio/mpeg")
	
	if err != nil {
		t.Errorf("SendAudio() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("SendAudio() returned nil message")
	}
}

func TestMessageService_SendAudio_ValidationError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	service, _ := setupTestMessageService(t, botToken)
	
	tests := []struct {
		name     string
		chatID   string
		audioURL string
		mimeType string
		wantErr  bool
	}{
		{
			name:     "missing chat ID",
			chatID:   "",
			audioURL: "https://example.com/audio.mp3",
			mimeType: "audio/mpeg",
			wantErr:  true,
		},
		{
			name:     "missing audio URL",
			chatID:   "user123",
			audioURL: "",
			mimeType: "audio/mpeg",
			wantErr:  true,
		},
		{
			name:     "invalid MIME type",
			chatID:   "user123",
			audioURL: "https://example.com/audio.mp3",
			mimeType: "video/mp4",
			wantErr:  true,
		},
	}
	
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.SendAudio(ctx, tt.chatID, tt.audioURL, tt.mimeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendAudio() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageService_SendTemplate_Success(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL
		expectedPath := "/bot" + botToken + "/sendTemplate"
		if r.URL.Path != expectedPath {
			t.Errorf("Request path = %v, want %v", r.URL.Path, expectedPath)
		}
		
		var payload map[string]interface{}
		json.NewDecoder(r.Body).Decode(&payload)
		
		// Verify structured message is present
		if _, ok := payload["structured_message"]; !ok {
			t.Error("Expected structured_message in payload")
		}
		
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	structuredConfig := types.StructuredMessageConfig{
		ChatID: "user123",
		StructuredMessage: types.StructuredMessage{
			Type: types.StructuredMessageTypeTemplate,
			Elements: []types.MessageElement{
				{
					Title:    "Test Element",
					Subtitle: "Test Subtitle",
					Buttons: []types.Button{
						{
							Type:    types.ButtonTypePostback,
							Title:   "Click Me",
							Payload: "test_payload",
						},
					},
				},
			},
			QuickReplies: []types.QuickReply{
				{
					ContentType: types.QuickReplyTypeText,
					Title:       "Quick Reply",
					Payload:     "qr_payload",
				},
			},
		},
	}
	
	ctx := context.Background()
	message, err := service.SendTemplate(ctx, structuredConfig)
	
	if err != nil {
		t.Errorf("SendTemplate() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("SendTemplate() returned nil message")
	}
}

func TestMessageService_SendTemplate_ValidationError(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	service, _ := setupTestMessageService(t, botToken)
	
	tests := []struct {
		name    string
		config  types.StructuredMessageConfig
		wantErr bool
	}{
		{
			name: "missing chat ID",
			config: types.StructuredMessageConfig{
				StructuredMessage: types.StructuredMessage{
					Type: types.StructuredMessageTypeTemplate,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid structured message type",
			config: types.StructuredMessageConfig{
				ChatID: "user123",
				StructuredMessage: types.StructuredMessage{
					Type: types.StructuredMessageType("invalid"),
				},
			},
			wantErr: true,
		},
		{
			name: "invalid element",
			config: types.StructuredMessageConfig{
				ChatID: "user123",
				StructuredMessage: types.StructuredMessage{
					Type: types.StructuredMessageTypeTemplate,
					Elements: []types.MessageElement{
						{
							Title: "", // Empty title
						},
					},
				},
			},
			wantErr: true,
		},
	}
	
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.SendTemplate(ctx, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendTemplate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageService_SendStructuredMessage_Alias(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	structuredConfig := types.StructuredMessageConfig{
		ChatID: "user123",
		StructuredMessage: types.StructuredMessage{
			Type: types.StructuredMessageTypeButton,
			Elements: []types.MessageElement{
				{
					Title: "Test",
				},
			},
		},
	}
	
	ctx := context.Background()
	message, err := service.SendStructuredMessage(ctx, structuredConfig)
	
	if err != nil {
		t.Errorf("SendStructuredMessage() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("SendStructuredMessage() returned nil message")
	}
}

func TestMessageService_SendImageMessage_Alias(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	imageConfig := types.ImageMessageConfig{
		ChatID:   "user123",
		ImageURL: "https://example.com/image.jpg",
	}
	
	ctx := context.Background()
	message, err := service.SendImageMessage(ctx, imageConfig)
	
	if err != nil {
		t.Errorf("SendImageMessage() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("SendImageMessage() returned nil message")
	}
}

func TestMessageService_SendFileMessage_Alias(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := APIResponse{
			OK:     true,
			Result: json.RawMessage(`{"message_id": "msg123", "chat": {"id": "user123", "type": "private"}, "date": "2024-01-01T00:00:00Z"}`),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	
	service, config := setupTestMessageService(t, botToken)
	config.BaseURL = server.URL
	authService, _ := auth.NewAuthService(config)
	service.authService = authService
	
	fileConfig := types.FileMessageConfig{
		ChatID:   "user123",
		FileURL:  "https://example.com/file.pdf",
		FileName: "file.pdf",
	}
	
	ctx := context.Background()
	message, err := service.SendFileMessage(ctx, fileConfig)
	
	if err != nil {
		t.Errorf("SendFileMessage() error = %v", err)
	}
	
	if message == nil {
		t.Fatal("SendFileMessage() returned nil message")
	}
}

func TestValidateRecipientID(t *testing.T) {
	tests := []struct {
		name        string
		recipientID string
		wantErr     bool
	}{
		{
			name:        "valid recipient ID",
			recipientID: "user123",
			wantErr:     false,
		},
		{
			name:        "empty recipient ID",
			recipientID: "",
			wantErr:     true,
		},
		{
			name:        "recipient ID too long",
			recipientID: string(make([]byte, 101)),
			wantErr:     true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRecipientID(tt.recipientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRecipientID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
