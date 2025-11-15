# Polling Example

This example demonstrates how to use the Zalo Bot SDK with polling to receive and handle updates.

## Features

- Bot initialization with bot token authentication
- Long polling for updates
- Message handling and response generation
- Command processing (/start, /help, /echo, /profile)
- User profile retrieval
- Postback event handling
- User action handling (join, leave, block)
- Graceful shutdown

## Prerequisites

- Go 1.16 or higher
- A Zalo Bot token (format: `123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11`)

## Setup

1. Set your bot token as an environment variable:

```bash
export ZALO_BOT_TOKEN="your-bot-token-here"
```

2. Run the example:

```bash
go run main.go
```

## How It Works

1. **Bot Initialization**: The bot is created using `zalobot.New()` with the bot token
2. **Polling Configuration**: Updates are fetched using long polling with a 30-second timeout
3. **Update Processing**: Incoming updates are processed in a loop
4. **Message Handling**: Different message types are handled appropriately
5. **Response Generation**: The bot responds to user messages based on content

## Supported Commands

- `/start` - Display welcome message and available commands
- `/help` - Show help information
- `/profile` - Retrieve and display user profile
- Any other text - Echo the message back

## Code Structure

- `main()` - Initializes the bot and starts polling
- `handleUpdate()` - Routes updates to appropriate handlers
- `handleTextMessage()` - Processes text messages and generates responses
- `handlePostback()` - Handles button click events
- `handleUserAction()` - Processes user actions (join, leave, block)

## Graceful Shutdown

The example includes graceful shutdown handling. Press `Ctrl+C` to stop the bot cleanly.

## Error Handling

The example includes basic error handling for:
- Bot initialization failures
- Message sending failures
- User profile retrieval failures

## Next Steps

- Explore the webhook example for an alternative update mechanism
- Check out the advanced examples for rich media messages
- Review the SDK documentation for more features
