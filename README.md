# Go Zalo Bot SDK

A comprehensive Go SDK for building Zalo Bot applications with support for bot token authentication, message handling, webhooks, and rich media content.

[![Go Reference](https://pkg.go.dev/badge/github.com/vkhangstack/go-zalo-bot.svg)](https://pkg.go.dev/github.com/vkhangstack/go-zalo-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/vkhangstack/go-zalo-bot)](https://goreportcard.com/report/github.com/vkhangstack/go-zalo-bot)
[![CI](https://github.com/vkhangstack/go-zalo-bot/actions/workflows/ci.yml/badge.svg)](https://github.com/vkhangstack/go-zalo-bot/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Features

- 🔐 **Bot Token Authentication** - Simple authentication using bot tokens (similar to Telegram Bot API)
- 💬 **Message Handling** - Send and receive text messages with Unicode support (including Vietnamese)
- 📸 **Rich Media** - Support for images, files, videos, audio, and structured messages
- 🎯 **Interactive Elements** - Buttons, quick replies, and postback events
- 🔔 **Webhooks** - Real-time event processing with signature validation
- 📊 **Polling** - Alternative update mechanism with long polling support
- 👤 **User Profiles** - Retrieve and manage user profile information
- 🔄 **Retry Mechanisms** - Automatic retry with exponential backoff for transient failures
- ⚠️ **Error Handling** - Comprehensive error types and handling strategies
- 📝 **Logging** - Configurable debug logging for troubleshooting
- 🧪 **Well Tested** - Extensive unit and integration tests

## Installation

### Requirements

- Go 1.20 or higher
- Git

### Install SDK

```bash
go get github.com/vkhangstack/go-zalo-bot
```

### Development Tools (Optional)

For contributors and developers:

```bash
# Install golangci-lint for code quality checks
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Windows
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Quick Start

### Basic Bot with Polling

```go
package main

import (
    "log"
    
    zalobot "github.com/vkhangstack/go-zalo-bot"
    "github.com/vkhangstack/go-zalo-bot/types"
)

func main() {
    // Create bot with bot token
    bot, err := zalobot.New("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
    if err != nil {
        log.Fatal(err)
    }
    defer bot.Close()

    // Configure polling
    updateConfig := types.UpdateConfig{
        Limit:   100,
        Timeout: 30, // Long polling timeout
    }

    // Get updates channel
    updates := bot.GetUpdatesChan(updateConfig)

    // Process updates
    for update := range updates {
        if update.Message != nil {
            // Echo the message
            bot.SendMessage(types.MessageConfig{
                ChatID: update.Message.From.ID,
                Text:   "You said: " + update.Message.Text,
            })
        }
    }
}
```

### Basic Bot with Webhooks

```go
package main

import (
    "encoding/json"
    "io"
    "log"
    "net/http"
    
    zalobot "github.com/vkhangstack/go-zalo-bot"
)

var bot *zalobot.BotAPI

func main() {
    var err error
    bot, err = zalobot.New("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
    if err != nil {
        log.Fatal(err)
    }
    defer bot.Close()

    // Set webhook secret for signature validation
    bot.SetWebhookSecretToken("your-webhook-secret")

    // Set up webhook handler
    http.HandleFunc("/webhook", webhookHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    signature := r.Header.Get("X-Zalo-Signature")

    // Process webhook with validation
    update, err := bot.ProcessWebhook(body, signature)
    if err != nil {
        http.Error(w, "Invalid signature", http.StatusUnauthorized)
        return
    }

    // Handle the update
    if update.Message != nil {
        bot.SendMessage(types.MessageConfig{
            ChatID: update.Message.From.ID,
            Text:   "You said: " + update.Message.Text,
        })
    }

    json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
```

## Configuration

### Bot Options

```go
bot, err := zalobot.New(botToken,
    zalobot.WithDebug(),                           // Enable debug logging
    zalobot.WithTimeout(30*time.Second),           // Set request timeout
    zalobot.WithRetries(5),                        // Set max retries
    zalobot.WithEnvironment(types.Production),     // Set environment
    zalobot.WithBaseURL("https://custom-api.com"), // Custom base URL
)
```

### Retry Configuration

```go
retryConfig := &types.RetryConfig{
    MaxRetries:    5,
    InitialDelay:  1 * time.Second,
    MaxDelay:      30 * time.Second,
    BackoffFactor: 2.0,
    RetryableErrors: []types.ErrorType{
        types.ErrorTypeNetwork,
        types.ErrorTypeRateLimit,
    },
}

bot, err := zalobot.New(botToken, zalobot.WithRetryConfig(retryConfig))
```

## Usage Examples

### Sending Messages

#### Text Message

```go
message, err := bot.SendMessage(types.MessageConfig{
    ChatID: "user123",
    Text:   "Hello, World! 👋",
})
```

#### Image Message

```go
message, err := bot.SendImage(types.ImageMessageConfig{
    ChatID:   "user123",
    ImageURL: "https://example.com/image.jpg",
    Caption:  "Check out this image!",
    MimeType: "image/jpeg",
})
```

#### File Message

```go
message, err := bot.SendFile(types.FileMessageConfig{
    ChatID:   "user123",
    FileURL:  "https://example.com/document.pdf",
    FileName: "document.pdf",
    MimeType: "application/pdf",
    Size:     1024 * 1024, // 1MB
})
```

#### Video Message

```go
message, err := bot.SendVideo("user123", "https://example.com/video.mp4", "video/mp4")
```

#### Audio Message

```go
message, err := bot.SendAudio("user123", "https://example.com/audio.mp3", "audio/mpeg")
```

### Structured Messages

#### Button Message

```go
message, err := bot.SendTemplate(types.StructuredMessageConfig{
    ChatID: "user123",
    StructuredMessage: types.StructuredMessage{
        Type: types.StructuredMessageTypeButton,
        Elements: []types.MessageElement{
            {
                Title:    "Choose an option",
                Subtitle: "Please select:",
                Buttons: []types.Button{
                    {
                        Type:    types.ButtonTypePostback,
                        Title:   "Option 1",
                        Payload: "OPTION_1",
                    },
                    {
                        Type:  types.ButtonTypeWebURL,
                        Title: "Visit Website",
                        URL:   "https://example.com",
                    },
                },
            },
        },
    },
})
```

#### Carousel Message

```go
message, err := bot.SendStructuredMessage(types.StructuredMessageConfig{
    ChatID: "user123",
    StructuredMessage: types.StructuredMessage{
        Type: types.StructuredMessageTypeTemplate,
        Elements: []types.MessageElement{
            {
                Title:    "Product 1",
                Subtitle: "Description of product 1",
                ImageURL: "https://example.com/product1.jpg",
                Buttons: []types.Button{
                    {Type: types.ButtonTypePostback, Title: "Buy", Payload: "BUY_1"},
                },
            },
            {
                Title:    "Product 2",
                Subtitle: "Description of product 2",
                ImageURL: "https://example.com/product2.jpg",
                Buttons: []types.Button{
                    {Type: types.ButtonTypePostback, Title: "Buy", Payload: "BUY_2"},
                },
            },
        },
    },
})
```

#### Quick Replies

```go
message, err := bot.SendStructuredMessage(types.StructuredMessageConfig{
    ChatID: "user123",
    StructuredMessage: types.StructuredMessage{
        Type: types.StructuredMessageTypeTemplate,
        QuickReplies: []types.QuickReply{
            {
                ContentType: types.QuickReplyTypeText,
                Title:       "Yes",
                Payload:     "QUICK_YES",
            },
            {
                ContentType: types.QuickReplyTypeText,
                Title:       "No",
                Payload:     "QUICK_NO",
            },
        },
    },
})
```

### User Profile

```go
profile, err := bot.GetUserProfile("user123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Name: %s\n", profile.Name)
fmt.Printf("ID: %s\n", profile.ID)
fmt.Printf("Avatar: %s\n", profile.Avatar)
```

### Webhook Management

#### Set Webhook

```go
err := bot.SetWebhook(types.WebhookConfig{
    URL:         "https://your-domain.com/webhook",
    SecretToken: "your-webhook-secret",
})
```

#### Delete Webhook

```go
err := bot.DeleteWebhook()
```

#### Get Webhook Info

```go
info, err := bot.GetWebhookInfo()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("URL: %s\n", info.URL)
fmt.Printf("Has Custom Certificate: %v\n", info.HasCustomCertificate)
```

## Error Handling

The SDK provides typed errors for better error handling:

```go
message, err := bot.SendMessage(config)
if err != nil {
    if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
        switch zaloBotErr.Type {
        case types.ErrorTypeAPI:
            log.Printf("API Error: %s (Code: %d)", zaloBotErr.Message, zaloBotErr.Code)
        case types.ErrorTypeNetwork:
            log.Printf("Network Error: %s", zaloBotErr.Message)
        case types.ErrorTypeValidation:
            log.Printf("Validation Error: %s", zaloBotErr.Message)
        case types.ErrorTypeRateLimit:
            log.Printf("Rate Limit Error: %s", zaloBotErr.Message)
        case types.ErrorTypeBotToken:
            log.Printf("Bot Token Error: %s", zaloBotErr.Message)
        }
        
        // Check if error is retryable
        if zaloBotErr.IsRetryable() {
            // Implement custom retry logic
        }
    }
}
```

## Authentication

The SDK uses bot token authentication similar to Telegram Bot API:

1. **Token Format**: `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`
2. **URL Pattern**: The bot token is embedded in the API endpoint URL:
   ```
   https://bot-api.zapps.me/bot${BOT_TOKEN}/method
   ```
3. **No Additional Headers**: No Authorization headers required
4. **Long-lived Tokens**: Bot tokens don't expire automatically

### Environment Configuration

```go
// Production (default)
bot, err := zalobot.New(botToken, zalobot.WithEnvironment(types.Production))

// Development
bot, err := zalobot.New(botToken, 
    zalobot.WithEnvironment(types.Development),
    zalobot.WithBaseURL("https://dev-bot-api.zapps.me"),
)
```

## Examples

The `examples/` directory contains complete working examples:

- **[Polling Example](examples/polling/)** - Bot using polling for updates
- **[Webhook Example](examples/webhook/)** - Bot using webhooks for real-time updates
- **[Advanced Examples](examples/advanced/)** - Rich media, user profiles, error handling

## API Reference

### Core Types

- `BotAPI` - Main bot client
- `MessageConfig` - Configuration for sending text messages
- `ImageMessageConfig` - Configuration for sending images
- `FileMessageConfig` - Configuration for sending files
- `StructuredMessageConfig` - Configuration for structured messages
- `UpdateConfig` - Configuration for polling updates
- `WebhookConfig` - Configuration for webhook setup
- `UserProfile` - User profile information
- `Update` - Incoming update from Zalo
- `Message` - Message structure
- `ZaloBotError` - Typed error structure

### Main Methods

#### Bot Management
- `New(botToken string, options ...BotOption) (*BotAPI, error)` - Create new bot instance
- `Close()` - Close bot and release resources
- `GetBotToken() string` - Get bot token
- `GetAPIEndpoint(method string) string` - Get API endpoint URL

#### Messaging
- `SendMessage(config MessageConfig) (*Message, error)` - Send text message
- `SendImage(config ImageMessageConfig) (*Message, error)` - Send image
- `SendFile(config FileMessageConfig) (*Message, error)` - Send file
- `SendVideo(chatID, videoURL, mimeType string) (*Message, error)` - Send video
- `SendAudio(chatID, audioURL, mimeType string) (*Message, error)` - Send audio
- `SendTemplate(config StructuredMessageConfig) (*Message, error)` - Send structured message
- `SendStructuredMessage(config StructuredMessageConfig) (*Message, error)` - Send structured message (alias)

#### Updates
- `GetUpdates(config UpdateConfig) ([]Update, error)` - Get updates via polling
- `GetUpdatesChan(config UpdateConfig) <-chan Update` - Get updates channel for polling
- `IsPolling() bool` - Check if bot is currently polling
- `StopPolling()` - Stop polling loop

#### Webhooks
- `SetWebhook(config WebhookConfig) error` - Set webhook URL
- `DeleteWebhook() error` - Delete webhook
- `GetWebhookInfo() (*WebhookInfo, error)` - Get webhook information
- `ProcessWebhook(payload []byte, signature string) (*Update, error)` - Process webhook with validation
- `ValidateWebhookSignature(payload []byte, signature string) error` - Validate webhook signature
- `SetWebhookSecretToken(token string)` - Set webhook secret token

#### User Management
- `GetUserProfile(userID string) (*UserProfile, error)` - Get user profile

## Development

### Automation Scripts

The project includes automation scripts for development and release workflows:

```bash
# Build validation
./scripts/build.sh

# Run tests with coverage
./scripts/test.sh --coverage

# Run linting and static analysis
./scripts/lint.sh

# Pre-release validation
./scripts/validate.sh

# Create a release
./scripts/release.sh v1.2.3
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed documentation on using these scripts.

### Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./types
go test ./services
go test ./auth

# Using the test script (recommended)
./scripts/test.sh --coverage --race
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and test locally:
   ```bash
   ./scripts/build.sh
   ./scripts/test.sh --coverage
   ./scripts/lint.sh --fix
   ./scripts/validate.sh
   ```
4. Commit your changes (`git commit -m 'Add some amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on:
- Development setup and tools
- Using automation scripts
- CI/CD workflows
- Troubleshooting common issues
- Creating releases

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://pkg.go.dev/github.com/vkhangstack/go-zalo-bot)
- 💬 [Issues](https://github.com/vkhangstack/go-zalo-bot/issues)
- 📧 Email: support@example.com

## Acknowledgments

- Inspired by the [Telegram Bot API](https://core.telegram.org/bots/api)
- Built with ❤️ for the Go community

## Roadmap

- [ ] Add support for more message types
- [ ] Implement message editing and deletion
- [ ] Add support for group chats
- [ ] Implement file upload from local filesystem
- [ ] Add middleware support for message processing
- [ ] Create CLI tool for bot management
- [ ] Add more examples and tutorials

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes.
