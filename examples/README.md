# Go Zalo Bot SDK Examples

This directory contains comprehensive examples demonstrating how to use the Go Zalo Bot SDK.

## Overview

The examples are organized into three categories:

1. **Polling** - Bot using polling to receive updates
2. **Webhook** - Bot using webhooks for real-time updates
3. **Advanced** - Advanced features including rich media, user profiles, and error handling

## Quick Start

### Prerequisites

Before running any example, you need:

1. A Zalo Bot token (format: `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`)
2. Go 1.16 or higher installed

### Set Environment Variables

```bash
export ZALO_BOT_TOKEN="your-bot-token-here"
```

For webhook examples, also set:

```bash
export ZALO_WEBHOOK_SECRET="your-webhook-secret-here"
export PORT="8080"  # Optional, defaults to 8080
```

## Examples

### 1. Polling Example

**Location**: `polling/`

**Description**: Demonstrates how to build a bot that uses polling to receive updates from Zalo.

**Features**:
- Bot initialization with bot token
- Long polling for updates
- Message handling and response generation
- Command processing (/start, /help, /profile)
- User profile retrieval
- Graceful shutdown

**Run**:
```bash
cd polling
go run main.go
```

**Learn More**: See [polling/README.md](polling/README.md)

### 2. Webhook Example

**Location**: `webhook/`

**Description**: Demonstrates how to build a bot that uses webhooks for real-time event processing.

**Features**:
- HTTP server setup for webhook endpoint
- Webhook signature validation
- Real-time event processing
- Message handling and response generation
- Health check endpoint
- Graceful shutdown

**Run**:
```bash
cd webhook
go run main.go
```

**Learn More**: See [webhook/README.md](webhook/README.md)

### 3. Advanced Examples

**Location**: `advanced/`

**Description**: Demonstrates advanced SDK features including rich media, user profiles, and error handling.

#### 3.1 Rich Media Messages (`rich_media.go`)

**Features**:
- Sending images with captions
- Sending files with size validation
- Sending videos and audio
- Structured messages with buttons
- Carousel templates
- Quick replies
- Media sequences

**Run**:
```bash
cd advanced
go run rich_media.go
```

#### 3.2 User Profile Handling (`user_profile.go`)

**Features**:
- Retrieving user profiles
- Sending personalized messages
- Handling multiple users
- Building a user directory
- Profile caching for performance

**Run**:
```bash
cd advanced
go run user_profile.go
```

#### 3.3 Error Handling (`error_handling.go`)

**Features**:
- Basic error handling patterns
- Automatic retry mechanisms
- Timeout handling with context
- Rate limit handling with exponential backoff
- Validation error handling
- Debug logging
- Custom error handlers
- Fallback mechanisms

**Run**:
```bash
cd advanced
go run error_handling.go
```

**Learn More**: See [advanced/README.md](advanced/README.md)

## Example Comparison

| Feature | Polling | Webhook | Advanced |
|---------|---------|---------|----------|
| Update Mechanism | Long polling | Real-time webhooks | Both |
| Message Handling | ✓ | ✓ | ✓ |
| Rich Media | Basic | Basic | ✓ |
| User Profiles | ✓ | ✓ | ✓ |
| Error Handling | Basic | Basic | ✓ |
| Structured Messages | - | - | ✓ |
| HTTP Server | - | ✓ | - |
| Signature Validation | - | ✓ | - |

## Choosing Between Polling and Webhooks

### Use Polling When:
- You're developing locally without a public URL
- You want simpler setup without HTTP server
- You need to process updates at your own pace
- You're behind a firewall or NAT

### Use Webhooks When:
- You need real-time event processing
- You have a publicly accessible HTTPS endpoint
- You want to reduce server load
- You're deploying to production

## Common Patterns

### 1. Basic Message Echo Bot

```go
bot, _ := zalobot.New(botToken)
defer bot.Close()

updates := bot.GetUpdatesChan(types.UpdateConfig{Limit: 100, Timeout: 30})

for update := range updates {
    if update.Message != nil {
        bot.SendMessage(types.MessageConfig{
            ChatID: update.Message.From.ID,
            Text:   "You said: " + update.Message.Text,
        })
    }
}
```

