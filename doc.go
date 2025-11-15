// Package zalobot provides a comprehensive Go SDK for building Zalo Bot applications.
//
// # Overview
//
// The Go Zalo Bot SDK enables developers to build chatbots and automated messaging
// solutions using the Zalo platform. It provides a simple, idiomatic Go interface
// for interacting with the Zalo Bot API.
//
// # Features
//
//   - Bot Token Authentication - Simple authentication using bot tokens
//   - Message Handling - Send and receive text messages with Unicode support
//   - Rich Media - Support for images, files, videos, audio, and structured messages
//   - Interactive Elements - Buttons, quick replies, and postback events
//   - Webhooks - Real-time event processing with signature validation
//   - Polling - Alternative update mechanism with long polling support
//   - User Profiles - Retrieve and manage user profile information
//   - Retry Mechanisms - Automatic retry with exponential backoff
//   - Error Handling - Comprehensive error types and handling strategies
//   - Logging - Configurable debug logging for troubleshooting
//
// # Installation
//
// Install the SDK using go get:
//
//	go get github.com/vkhangstack/go-zalo-bot
//
// # Quick Start
//
// Create a bot and send a message:
//
//	package main
//
//	import (
//	    "log"
//	    zalobot "github.com/vkhangstack/go-zalo-bot"
//	    "github.com/vkhangstack/go-zalo-bot/types"
//	)
//
//	func main() {
//	    // Create bot with bot token
//	    bot, err := zalobot.New("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer bot.Close()
//
//	    // Send a message
//	    message, err := bot.SendMessage(types.MessageConfig{
//	        ChatID: "user123",
//	        Text:   "Hello, World!",
//	    })
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    log.Printf("Message sent: %s", message.MessageID)
//	}
//
// # Authentication
//
// The SDK uses bot token authentication similar to Telegram Bot API.
// Bot tokens are in the format: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
//
// The bot token is embedded directly in the API endpoint URL:
//
//	https://bot-api.zapps.me/bot${BOT_TOKEN}/method
//
// No additional Authorization headers are required.
//
// # Polling for Updates
//
// Use polling to receive updates from Zalo:
//
//	updateConfig := types.UpdateConfig{
//	    Limit:   100,
//	    Timeout: 30, // Long polling timeout in seconds
//	}
//
//	updates := bot.GetUpdatesChan(updateConfig)
//
//	for update := range updates {
//	    if update.Message != nil {
//	        // Handle text message
//	        bot.SendMessage(types.MessageConfig{
//	            ChatID: update.Message.From.ID,
//	            Text:   "You said: " + update.Message.Text,
//	        })
//	    }
//	}
//
// # Webhooks
//
// Set up webhooks for real-time event processing:
//
//	// Set webhook secret token for signature validation
//	bot.SetWebhookSecretToken("your-webhook-secret")
//
//	// Configure webhook
//	err := bot.SetWebhook(types.WebhookConfig{
//	    URL:         "https://your-domain.com/webhook",
//	    SecretToken: "your-webhook-secret",
//	})
//
// Process webhook requests in your HTTP handler:
//
//	func webhookHandler(w http.ResponseWriter, r *http.Request) {
//	    body, _ := io.ReadAll(r.Body)
//	    signature := r.Header.Get("X-Zalo-Signature")
//
//	    update, err := bot.ProcessWebhook(body, signature)
//	    if err != nil {
//	        http.Error(w, "Invalid signature", http.StatusUnauthorized)
//	        return
//	    }
//
//	    // Handle the update
//	    if update.Message != nil {
//	        bot.SendMessage(types.MessageConfig{
//	            ChatID: update.Message.From.ID,
//	            Text:   "Received: " + update.Message.Text,
//	        })
//	    }
//
//	    json.NewEncoder(w).Encode(map[string]bool{"ok": true})
//	}
//
// # Rich Media Messages
//
// Send images, files, videos, and audio:
//
//	// Send image
//	bot.SendImage(types.ImageMessageConfig{
//	    ChatID:   "user123",
//	    ImageURL: "https://example.com/image.jpg",
//	    Caption:  "Check this out!",
//	    MimeType: "image/jpeg",
//	})
//
//	// Send file
//	bot.SendFile(types.FileMessageConfig{
//	    ChatID:   "user123",
//	    FileURL:  "https://example.com/document.pdf",
//	    FileName: "document.pdf",
//	    MimeType: "application/pdf",
//	    Size:     1024 * 1024,
//	})
//
//	// Send video
//	bot.SendVideo("user123", "https://example.com/video.mp4", "video/mp4")
//
// # Structured Messages
//
// Send interactive messages with buttons and quick replies:
//
//	bot.SendTemplate(types.StructuredMessageConfig{
//	    ChatID: "user123",
//	    StructuredMessage: types.StructuredMessage{
//	        Type: types.StructuredMessageTypeButton,
//	        Elements: []types.MessageElement{
//	            {
//	                Title:    "Choose an option",
//	                Subtitle: "Please select:",
//	                Buttons: []types.Button{
//	                    {
//	                        Type:    types.ButtonTypePostback,
//	                        Title:   "Option 1",
//	                        Payload: "OPTION_1",
//	                    },
//	                    {
//	                        Type:  types.ButtonTypeWebURL,
//	                        Title: "Visit Website",
//	                        URL:   "https://example.com",
//	                    },
//	                },
//	            },
//	        },
//	    },
//	})
//
// # User Profiles
//
// Retrieve user profile information:
//
//	profile, err := bot.GetUserProfile("user123")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Name: %s\n", profile.Name)
//	fmt.Printf("ID: %s\n", profile.ID)
//	fmt.Printf("Avatar: %s\n", profile.Avatar)
//
// # Error Handling
//
// The SDK provides typed errors for better error handling:
//
//	message, err := bot.SendMessage(config)
//	if err != nil {
//	    if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
//	        switch zaloBotErr.Type {
//	        case types.ErrorTypeAPI:
//	            log.Printf("API Error: %s", zaloBotErr.Message)
//	        case types.ErrorTypeNetwork:
//	            log.Printf("Network Error: %s", zaloBotErr.Message)
//	        case types.ErrorTypeValidation:
//	            log.Printf("Validation Error: %s", zaloBotErr.Message)
//	        case types.ErrorTypeRateLimit:
//	            log.Printf("Rate Limit Error: %s", zaloBotErr.Message)
//	        }
//
//	        if zaloBotErr.IsRetryable() {
//	            // Implement custom retry logic
//	        }
//	    }
//	}
//
// # Configuration
//
// Customize bot behavior with options:
//
//	bot, err := zalobot.New(botToken,
//	    zalobot.WithDebug(),                           // Enable debug logging
//	    zalobot.WithTimeout(30*time.Second),           // Set request timeout
//	    zalobot.WithRetries(5),                        // Set max retries
//	    zalobot.WithEnvironment(types.Production),     // Set environment
//	    zalobot.WithBaseURL("https://custom-api.com"), // Custom base URL
//	)
//
// Configure retry behavior:
//
//	retryConfig := &types.RetryConfig{
//	    MaxRetries:    5,
//	    InitialDelay:  1 * time.Second,
//	    MaxDelay:      30 * time.Second,
//	    BackoffFactor: 2.0,
//	    RetryableErrors: []types.ErrorType{
//	        types.ErrorTypeNetwork,
//	        types.ErrorTypeRateLimit,
//	    },
//	}
//
//	bot, err := zalobot.New(botToken, zalobot.WithRetryConfig(retryConfig))
//
// # Examples
//
// The SDK includes comprehensive examples:
//
//   - Polling Example - Bot using polling for updates
//   - Webhook Example - Bot using webhooks for real-time updates
//   - Advanced Examples - Rich media, user profiles, error handling
//
// See the examples directory for complete working examples.
//
// # Package Structure
//
// The SDK is organized into several packages:
//
//   - zalobot - Main package with BotAPI client
//   - types - Type definitions for messages, users, configs, and errors
//   - services - Service implementations for messages, users, and webhooks
//   - auth - Authentication and token management
//   - utils - Utility functions and helpers
//
// # Best Practices
//
// 1. Always close the bot when done:
//
//	defer bot.Close()
//
// 2. Handle errors appropriately:
//
//	if err != nil {
//	    if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
//	        // Handle typed error
//	    }
//	}
//
// 3. Use context for timeout control:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	messageService.Send(ctx, config)
//
// 4. Implement rate limiting to avoid API limits:
//
//	time.Sleep(100 * time.Millisecond) // Between requests
//
// 5. Cache user profiles to reduce API calls:
//
//	// Implement a simple cache
//	cache := make(map[string]*types.UserProfile)
//
// # Support
//
// For more information, examples, and documentation:
//
//   - Documentation: https://pkg.go.dev/github.com/vkhangstack/go-zalo-bot
//   - Examples: https://github.com/vkhangstack/go-zalo-bot/tree/main/examples
//   - Issues: https://github.com/vkhangstack/go-zalo-bot/issues
//
// # License
//
// This project is licensed under the MIT License.
package zalobot
