// Package main demonstrates error handling, retry mechanisms, and logging with the Zalo Bot SDK
// This example shows best practices for handling errors and implementing robust bot applications
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	zalobot "github.com/vkhangstack/go-zalo-bot"
	"github.com/vkhangstack/go-zalo-bot/types"
)

func main() {
	// Get bot token from environment variable
	botToken := os.Getenv("ZALO_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("ZALO_BOT_TOKEN environment variable is required")
	}

	// Create bot instance with custom retry configuration
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

	bot, err := zalobot.New(botToken,
		types.WithDebug(),
		types.WithTimeout(30*time.Second),
		types.WithRetries(5),
		types.WithRetryConfig(retryConfig),
	)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	defer bot.Close()

	// Example chat ID
	chatID := "example_user_id"

	// Demonstrate error handling patterns
	demonstrateBasicErrorHandling(bot, chatID)
	demonstrateRetryMechanism(bot, chatID)
	demonstrateTimeoutHandling(bot, chatID)
	demonstrateRateLimitHandling(bot, chatID)
	demonstrateValidationErrors(bot)
	demonstrateLogging(bot, chatID)
}

// demonstrateBasicErrorHandling shows basic error handling patterns
func demonstrateBasicErrorHandling(bot *zalobot.BotAPI, chatID string) {
	fmt.Println("\n=== Basic Error Handling ===")

	// Attempt to send a message
	_, err := bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   "Hello, World!",
	})

	if err != nil {
		// Check error type
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
			default:
				log.Printf("Unknown Error: %s", zaloBotErr.Message)
			}

			// Check if error is retryable
			if zaloBotErr.IsRetryable() {
				log.Println("This error is retryable")
			}
		} else {
			log.Printf("Unexpected error: %v", err)
		}
		return
	}

	log.Println("Message sent successfully")
	fmt.Println("============================\n")
}

// demonstrateRetryMechanism shows how the SDK handles retries automatically
func demonstrateRetryMechanism(bot *zalobot.BotAPI, chatID string) {
	fmt.Println("\n=== Retry Mechanism ===")

	// The SDK automatically retries on transient failures
	// This is configured via WithRetries() and WithRetryConfig()
	_, err := bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   "This message will be retried on transient failures",
	})

	if err != nil {
		log.Printf("Failed after retries: %v", err)
		return
	}

	log.Println("Message sent (with automatic retries if needed)")
	fmt.Println("=======================\n")
}

// demonstrateTimeoutHandling shows how to handle timeout scenarios
func demonstrateTimeoutHandling(bot *zalobot.BotAPI, chatID string) {
	fmt.Println("\n=== Timeout Handling ===")

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the message service directly with context
	messageService := bot.GetMessageService()
	_, err := messageService.Send(ctx, types.MessageConfig{
		ChatID: chatID,
		Text:   "This message has a 5-second timeout",
	})

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
		} else {
			log.Printf("Error: %v", err)
		}
		return
	}

	log.Println("Message sent within timeout")
	fmt.Println("========================\n")
}

// demonstrateRateLimitHandling shows how to handle rate limiting
func demonstrateRateLimitHandling(bot *zalobot.BotAPI, chatID string) {
	fmt.Println("\n=== Rate Limit Handling ===")

	// Send multiple messages in quick succession
	for i := 0; i < 5; i++ {
		_, err := bot.SendMessage(types.MessageConfig{
			ChatID: chatID,
			Text:   fmt.Sprintf("Message %d", i+1),
		})

		if err != nil {
			if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
				if zaloBotErr.Type == types.ErrorTypeRateLimit {
					log.Printf("Rate limited on message %d, waiting before retry...", i+1)
					// The SDK handles exponential backoff automatically
					// But you can also implement custom rate limiting here
					time.Sleep(2 * time.Second)
					continue
				}
			}
			log.Printf("Failed to send message %d: %v", i+1, err)
			continue
		}

		log.Printf("Message %d sent successfully", i+1)
		
		// Add a small delay between messages to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("===========================\n")
}

// demonstrateValidationErrors shows how to handle validation errors
func demonstrateValidationErrors(bot *zalobot.BotAPI) {
	fmt.Println("\n=== Validation Errors ===")

	// Try to send a message with invalid configuration
	_, err := bot.SendMessage(types.MessageConfig{
		ChatID: "", // Empty chat ID - validation error
		Text:   "This will fail validation",
	})

	if err != nil {
		if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
			if zaloBotErr.Type == types.ErrorTypeValidation {
				log.Printf("Validation failed: %s", zaloBotErr.Message)
			}
		}
	}

	// Try to send an image with invalid MIME type
	_, err = bot.SendImage(types.ImageMessageConfig{
		ChatID:   "user123",
		ImageURL: "https://example.com/image.jpg",
		MimeType: "invalid/mime", // Invalid MIME type
	})

	if err != nil {
		log.Printf("Image validation failed: %v", err)
	}

	fmt.Println("=========================\n")
}

// demonstrateLogging shows how to use debug logging
func demonstrateLogging(bot *zalobot.BotAPI, chatID string) {
	fmt.Println("\n=== Logging Demo ===")

	// Debug mode is enabled via WithDebug() option
	// The SDK will log detailed information about requests and responses

	log.Println("Sending message with debug logging enabled...")

	_, err := bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   "This message is sent with debug logging",
	})

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	log.Println("Check the logs above for detailed debug information")
	fmt.Println("====================\n")
}

// Example: Custom error handler wrapper
type ErrorHandler struct {
	bot *zalobot.BotAPI
}

func NewErrorHandler(bot *zalobot.BotAPI) *ErrorHandler {
	return &ErrorHandler{bot: bot}
}

func (eh *ErrorHandler) SendMessageWithRetry(chatID, text string, maxRetries int) error {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		_, err := eh.bot.SendMessage(types.MessageConfig{
			ChatID: chatID,
			Text:   text,
		})

		if err == nil {
			return nil
		}

		lastErr = err

		// Check if error is retryable
		if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
			if !zaloBotErr.IsRetryable() {
				return err // Don't retry non-retryable errors
			}
		}

		// Calculate backoff delay
		delay := time.Duration(attempt+1) * time.Second
		log.Printf("Attempt %d failed, retrying in %v...", attempt+1, delay)
		time.Sleep(delay)
	}

	return fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// Example: Error recovery with fallback
func sendMessageWithFallback(bot *zalobot.BotAPI, chatID, primaryText, fallbackText string) error {
	// Try to send primary message
	_, err := bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   primaryText,
	})

	if err != nil {
		log.Printf("Primary message failed: %v", err)
		log.Println("Attempting fallback message...")

		// Try fallback message
		_, err = bot.SendMessage(types.MessageConfig{
			ChatID: chatID,
			Text:   fallbackText,
		})

		if err != nil {
			return fmt.Errorf("both primary and fallback messages failed: %w", err)
		}

		log.Println("Fallback message sent successfully")
		return nil
	}

	log.Println("Primary message sent successfully")
	return nil
}

// Example: Batch operation with error collection
func sendBatchMessages(bot *zalobot.BotAPI, messages []types.MessageConfig) []error {
	errors := make([]error, 0)

	for i, config := range messages {
		_, err := bot.SendMessage(config)
		if err != nil {
			log.Printf("Failed to send message %d: %v", i+1, err)
			errors = append(errors, fmt.Errorf("message %d: %w", i+1, err))
			continue
		}

		log.Printf("Message %d sent successfully", i+1)
	}

	return errors
}