### 2. Command Handler

```go
func handleCommand(bot *zalobot.BotAPI, message *types.Message) {
    switch message.Text {
    case "/start":
        bot.SendMessage(types.MessageConfig{
            ChatID: message.From.ID,
            Text:   "Welcome! Use /help for commands.",
        })
    case "/help":
        bot.SendMessage(types.MessageConfig{
            ChatID: message.From.ID,
            Text:   "Available commands: /start, /help, /profile",
        })
    case "/profile":
        profile, _ := bot.GetUserProfile(message.From.ID)
        bot.SendMessage(types.MessageConfig{
            ChatID: message.From.ID,
            Text:   fmt.Sprintf("Your name: %s", profile.Name),
        })
    }
}
```

### 3. Error Handling

```go
message, err := bot.SendMessage(config)
if err != nil {
    if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
        switch zaloBotErr.Type {
        case types.ErrorTypeRateLimit:
            // Wait and retry
            time.Sleep(2 * time.Second)
            bot.SendMessage(config)
        case types.ErrorTypeValidation:
            // Fix input and retry
            log.Printf("Validation error: %s", zaloBotErr.Message)
        default:
            log.Printf("Error: %s", zaloBotErr.Message)
        }
    }
}
```

### 4. Rich Media

```go
// Send image
bot.SendImage(types.ImageMessageConfig{
    ChatID:   "user123",
    ImageURL: "https://example.com/image.jpg",
    Caption:  "Check this out!",
    MimeType: "image/jpeg",
})

// Send structured message with buttons
bot.SendTemplate(types.StructuredMessageConfig{
    ChatID: "user123",
    StructuredMessage: types.StructuredMessage{
        Type: types.StructuredMessageTypeButton,
        Elements: []types.MessageElement{
            {
                Title: "Choose an option",
                Buttons: []types.Button{
                    {Type: types.ButtonTypePostback, Title: "Option 1", Payload: "OPT1"},
                    {Type: types.ButtonTypeWebURL, Title: "Visit", URL: "https://example.com"},
                },
            },
        },
    },
})
```

## Development Tips

### 1. Testing Locally

Use ngrok to expose your local webhook endpoint:

```bash
# Start your bot
go run webhook/main.go

# In another terminal
ngrok http 8080

# Use the ngrok HTTPS URL as your webhook URL
```

### 2. Debug Mode

Enable debug mode to see detailed logs:

```go
bot, err := zalobot.New(botToken, zalobot.WithDebug())
```

### 3. Custom Timeout

Set custom timeout for requests:

```go
bot, err := zalobot.New(botToken, zalobot.WithTimeout(30*time.Second))
```

### 4. Rate Limiting

Add delays between requests to avoid rate limiting:

```go
for _, message := range messages {
    bot.SendMessage(message)
    time.Sleep(100 * time.Millisecond)
}
```

## Troubleshooting

### Common Issues

1. **"Bot token is required" error**
   - Make sure `ZALO_BOT_TOKEN` environment variable is set
   - Verify the token format is correct

2. **"Invalid signature" error (webhooks)**
   - Ensure `ZALO_WEBHOOK_SECRET` matches the secret in webhook config
   - Check that the signature header is being sent correctly

3. **Timeout errors**
   - Increase timeout with `WithTimeout()` option
   - Check network connectivity

4. **Rate limiting errors**
   - Add delays between requests
   - Implement exponential backoff (see error_handling.go)

5. **File size errors**
   - Ensure files are under 50MB limit
   - Check file MIME type is supported

## Next Steps

1. Review the [main README](../README.md) for SDK overview
2. Check the [API documentation](https://pkg.go.dev/github.com/vkhangstack/go-zalo-bot)
3. Read the [contributing guide](../CONTRIBUTING.md) to contribute
4. Explore the source code for implementation details

## Support

- **Documentation**: https://pkg.go.dev/github.com/vkhangstack/go-zalo-bot
- **Issues**: https://github.com/vkhangstack/go-zalo-bot/issues
- **Examples**: https://github.com/vkhangstack/go-zalo-bot/tree/main/examples

## License

All examples are licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.
