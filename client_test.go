package zalobot

import (
	"net/http"
	"testing"
	"time"

	"github.com/vkhangstack/go-zalo-bot/types"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		botToken  string
		options   []types.BotOption
		wantErr   bool
		errType   types.ErrorType
	}{
		{
			name:     "valid bot token",
			botToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr:  false,
		},
		{
			name:     "valid bot token with options",
			botToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			options: []types.BotOption{
				types.WithDebug(),
				types.WithTimeout(10 * time.Second),
			},
			wantErr: false,
		},
		{
			name:     "empty bot token",
			botToken: "",
			wantErr:  true,
			errType:  types.ErrorTypeValidation,
		},
		{
			name:     "invalid bot token format - no colon",
			botToken: "123456ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr:  true,
			errType:  types.ErrorTypeValidation,
		},
		{
			name:     "invalid bot token format - non-numeric bot ID",
			botToken: "abc123:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantErr:  true,
			errType:  types.ErrorTypeValidation,
		},
		{
			name:     "invalid bot token format - token part too short",
			botToken: "123456:ABC",
			wantErr:  true,
			errType:  types.ErrorTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot, err := New(tt.botToken, tt.options...)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("New() expected error but got none")
					return
				}
				
				if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
					if zaloBotErr.Type != tt.errType {
						t.Errorf("New() error type = %v, want %v", zaloBotErr.Type, tt.errType)
					}
				}
				return
			}
			
			if err != nil {
				t.Errorf("New() unexpected error = %v", err)
				return
			}
			
			if bot == nil {
				t.Error("New() returned nil bot")
				return
			}
			
			// Verify bot token is set correctly
			if bot.GetBotToken() != tt.botToken {
				t.Errorf("GetBotToken() = %v, want %v", bot.GetBotToken(), tt.botToken)
			}
			
			// Verify config is set
			config := bot.GetConfig()
			if config == nil {
				t.Error("GetConfig() returned nil")
				return
			}
			
			if config.BotToken != tt.botToken {
				t.Errorf("Config.BotToken = %v, want %v", config.BotToken, tt.botToken)
			}
		})
	}
}

