# Advanced Examples

This directory contains advanced usage examples for the Zalo Bot SDK, demonstrating rich media messages, user profile handling, error handling, retry mechanisms, and logging.

## Examples

### 1. Rich Media Messages (`rich_media.go`)

Demonstrates how to send various types of rich media content:

- **Image Messages**: Send images with captions and MIME type validation
- **File Messages**: Send files with size limits and format validation
- **Video Messages**: Send video content with proper MIME types
- **Audio Messages**: Send audio files
- **Structured Messages**: Send carousel-style messages with multiple elements
- **Button Messages**: Send messages with interactive buttons
- **Quick Reply Messages**: Send messages with quick reply options
- **Media Sequences**: Send multiple media items in sequence

**Key Features:**
- MIME type validation for different media types
- File size checking (max 50MB for files)
- Support for captions and metadata
- Interactive elements (buttons, quick replies)
- Postback handling for button clicks

### 2. User Profile Handling (`user_profile.go`)

Demonstrates how to work with user profile information:

- **Profile Retrieval**: Get user profile data including name, ID, and avatar
- **Personalized Messages**: Send customized messages using profile data
- **Multiple Users**: Handle profiles for multiple users efficiently
- **User Directory**: Build and manage a user directory
- **Profile Caching**: Implement caching to reduce API calls

**Key Features:**
- Efficient profile data retrieval
- Profile caching for performance optimization
- Batch user processing
- Personalization strategies

### 3. Error Handling and Retry Mechanisms (`error_handling.go`)

Demonstrates robust error handling and retry strategies:

- **Basic Error Handling**: Handle different error types appropriately
- **Automatic Retries**: Configure and use automatic retry mechanisms
- **Timeout Handling**: Implement request timeouts with context
- **Rate Limit Handling**: Handle API rate limits with exponential backoff
- **Validation Errors**: Catch and handle validation errors early
- **Debug Logging**: Enable detailed logging for troubleshooting

**Key Features:**
- Typed error handling with `ZaloBotError`
- Configurable retry strategies
- Exponential backoff for rate limiting
- Context-based timeout control
- Custom error handlers
- Fallback mechanisms
- Batch operation error collection

## Running the Examples

### Prerequisites

```bash
export ZALO_BOT_TOKEN="your-bot-token-here"
```

### Rich Media Example

```bash
go run rich_media.go
```

This will demonstrate sending various types of rich media messages. Make sure to replace `example_user_id` with an actual user ID.

### User Profile Example

```bash
go run user_profile.go
```

This will show how to retrieve and work with user profiles. Replace `example_user_id` with an actual user ID.

### Error Handling Example

```bash
go run error_handling.go
```

This will demonstrate various error handling patterns and retry mechanisms.

## Configuration Options

### Bot Initialization with Custom Settings

```go
bot, err := zalobot.New(botToken,
    zalobot.WithDebug(),                    // Enable debug logging
    zalobot.WithTimeout(30*time.Second),    // Set request timeout
    zalobot.WithRetries(5),                 // Set max retries
    zalobot.WithEnvironment(types.Production), // Set environment
    zalobot.WithRetryConfig(retryConfig),   // Custom retry config
)
```

### Custom Retry Configuration

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
```

## Error Types

The SDK provides typed errors for better error handling:

- `ErrorTypeAPI`: API-related errors from Zalo
- `ErrorTypeNetwork`: Network connectivity issues
- `ErrorTypeValidation`: Input validation failures
- `ErrorTypeRateLimit`: Rate limiting errors
- `ErrorTypeBotToken`: Bot token authentication errors

## Best Practices

### 1. Error Handling

Always check error types and handle them appropriately:

```go
if err != nil {
    if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
        switch zaloBotErr.Type {
        case types.ErrorTypeRateLimit:
            // Handle rate limiting
        case types.ErrorTypeValidation:
            // Handle validation errors
        default:
            // Handle other errors
        }
    }
}
```

### 2. Rate Limiting

Implement delays between requests to avoid rate limiting:

```go
for _, message := range messages {
    bot.SendMessage(message)
    time.Sleep(100 * time.Millisecond) // Small delay
}
```

### 3. Profile Caching

Cache user profiles to reduce API calls:

```go
cache := NewProfileCache(bot)
profile, err := cache.GetProfile(userID) // Uses cache if available
```

### 4. Timeout Control

Use context for timeout control:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

messageService.Send(ctx, config)
```

### 5. Validation

Validate input before making API calls:

```go
if err := config.Validate(); err != nil {
    // Handle validation error early
    return err
}
```

## Media Type Support

### Images
- MIME types: `image/jpeg`, `image/png`, `image/gif`, `image/webp`
- Recommended: Use JPEG for photos, PNG for graphics

### Files
- Maximum size: 50MB
- Common MIME types: `application/pdf`, `application/zip`, etc.

### Videos
- MIME types: `video/mp4`, `video/mpeg`, `video/quicktime`, `video/webm`
- Recommended: Use MP4 for best compatibility

### Audio
- MIME types: `audio/mpeg`, `audio/mp4`, `audio/ogg`, `audio/wav`
- Recommended: Use MP3 (audio/mpeg) for best compatibility

## Structured Message Patterns

### Button Template

```go
types.StructuredMessage{
    Type: types.StructuredMessageTypeButton,
    Elements: []types.MessageElement{
        {
            Title: "Choose an option",
            Buttons: []types.Button{
                {Type: types.ButtonTypePostback, Title: "Option 1", Payload: "OPT1"},
                {Type: types.ButtonTypeWebURL, Title: "Visit", URL: "https://..."},
            },
        },
    },
}
```

### Carousel Template

```go
types.StructuredMessage{
    Type: types.StructuredMessageTypeTemplate,
    Elements: []types.MessageElement{
        {Title: "Item 1", Subtitle: "Description", ImageURL: "...", Buttons: [...]},
        {Title: "Item 2", Subtitle: "Description", ImageURL: "...", Buttons: [...]},
    },
}
```

### Quick Replies

```go
types.StructuredMessage{
    QuickReplies: []types.QuickReply{
        {ContentType: types.QuickReplyTypeText, Title: "Yes", Payload: "YES"},
        {ContentType: types.QuickReplyTypeText, Title: "No", Payload: "NO"},
        {ContentType: types.QuickReplyTypeLocation, Title: "Share Location"},
    },
}
```

## Troubleshooting

### Common Issues

1. **Rate Limiting**: Add delays between requests or implement exponential backoff
2. **Timeout Errors**: Increase timeout duration or check network connectivity
3. **Validation Errors**: Verify input data before sending
4. **File Size Errors**: Ensure files are under 50MB limit
5. **MIME Type Errors**: Use supported MIME types for media

### Debug Mode

Enable debug mode to see detailed logs:

```go
bot, err := zalobot.New(botToken, zalobot.WithDebug())
```

This will log:
- Request URLs and payloads
- Response data
- Error details
- Retry attempts

## Next Steps

- Review the polling and webhook examples for update handling
- Check the main SDK documentation for complete API reference
- Explore the source code for implementation details
- Join the community for support and discussions
