// Package main demonstrates how to use the Zalo Bot SDK with webhooks
// This example shows bot setup with webhooks using bot token,
// webhook secret token validation, and event processing
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

	// Set webhook secret token, compared against the X-Bot-Api-Secret-Token
	// header Zalo sends with every webhook request
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

	// Get the secret token from the header Zalo sends with every webhook
	// request (X-Bot-Api-Secret-Token) and re-validate it before processing.
	secretToken := r.Header.Get(bot.GetFieldSecretToken())
	if secretToken == "" {
		log.Println("Missing webhook secret token")
		http.Error(w, "Missing secret token", http.StatusForbidden)
		return
	}

	// Process webhook with secret token validation
	update, err := bot.ProcessWebhook(body, secretToken)
	if err != nil {
		log.Printf("Webhook validation failed: %v", err)
		http.Error(w, "Invalid secret token", http.StatusForbidden)
		return
	}

	// Log the received update
	log.Printf("Received webhook event: %s", update.EventName)

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

// handleUpdate processes incoming updates.
//
// Per https://bot.zapps.me/docs/webhook/, webhooks only ever deliver
// message.* events - there is no postback or user_action event from the
// webhook itself.
func handleUpdate(update *types.Update) {
	switch update.EventName {
	case types.EventMessageText:
		handleTextMessage(update.Message)

	case types.EventMessageImage:
		handleImageMessage(update.Message)

	case types.EventMessageSticker:
		handleStickerMessage(update.Message)

	case types.EventMessageVoice:
		handleVoiceMessage(update.Message)

	case types.EventMessageUnsupported:
		// Content withheld for a "special audience" account (minors, etc.)
		// to comply with applicable regulations - nothing to process.
		log.Printf("Received unsupported/withheld message event")

	default:
		log.Printf("Received unrecognized webhook event: %s", update.EventName)
	}
}

// handleTextMessage processes text messages and sends responses
func handleTextMessage(message *types.Message) {
	log.Printf("Received message from %s: %s", message.From.ID, message.Text)

	chatID := chatIDFor(message)

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

	sendReply(chatID, responseText)
}

// handleImageMessage processes incoming images (sent as a plain URL)
func handleImageMessage(message *types.Message) {
	imageURL := ""
	if message.Photo != nil {
		imageURL = message.Photo.URL
	}
	log.Printf("Received image from %s: %s (caption: %q)", message.From.ID, imageURL, message.Caption)
	sendReply(chatIDFor(message), "Thanks for the image!")
}

// handleStickerMessage processes incoming stickers
func handleStickerMessage(message *types.Message) {
	stickerID, stickerURL := "", ""
	if message.Sticker != nil {
		stickerID = message.Sticker.FileID
		stickerURL = message.Sticker.URL
	}
	log.Printf("Received sticker from %s: %s (%s)", message.From.ID, stickerID, stickerURL)
	sendReply(chatIDFor(message), "Nice sticker!")
}

// handleVoiceMessage processes incoming voice messages
func handleVoiceMessage(message *types.Message) {
	log.Printf("Received voice message from %s: %s", message.From.ID, message.VoiceURL)
	sendReply(chatIDFor(message), "Got your voice message!")
}

// chatIDFor returns the chat ID to reply to for an incoming message
func chatIDFor(message *types.Message) string {
	if message.Chat != nil {
		return message.Chat.ID
	}
	return message.From.ID
}

// sendReply sends a text response back to the given chat
func sendReply(chatID, text string) {
	_, err := bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	} else {
		log.Printf("Sent response to %s", chatID)
	}
}
