# Webhook Example

This example demonstrates how to use the Zalo Bot SDK with webhooks to receive and handle updates in real-time.

## Features

- Bot initialization with bot token authentication
- HTTP server for webhook endpoint
- Webhook signature validation for security
- Real-time event processing
- Message handling and response generation
- Command processing (/start, /help, /profile)
- User profile retrieval
- Postback event handling
- User action handling (join, leave, block)
- Health check endpoint
- Graceful shutdown

## Prerequisites

- Go 1.16 or higher
- A Zalo Bot token (format: `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`)
- A webhook secret token for signature validation
- A publicly accessible HTTPS URL for the webhook endpoint (for production)

## Setup

1. Set your bot token and webhook secret as environment variables:

```bash
export ZALO_BOT_TOKEN="your-bot-token-here"
export ZALO_WEBHOOK_SECRET="your-webhook-secret-here"
export PORT="8080"  # Optional, defaults to 8080
```

2. Run the example:

```bash
go run main.go
```

3. Configure your webhook URL in Zalo Bot settings:

```go
// Use the SDK to set the webhook URL
bot.SetWebhook(types.WebhookConfig{
    URL:         "https://your-domain.com/webhook",
    SecretToken: "your-webhook-secret-here",
})
```

## How It Works

1. **Bot Initialization**: The bot is created using `zalobot.New()` with the bot token
2. **Webhook Secret**: The webhook secret token is set for signature validation
3. **HTTP Server**: An HTTP server is started to listen for webhook requests
4. **Signature Validation**: Each webhook request is validated using the signature
5. **Event Processing**: Valid updates are processed asynchronously
6. **Response Generation**: The bot responds to user messages based on content

## Endpoints

- `POST /webhook` - Webhook endpoint for receiving updates from Zalo
- `GET /health` - Health check endpoint

## Webhook Security

The example implements webhook signature validation to ensure requests are authentic:

1. Zalo sends a signature in the `X-Zalo-Signature` header
2. The SDK validates the signature using the webhook secret token
3. Invalid signatures are rejected with a 401 Unauthorized response

## Supported Commands

- `/start` - Display welcome message and available commands
- `/help` - Show help information
- `/profile` - Retrieve and display user profile
- Any other text - Echo the message back

## Code Structure

- `main()` - Initializes the bot and starts the HTTP server
- `webhookHandler()` - Handles incoming webhook requests and validates signatures
- `healthHandler()` - Provides a health check endpoint
- `handleUpdate()` - Routes updates to appropriate handlers
- `handleTextMessage()` - Processes text messages and generates responses
- `handlePostback()` - Handles button click events
- `handleUserAction()` - Processes user actions (join, leave, block)

## Development vs Production

### Development

For local development, you can use tools like [ngrok](https://ngrok.com/) to expose your local server:

```bash
# Start your bot
go run main.go

# In another terminal, start ngrok
ngrok http 8080

# Use the ngrok HTTPS URL as your webhook URL
```

### Production

For production deployment:

1. Deploy your application to a server with a valid SSL certificate
2. Use HTTPS for the webhook URL (required by Zalo)
3. Set appropriate environment variables
4. Consider using a reverse proxy (nginx, Caddy) for SSL termination
5. Implement proper logging and monitoring

## Error Handling

The example includes error handling for:
- Invalid HTTP methods
- Missing or invalid signatures
- Request body reading failures
- Message sending failures
- User profile retrieval failures

## Graceful Shutdown

The example includes graceful shutdown handling. Press `Ctrl+C` to stop the server cleanly.

## Testing

You can test the webhook endpoint locally:

```bash
# Test health endpoint
curl http://localhost:8080/health

# Test webhook endpoint (requires valid signature)
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-Zalo-Signature: your-signature" \
  -d '{"update_id": 1, "message": {"from": {"id": "123"}, "text": "Hello"}}'
```

## Next Steps

- Explore the polling example for an alternative update mechanism
- Check out the advanced examples for rich media messages
- Review the SDK documentation for more features
- Implement rate limiting and request throttling for production
- Add monitoring and alerting for webhook failures
