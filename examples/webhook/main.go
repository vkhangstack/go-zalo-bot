// Package main demonstrates how to use the Zalo Bot SDK with webhooks
// This example shows bot setup with webhooks using bot token,
// webhook signature validation, and event processing
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	zalobot "github.com/vkhangstack/go-zalo-bot"
	"github.com/vkhangstack/go-zalo-bot/types"
)

var bot *zalobot.BotAPI

func main() {
	// Get bot token from environment variable
	botToken := os.Getenv("ZALO_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("ZALO_BOT_TOKEN environment variable is required")
	}

	// Get webhook secret token from environment variable
	webhookSecret := os.Getenv("ZALO_WEBHOOK_SECRET")
	if webhookSecret == "" {
		log.Fatal("ZALO_WEBHOOK_SECRET environment variable is required")
	}

	// Get server port from environment variable (default: 8080)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new bot instance with bot token authentication
	var err error
	bot, err = zalobot.New(botToken,
		types.WithDebug(),
		types.WithEnvironment(types.Production),
	)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	defer bot.Close()

	// Set webhook secret token for signature validation
	bot.SetWebhookSecretToken(webhookSecret)

	log.Printf("Bot started successfully with token: %s...", botToken[:10])
	log.Printf("Webhook secret configured")

	// Set up HTTP server for webhook
	http.HandleFunc("/webhook", webhookHandler)
	http.HandleFunc("/health", healthHandler)

	// Start server in a goroutine
	server := &http.Server{
		Addr: ":" + port,
	}

	go func() {
		log.Printf("Starting webhook server on port %s", port)
		log.Printf("Webhook endpoint: http://localhost:%s/webhook", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	if err := server.Close(); err != nil {
		log.Printf("Error closing server: %v", err)
	}
}

// webhookHandler handles incoming webhook requests
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get signature from header
	signature := r.Header.Get("X-Zalo-Signature")
	if signature == "" {
		log.Println("Missing webhook signature")
		http.Error(w, "Missing signature", http.StatusUnauthorized)
		return
	}

	// Process webhook with signature validation
	update, err := bot.ProcessWebhook(body, signature)
	if err != nil {
		log.Printf("Webhook validation failed: %v", err)
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Log the received update
	log.Printf("Received webhook update #%d", update.UpdateID)

	// Handle the update asynchronously
	go handleUpdate(update)

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok": true,
	})
}

// healthHandler handles health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "healthy",
		"bot":    "running",
	})
}

// handleUpdate processes incoming updates
func handleUpdate(update *types.Update) {
	// Handle text messages
	if update.Message != nil && update.Message.Text != "" {
		handleTextMessage(update.Message)
		return
	}

	// Handle postback events (button clicks)
	if update.PostbackEvent != nil {
		handlePostback(update)
		return
	}

	// Handle user actions (join, leave, block)
	if update.UserAction != nil {
		handleUserAction(update)
		return
	}

	log.Printf("Received update #%d with no recognized content", update.UpdateID)
}

// handleTextMessage processes text messages and sends responses
func handleTextMessage(message *types.Message) {
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
		responseText = "Welcome to Zalo Bot! 👋\n\nThis bot uses webhooks for real-time updates.\n\nAvailable commands:\n/start - Show this message\n/help - Get help\n/echo <text> - Echo your message\n/profile - Get your profile"

	case "/help":
		responseText = "This is a webhook-based demo bot built with Go Zalo Bot SDK.\n\nAll messages are processed in real-time via webhooks!"

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
func handlePostback(update *types.Update) {
	log.Printf("Received postback: %s", update.PostbackEvent.Payload)

	// Process the postback payload
	switch update.PostbackEvent.Payload {
	case "GET_STARTED":
		log.Println("User clicked Get Started button")
	case "HELP":
		log.Println("User clicked Help button")
	default:
		log.Printf("Unknown postback payload: %s", update.PostbackEvent.Payload)
	}
}

// handleUserAction processes user actions like join, leave, block
func handleUserAction(update *types.Update) {
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
