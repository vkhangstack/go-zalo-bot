// Package main demonstrates how to use the Zalo Bot SDK with polling
// This example shows bot setup with polling using bot token,
// message handling, and response generation
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	zalobot "github.com/vkhangstack/go-zalo-bot"
	"github.com/vkhangstack/go-zalo-bot/types"
)

func main() {
	// Get bot token from environment variable
	botToken := os.Getenv("ZALO_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("ZALO_BOT_TOKEN environment variable is required")
	}

	// Create a new bot instance with bot token authentication
	// The bot token should be in format: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	bot, err := zalobot.New(botToken,
		types.WithDebug(),
		types.WithEnvironment(types.Production),
	)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	defer bot.Close()

	log.Printf("Bot started successfully with token: %s...", botToken[:10])
	log.Println("Polling for updates...")

	// Configure update polling
	updateConfig := types.UpdateConfig{
		Offset:  0,
		Limit:   100,
		Timeout: 30, // Long polling timeout in seconds
	}

	// Get updates channel for polling
	updates := bot.GetUpdatesChan(updateConfig)

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Process updates
	for {
		select {
		case update, ok := <-updates:
			if !ok {
				log.Println("Updates channel closed")
				return
			}

			// Handle the update
			handleUpdate(bot, update)

		case <-sigChan:
			log.Println("Shutting down gracefully...")
			bot.StopPolling()
			return
		}
	}
}

// handleUpdate processes incoming updates and generates responses
func handleUpdate(bot *zalobot.BotAPI, update types.Update) {
	// Handle text messages
	if update.Message != nil && update.Message.Text != "" {
		handleTextMessage(bot, update.Message)
		return
	}

	// Handle postback events (button clicks)
	if update.PostbackEvent != nil {
		handlePostback(bot, update)
		return
	}

	// Handle user actions (join, leave, block)
	if update.UserAction != nil {
		handleUserAction(bot, update)
		return
	}

	log.Printf("Received update #%d with no recognized content", update.UpdateID)
}

// handleTextMessage processes text messages and sends responses
func handleTextMessage(bot *zalobot.BotAPI, message *types.Message) {
	log.Printf("Received message from %s: %s", message.From.ID, message.Text)

	// Get chat ID for response
	chatID := message.From.ID
	if message.Chat != nil {
		chatID = message.Chat.ID
	}

	// Generate response based on message content
	var responseText string
	switch message.Text {
	case "/start":
		responseText = "Welcome to Zalo Bot! 👋\n\nAvailable commands:\n/start - Show this message\n/help - Get help\n/echo <text> - Echo your message\n/profile - Get your profile"

	case "/help":
		responseText = "This is a demo bot built with Go Zalo Bot SDK.\n\nSend me any message and I'll respond!"

	case "/profile":
		// Get user profile
		profile, err := bot.GetUserProfile(message.From.ID)
		if err != nil {
			log.Printf("Failed to get user profile: %v", err)
			responseText = "Sorry, I couldn't retrieve your profile."
		} else {
			responseText = fmt.Sprintf("Your Profile:\nName: %s\nID: %s", profile.Name, profile.ID)
		}

	default:
		// Echo the message back
		responseText = fmt.Sprintf("You said: %s", message.Text)
	}

	// Send response
	_, err := bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   responseText,
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		log.Printf("Sent response to %s", chatID)
	}
}

// handlePostback processes postback events from button clicks
func handlePostback(bot *zalobot.BotAPI, update types.Update) {
	log.Printf("Received postback: %s", update.PostbackEvent.Payload)

	// Extract user ID from the update
	// In a real scenario, you'd get this from the postback event
	// For now, we'll skip sending a response
	log.Println("Postback event processed")
}

// handleUserAction processes user actions like join, leave, block
func handleUserAction(bot *zalobot.BotAPI, update types.Update) {
	log.Printf("User action: %s by user %s", update.UserAction.Type, update.UserAction.UserID)

	// Handle different user actions
	switch update.UserAction.Type {
	case types.UserActionTypeJoin:
		// Send welcome message
		_, err := bot.SendMessage(types.MessageConfig{
			ChatID: update.UserAction.UserID,
			Text:   "Welcome! Thanks for joining! 🎉",
		})
		if err != nil {
			log.Printf("Failed to send welcome message: %v", err)
		}

	case types.UserActionTypeLeave:
		log.Printf("User %s left", update.UserAction.UserID)

	case types.UserActionTypeBlock:
		log.Printf("User %s blocked the bot", update.UserAction.UserID)
	}
}
