# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-XX

### Added

#### Core Features
- Bot token authentication system (similar to Telegram Bot API)
- Main `BotAPI` client with comprehensive bot management
- Support for both polling and webhook update mechanisms
- Graceful shutdown and resource cleanup

#### Messaging
- Text message sending with Unicode support (including Vietnamese characters)
- Rich media message support:
  - Image messages with MIME type validation
  - File messages with size limits (max 50MB)
  - Video messages
  - Audio messages
- Structured messages with interactive elements:
  - Button templates
  - Carousel templates
  - Quick replies
- Message configuration validation

#### Webhooks
- Webhook setup and management (`SetWebhook`, `DeleteWebhook`, `GetWebhookInfo`)
- Webhook signature validation for security
- Support for multiple event types:
  - Text messages
  - Attachments
  - Postback events
  - User actions (join, leave, block)
- Real-time event processing

#### Polling
- Long polling support with configurable timeout
- Update channel for continuous polling
- Automatic offset management
- Polling state management (`IsPolling`, `StopPolling`)

#### User Management
- User profile retrieval
- Profile data including name, ID, and avatar
- Rate limiting handling for user service

#### Error Handling
- Typed error system with `ZaloBotError`
- Error types:
  - `ErrorTypeAPI` - API-related errors
  - `ErrorTypeNetwork` - Network connectivity issues
  - `ErrorTypeValidation` - Input validation failures
  - `ErrorTypeRateLimit` - Rate limiting errors
  - `ErrorTypeBotToken` - Bot token authentication errors
- Retryable error detection
- Descriptive error messages

#### Retry Mechanisms
- Automatic retry with exponential backoff
- Configurable retry behavior:
  - Max retries
  - Initial delay
  - Max delay
  - Backoff factor
  - Retryable error types
- Rate limit handling with exponential backoff

#### Configuration
- Flexible bot configuration with options pattern
- Configuration options:
  - `WithDebug()` - Enable debug logging
  - `WithTimeout()` - Set request timeout
  - `WithRetries()` - Set max retries
  - `WithEnvironment()` - Set environment (development/production)
  - `WithBaseURL()` - Custom base URL
  - `WithHTTPClient()` - Custom HTTP client
  - `WithRetryConfig()` - Custom retry configuration
- Environment support (development and production)

#### Validation
- Input validation for all message types
- Recipient ID validation
- Bot token format validation
- File size and MIME type validation
- Webhook URL validation (HTTPS required)
- User ID validation

#### Logging
- Optional debug logging
- Configurable log levels
- Structured logging support

#### Testing
- Comprehensive unit tests for all packages
- Integration tests for core functionality
- Test coverage for:
  - Authentication and token management
  - Message sending and receiving
  - Webhook processing
  - User profile retrieval
  - Error handling
  - Configuration validation

#### Documentation
- Complete README with quick start guide
- Comprehensive API documentation with GoDoc comments
- Usage examples:
  - Polling example application
  - Webhook example application
  - Advanced examples (rich media, user profiles, error handling)
- Package-level documentation
- Inline code documentation

#### Examples
- **Polling Example**: Complete bot using polling for updates
  - Message handling and response generation
  - Command processing
  - User profile retrieval
  - Graceful shutdown
- **Webhook Example**: Complete bot using webhooks
  - HTTP server setup
  - Signature validation
  - Real-time event processing
  - Health check endpoint
- **Advanced Examples**:
  - Rich media messages (images, files, videos, audio)
  - Structured messages (buttons, carousels, quick replies)
  - User profile handling and caching
  - Error handling patterns
  - Retry mechanisms
  - Logging and debugging

### Package Structure

```
go-zalo-bot/
├── client.go              # Main BotAPI client
├── doc.go                 # Package documentation
├── auth/
│   ├── token.go           # Bot token management
│   └── auth.go            # Authentication service
├── services/
│   ├── base.go            # Base service with common functionality
│   ├── message.go         # Message service
│   ├── user.go            # User service
│   └── webhook.go         # Webhook service
├── types/
│   ├── message.go         # Message types
│   ├── user.go            # User types
│   ├── webhook.go         # Webhook types
│   ├── config.go          # Configuration types
│   ├── error.go           # Error types
│   └── common.go          # Common types
├── utils/
│   ├── helpers.go         # Utility functions
│   └── logger.go          # Logging utilities
└── examples/
    ├── polling/           # Polling example
    ├── webhook/           # Webhook example
    └── advanced/          # Advanced examples
```

### Technical Details

#### Authentication
- Bot token format: `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`
- URL pattern: `https://bot-api.zapps.me/bot${BOT_TOKEN}/method`
- No Authorization headers required
- Long-lived tokens (don't expire automatically)

#### API Endpoints
- Base URL: `https://bot-api.zapps.me`
- Development URL: `https://dev-bot-api.zapps.me` (configurable)
- Methods: `sendMessage`, `sendImage`, `sendFile`, `sendTemplate`, `getUpdates`, `setWebhook`, `deleteWebhook`, `getWebhookInfo`, `getUserProfile`

#### Supported Media Types
- Images: `image/jpeg`, `image/png`, `image/gif`, `image/webp`
- Files: Any MIME type, max 50MB
- Videos: `video/mp4`, `video/mpeg`, `video/quicktime`, `video/webm`
- Audio: `audio/mpeg`, `audio/mp4`, `audio/ogg`, `audio/wav`

#### Requirements
- Go 1.16 or higher
- Valid Zalo Bot token
- HTTPS endpoint for webhooks (production)

### Dependencies
- Standard library only (no external dependencies)

### Notes
- This is the initial release of the Go Zalo Bot SDK
- The SDK follows the Telegram Bot API pattern for familiarity
- All public APIs are documented with GoDoc comments
- Comprehensive test coverage ensures reliability
- Examples demonstrate common use cases and best practices

## [Unreleased]

### Planned Features
- Message editing and deletion
- Group chat support
- File upload from local filesystem
- Middleware support for message processing
- CLI tool for bot management
- Additional message types
- Enhanced logging with structured output
- Metrics and monitoring support
- Rate limiting middleware
- Message queue support for high-volume bots

---

[1.0.0]: https://github.com/vkhangstack/go-zalo-bot/releases/tag/v1.0.0