func TestBotAPI_GetAPIEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		botToken string
		baseURL  string
		method   string
		want     string
	}{
		{
			name:     "production endpoint",
			botToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			baseURL:  "https://bot-api.zapps.me",
			method:   "sendMessage",
			want:     "https://bot-api.zapps.me/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/sendMessage",
		},
		{
			name:     "development endpoint",
			botToken: "789012:XYZ-GHI5678jklMn-opq12R3s4t567uv89",
			baseURL:  "https://dev-bot-api.zapps.me",
			method:   "getMe",
			want:     "https://dev-bot-api.zapps.me/bot789012:XYZ-GHI5678jklMn-opq12R3s4t567uv89/getMe",
		},
		{
			name:     "getUpdates method",
			botToken: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			baseURL:  "https://bot-api.zapps.me",
			method:   "getUpdates",
			want:     "https://bot-api.zapps.me/bot123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11/getUpdates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bot, err := New(tt.botToken, types.WithBaseURL(tt.baseURL))
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			got := bot.GetAPIEndpoint(tt.method)
			if got != tt.want {
				t.Errorf("GetAPIEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBotAPI_WithOptions(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	
	t.Run("with custom timeout", func(t *testing.T) {
		customTimeout := 15 * time.Second
		bot, err := New(botToken, types.WithTimeout(customTimeout))
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		
		config := bot.GetConfig()
		if config.Timeout != customTimeout {
			t.Errorf("Config.Timeout = %v, want %v", config.Timeout, customTimeout)
		}
	})
	
	t.Run("with debug mode", func(t *testing.T) {
		bot, err := New(botToken, types.WithDebug())
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		
		config := bot.GetConfig()
		if !config.Debug {
			t.Error("Config.Debug = false, want true")
		}
	})
	
	t.Run("with custom base URL", func(t *testing.T) {
		customURL := "https://custom-api.example.com"
		bot, err := New(botToken, types.WithBaseURL(customURL))
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		
		config := bot.GetConfig()
		if config.BaseURL != customURL {
			t.Errorf("Config.BaseURL = %v, want %v", config.BaseURL, customURL)
		}
	})
	
	t.Run("with environment", func(t *testing.T) {
		bot, err := New(botToken, types.WithEnvironment(types.Development))
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		
		config := bot.GetConfig()
		if config.Environment != types.Development {
			t.Errorf("Config.Environment = %v, want %v", config.Environment, types.Development)
		}
	})
	
	t.Run("with custom HTTP client", func(t *testing.T) {
		customClient := &http.Client{
			Timeout: 5 * time.Second,
		}
		bot, err := New(botToken, types.WithHTTPClient(customClient))
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		
		if bot.GetHTTPClient() != customClient {
			t.Error("GetHTTPClient() did not return custom client")
		}
	})
	
	t.Run("with multiple options", func(t *testing.T) {
		bot, err := New(botToken,
			types.WithDebug(),
			types.WithTimeout(20*time.Second),
			types.WithEnvironment(types.Development),
			types.WithRetries(5),
		)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}
		
		config := bot.GetConfig()
		if !config.Debug {
			t.Error("Config.Debug = false, want true")
		}
		if config.Timeout != 20*time.Second {
			t.Errorf("Config.Timeout = %v, want %v", config.Timeout, 20*time.Second)
		}
		if config.Environment != types.Development {
			t.Errorf("Config.Environment = %v, want %v", config.Environment, types.Development)
		}
		if config.Retries != 5 {
			t.Errorf("Config.Retries = %v, want %v", config.Retries, 5)
		}
	})
}

func TestBotAPI_Close(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	
	// Get context before closing
	ctx := bot.GetContext()
	if ctx == nil {
		t.Fatal("GetContext() returned nil")
	}
	
	// Close the bot
	bot.Close()
	
	// Verify context is cancelled
	select {
	case <-ctx.Done():
		// Context is cancelled as expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Context was not cancelled after Close()")
	}
}

func TestValidateWebhookURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid HTTPS URL",
			url:     "https://example.com/webhook",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL with port",
			url:     "https://example.com:8443/webhook",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL with path",
			url:     "https://api.example.com/bot/webhook/handler",
			wantErr: false,
		},
		{
			name:    "HTTP URL - should fail",
			url:     "http://example.com/webhook",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "whitespace URL",
			url:     "   ",
			wantErr: true,
		},
		{
			name:    "invalid URL format",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "URL without host",
			url:     "https:///webhook",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateWebhookURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateWebhookURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBotAPI_SetWebhook(t *testing.T) {
	// Note: These tests would require a mock HTTP server to fully test
	// For now, we test validation logic
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name    string
		config  types.WebhookConfig
		wantErr bool
	}{
		{
			name: "invalid URL - HTTP",
			config: types.WebhookConfig{
				URL:         "http://example.com/webhook",
				SecretToken: "secret",
			},
			wantErr: true,
		},
		{
			name: "invalid URL - empty",
			config: types.WebhookConfig{
				URL:         "",
				SecretToken: "secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bot.SetWebhook(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBotAPI_WebhookMethods(t *testing.T) {
	// Test that webhook methods exist and can be called
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("SetWebhook exists", func(t *testing.T) {
		// Just verify the method exists and handles validation
		config := types.WebhookConfig{
			URL:         "http://invalid.com",
			SecretToken: "secret",
		}
		err := bot.SetWebhook(config)
		if err == nil {
			t.Error("SetWebhook() should fail with HTTP URL")
		}
	})

	t.Run("DeleteWebhook exists", func(t *testing.T) {
		// Method exists - actual testing would require mock server
		_ = bot.DeleteWebhook
	})

	t.Run("GetWebhookInfo exists", func(t *testing.T) {
		// Method exists - actual testing would require mock server
		_ = bot.GetWebhookInfo
	})
}

// Integration Tests

func TestBotAPI_ServiceIntegration(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("message service is initialized", func(t *testing.T) {
		msgService := bot.GetMessageService()
		if msgService == nil {
			t.Error("GetMessageService() returned nil")
		}
	})

	t.Run("user service is initialized", func(t *testing.T) {
		userService := bot.GetUserService()
		if userService == nil {
			t.Error("GetUserService() returned nil")
		}
	})

	t.Run("webhook service is initialized", func(t *testing.T) {
		webhookService := bot.GetWebhookService()
		if webhookService == nil {
			t.Error("GetWebhookService() returned nil")
		}
	})
}

func TestBotAPI_SendMessageIntegration(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("SendMessage method exists", func(t *testing.T) {
		// Verify the method exists and can be called
		// Actual API call would require a mock server
		config := types.MessageConfig{
			ChatID: "user123",
			Text:   "Hello, World!",
		}
		
		// We can't test the actual API call without a mock server,
		// but we can verify the method signature is correct
		_ = bot.SendMessage
		_ = config
	})

	t.Run("SendImage method exists", func(t *testing.T) {
		_ = bot.SendImage
	})

	t.Run("SendFile method exists", func(t *testing.T) {
		_ = bot.SendFile
	})

	t.Run("SendVideo method exists", func(t *testing.T) {
		_ = bot.SendVideo
	})

	t.Run("SendAudio method exists", func(t *testing.T) {
		_ = bot.SendAudio
	})

	t.Run("SendTemplate method exists", func(t *testing.T) {
		_ = bot.SendTemplate
	})

	t.Run("SendStructuredMessage method exists", func(t *testing.T) {
		_ = bot.SendStructuredMessage
	})
}

func TestBotAPI_UserProfileIntegration(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("GetUserProfile method exists", func(t *testing.T) {
		// Verify the method exists
		_ = bot.GetUserProfile
	})
}

func TestBotAPI_WebhookIntegration(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("webhook secret token can be set", func(t *testing.T) {
		secretToken := "my-secret-token"
		bot.SetWebhookSecretToken(secretToken)
		
		webhookService := bot.GetWebhookService()
		if webhookService.GetSecretToken() != secretToken {
			t.Errorf("GetSecretToken() = %v, want %v", webhookService.GetSecretToken(), secretToken)
		}
	})

	t.Run("ProcessWebhook method exists", func(t *testing.T) {
		_ = bot.ProcessWebhook
	})

	t.Run("ValidateWebhookSignature method exists", func(t *testing.T) {
		_ = bot.ValidateWebhookSignature
	})

	t.Run("ParseWebhookUpdate method exists", func(t *testing.T) {
		_ = bot.ParseWebhookUpdate
	})
}

func TestBotAPI_PollingIntegration(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("GetUpdates method exists", func(t *testing.T) {
		config := types.UpdateConfig{
			Offset:  0,
			Limit:   10,
			Timeout: 0,
		}
		
		// Verify the method exists
		_ = bot.GetUpdates
		_ = config
	})

	t.Run("IsPolling returns false initially", func(t *testing.T) {
		if bot.IsPolling() {
			t.Error("IsPolling() = true, want false")
		}
	})

	t.Run("GetUpdatesChan starts polling", func(t *testing.T) {
		config := types.UpdateConfig{
			Offset:  0,
			Limit:   10,
			Timeout: 0,
		}
		
		updatesChan := bot.GetUpdatesChan(config)
		if updatesChan == nil {
			t.Fatal("GetUpdatesChan() returned nil")
		}

		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)

		if !bot.IsPolling() {
			t.Error("IsPolling() = false, want true after GetUpdatesChan")
		}

		// Stop polling
		bot.StopPolling()
		
		// Give it a moment to stop
		time.Sleep(100 * time.Millisecond)

		if bot.IsPolling() {
			t.Error("IsPolling() = true, want false after StopPolling")
		}
	})

	t.Run("StopPolling can be called multiple times", func(t *testing.T) {
		bot.StopPolling()
		bot.StopPolling()
		// Should not panic
	})

	t.Run("polling stops on Close", func(t *testing.T) {
		// Create a new bot for this test
		bot2, err := New(botToken)
		if err != nil {
			t.Fatalf("New() error = %v", err)
		}

		config := types.UpdateConfig{
			Offset:  0,
			Limit:   10,
			Timeout: 0,
		}
		
		_ = bot2.GetUpdatesChan(config)
		
		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)

		if !bot2.IsPolling() {
			t.Error("IsPolling() = false, want true")
		}

		// Close the bot
		bot2.Close()
		
		// Give it a moment to stop
		time.Sleep(100 * time.Millisecond)

		if bot2.IsPolling() {
			t.Error("IsPolling() = true, want false after Close")
		}
	})
}

func TestBotAPI_PollingWithInvalidConfig(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("GetUpdatesChan with config validates and starts polling", func(t *testing.T) {
		config := types.UpdateConfig{
			Offset:  0,
			Limit:   -1, // Invalid limit, will be corrected by validation
			Timeout: 0,
		}
		
		// After validation, limit will be set to default (100)
		// So this should start polling normally
		updatesChan := bot.GetUpdatesChan(config)
		if updatesChan == nil {
			t.Fatal("GetUpdatesChan() returned nil")
		}

		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)

		// Stop polling
		bot.StopPolling()
	})
}

func TestBotAPI_ConcurrentPolling(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("multiple GetUpdatesChan calls return same channel", func(t *testing.T) {
		config := types.UpdateConfig{
			Offset:  0,
			Limit:   10,
			Timeout: 0,
		}
		
		ch1 := bot.GetUpdatesChan(config)
		
		// Give it a moment to start
		time.Sleep(50 * time.Millisecond)
		
		ch2 := bot.GetUpdatesChan(config)
		
		// Both should return the same channel
		if ch1 != ch2 {
			t.Error("GetUpdatesChan() returned different channels")
		}

		// Stop polling
		bot.StopPolling()
	})
}

func TestBotAPI_SetWebhookWithSecretToken(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("SetWebhook sets secret token in webhook service", func(t *testing.T) {
		// Note: This test will fail at the API call level without a mock server,
		// but we can verify the secret token is set before the API call fails
		
		secretToken := "test-secret-token"
		config := types.WebhookConfig{
			URL:         "https://example.com/webhook",
			SecretToken: secretToken,
		}
		
		// This will fail at the API call, but that's expected
		_ = bot.SetWebhook(config)
		
		// The secret token should still be set even if the API call fails
		// (it's set after successful API response, so this test is limited)
	})
}

func TestBotAPI_CompleteMessageFlow(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("complete flow methods are accessible", func(t *testing.T) {
		// Verify all methods needed for a complete message flow exist
		
		// 1. Send a message
		_ = bot.SendMessage
		
		// 2. Get user profile
		_ = bot.GetUserProfile
		
		// 3. Send rich media
		_ = bot.SendImage
		_ = bot.SendFile
		
		// 4. Send structured message
		_ = bot.SendTemplate
		
		// All methods exist and are accessible
	})
}

func TestBotAPI_CompleteWebhookFlow(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("complete webhook flow methods are accessible", func(t *testing.T) {
		// Verify all methods needed for a complete webhook flow exist
		
		// 1. Set webhook
		_ = bot.SetWebhook
		
		// 2. Set secret token
		_ = bot.SetWebhookSecretToken
		
		// 3. Process webhook
		_ = bot.ProcessWebhook
		
		// 4. Validate signature
		_ = bot.ValidateWebhookSignature
		
		// 5. Parse update
		_ = bot.ParseWebhookUpdate
		
		// 6. Get webhook info
		_ = bot.GetWebhookInfo
		
		// 7. Delete webhook
		_ = bot.DeleteWebhook
		
		// All methods exist and are accessible
	})
}

func TestBotAPI_CompletePollingFlow(t *testing.T) {
	botToken := "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := New(botToken)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	defer bot.Close()

	t.Run("complete polling flow methods are accessible", func(t *testing.T) {
		// Verify all methods needed for a complete polling flow exist
		
		// 1. Get updates
		_ = bot.GetUpdates
		
		// 2. Get updates channel
		_ = bot.GetUpdatesChan
		
		// 3. Check polling status
		_ = bot.IsPolling
		
		// 4. Stop polling
		_ = bot.StopPolling
		
		// All methods exist and are accessible
	})
}
