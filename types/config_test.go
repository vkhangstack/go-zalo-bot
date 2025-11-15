package types

import (
	"testing"
)

func TestEnvironment_IsValid(t *testing.T) {
	tests := []struct {
		name string
		env  Environment
		want bool
	}{
		{"valid development", Development, true},
		{"valid production", Production, true},
		{"invalid environment", Environment("invalid"), false},
		{"empty environment", Environment(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.env.IsValid(); got != tt.want {
				t.Errorf("Environment.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessageType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		mt   MessageType
		want bool
	}{
		{"valid text", MessageTypeText, true},
		{"valid image", MessageTypeImage, true},
		{"valid file", MessageTypeFile, true},
		{"valid template", MessageTypeTemplate, true},
		{"valid interactive", MessageTypeInteractive, true},
		{"invalid type", MessageType("invalid"), false},
		{"empty type", MessageType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mt.IsValid(); got != tt.want {
				t.Errorf("MessageType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				BotToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			},
			wantErr: false,
		},
		{
			name: "missing bot token",
			config: &Config{
				BaseURL: "https://api.example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid bot token format - no colon",
			config: &Config{
				BotToken: "invalidtoken",
			},
			wantErr: true,
		},
		{
			name: "invalid bot token format - empty bot ID",
			config: &Config{
				BotToken: ":ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			},
			wantErr: true,
		},
		{
			name: "invalid bot token format - non-numeric bot ID",
			config: &Config{
				BotToken: "abc123:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			},
			wantErr: true,
		},
		{
			name: "invalid bot token format - short token part",
			config: &Config{
				BotToken: "123456:ABC",
			},
			wantErr: true,
		},
		{
			name: "invalid environment",
			config: &Config{
				BotToken:    "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
				Environment: Environment("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Check defaults are set
				if tt.config.BaseURL == "" {
					t.Error("Config.BaseURL should be set to default")
				}
				if tt.config.Timeout == 0 {
					t.Error("Config.Timeout should be set to default")
				}
				if tt.config.Retries == 0 {
					t.Error("Config.Retries should be set to default")
				}
				if tt.config.Environment == "" {
					t.Error("Config.Environment should be set to default")
				}
				if tt.config.HTTPClient == nil {
					t.Error("Config.HTTPClient should be set to default")
				}
				if tt.config.RetryConfig == nil {
					t.Error("Config.RetryConfig should be set to default")
				}
			}
		})
	}
}

func TestMessageConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *MessageConfig
		wantErr bool
	}{
		{
			name: "valid text message",
			config: &MessageConfig{
				ChatID:      "chat123",
				Text:        "Hello",
				MessageType: MessageTypeText,
			},
			wantErr: false,
		},
		{
			name: "missing chat ID",
			config: &MessageConfig{
				Text:        "Hello",
				MessageType: MessageTypeText,
			},
			wantErr: true,
		},
		{
			name: "invalid message type",
			config: &MessageConfig{
				ChatID:      "chat123",
				Text:        "Hello",
				MessageType: MessageType("invalid"),
			},
			wantErr: true,
		},
		{
			name: "text message without text",
			config: &MessageConfig{
				ChatID:      "chat123",
				MessageType: MessageTypeText,
			},
			wantErr: true,
		},
		{
			name: "image message without text",
			config: &MessageConfig{
				ChatID:      "chat123",
				MessageType: MessageTypeImage,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && tt.config.MessageType == "" {
				// Check default message type is set
				if tt.config.MessageType != MessageTypeText {
					t.Error("MessageConfig.MessageType should be set to default")
				}
			}
		})
	}
}

func TestWebhookConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *WebhookConfig
		wantErr bool
	}{
		{
			name: "valid webhook config",
			config: &WebhookConfig{
				URL:         "https://example.com/webhook",
				SecretToken: "secret",
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			config: &WebhookConfig{
				SecretToken: "secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("WebhookConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateConfig_Validate(t *testing.T) {
	tests := []struct {
		name   string
		config *UpdateConfig
		want   *UpdateConfig
	}{
		{
			name: "valid config",
			config: &UpdateConfig{
				Offset:  10,
				Limit:   50,
				Timeout: 30,
			},
			want: &UpdateConfig{
				Offset:  10,
				Limit:   50,
				Timeout: 30,
			},
		},
		{
			name: "zero limit gets default",
			config: &UpdateConfig{
				Offset:  10,
				Limit:   0,
				Timeout: 30,
			},
			want: &UpdateConfig{
				Offset:  10,
				Limit:   100,
				Timeout: 30,
			},
		},
		{
			name: "limit too high gets capped",
			config: &UpdateConfig{
				Offset:  10,
				Limit:   200,
				Timeout: 30,
			},
			want: &UpdateConfig{
				Offset:  10,
				Limit:   100,
				Timeout: 30,
			},
		},
		{
			name: "negative timeout gets reset",
			config: &UpdateConfig{
				Offset:  10,
				Limit:   50,
				Timeout: -10,
			},
			want: &UpdateConfig{
				Offset:  10,
				Limit:   50,
				Timeout: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err != nil {
				t.Errorf("UpdateConfig.Validate() error = %v", err)
				return
			}

			if tt.config.Offset != tt.want.Offset {
				t.Errorf("UpdateConfig.Offset = %v, want %v", tt.config.Offset, tt.want.Offset)
			}
			if tt.config.Limit != tt.want.Limit {
				t.Errorf("UpdateConfig.Limit = %v, want %v", tt.config.Limit, tt.want.Limit)
			}
			if tt.config.Timeout != tt.want.Timeout {
				t.Errorf("UpdateConfig.Timeout = %v, want %v", tt.config.Timeout, tt.want.Timeout)
			}
		})
	}
}

func TestConfig_GetAPIEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		method   string
		expected string
	}{
		{
			name: "default base URL",
			config: &Config{
				BotToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
				BaseURL:  "https://bot-api.zapps.me",
			},
			method:   "sendMessage",
			expected: "https://bot-api.zapps.me/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/sendMessage",
		},
		{
			name: "custom base URL",
			config: &Config{
				BotToken: "789012:XYZ-GHI5678jklMn-abc90D3f4g567hi89",
				BaseURL:  "https://dev-bot-api.zapps.me",
			},
			method:   "getMe",
			expected: "https://dev-bot-api.zapps.me/bot789012:XYZ-GHI5678jklMn-abc90D3f4g567hi89/getMe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetAPIEndpoint(tt.method)
			if result != tt.expected {
				t.Errorf("Config.GetAPIEndpoint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestImageMessageConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ImageMessageConfig
		wantErr bool
	}{
		{
			name: "valid image message config",
			config: &ImageMessageConfig{
				ChatID:   "chat123",
				ImageURL: "https://example.com/image.jpg",
				Caption:  "Test image",
				MimeType: "image/jpeg",
			},
			wantErr: false,
		},
		{
			name: "missing chat ID",
			config: &ImageMessageConfig{
				ImageURL: "https://example.com/image.jpg",
			},
			wantErr: true,
		},
		{
			name: "missing image URL",
			config: &ImageMessageConfig{
				ChatID: "chat123",
			},
			wantErr: true,
		},
		{
			name: "invalid MIME type",
			config: &ImageMessageConfig{
				ChatID:   "chat123",
				ImageURL: "https://example.com/image.jpg",
				MimeType: "text/plain",
			},
			wantErr: true,
		},
		{
			name: "valid without MIME type",
			config: &ImageMessageConfig{
				ChatID:   "chat123",
				ImageURL: "https://example.com/image.jpg",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageMessageConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileMessageConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *FileMessageConfig
		wantErr bool
	}{
		{
			name: "valid file message config",
			config: &FileMessageConfig{
				ChatID:   "chat123",
				FileURL:  "https://example.com/file.pdf",
				FileName: "document.pdf",
				MimeType: "application/pdf",
				Size:     1024,
			},
			wantErr: false,
		},
		{
			name: "missing chat ID",
			config: &FileMessageConfig{
				FileURL:  "https://example.com/file.pdf",
				FileName: "document.pdf",
			},
			wantErr: true,
		},
		{
			name: "missing file URL",
			config: &FileMessageConfig{
				ChatID:   "chat123",
				FileName: "document.pdf",
			},
			wantErr: true,
		},
		{
			name: "missing file name",
			config: &FileMessageConfig{
				ChatID:  "chat123",
				FileURL: "https://example.com/file.pdf",
			},
			wantErr: true,
		},
		{
			name: "file too large",
			config: &FileMessageConfig{
				ChatID:   "chat123",
				FileURL:  "https://example.com/file.pdf",
				FileName: "document.pdf",
				Size:     60 * 1024 * 1024, // 60MB
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("FileMessageConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBotOptions(t *testing.T) {
	config := &Config{}

	// Test WithBaseURL
	WithBaseURL("https://custom.api.com")(config)
	if config.BaseURL != "https://custom.api.com" {
		t.Errorf("WithBaseURL() failed, got %v", config.BaseURL)
	}

	// Test WithDebug
	WithDebug()(config)
	if !config.Debug {
		t.Error("WithDebug() failed, debug should be true")
	}

	// Test WithEnvironment
	WithEnvironment(Development)(config)
	if config.Environment != Development {
		t.Errorf("WithEnvironment() failed, got %v", config.Environment)
	}

	// Test WithRetries
	WithRetries(5)(config)
	if config.Retries != 5 {
		t.Errorf("WithRetries() failed, got %v", config.Retries)
	}
}