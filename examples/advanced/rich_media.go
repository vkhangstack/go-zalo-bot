// Package main demonstrates advanced usage of the Zalo Bot SDK
// This example shows how to send rich media messages including images, files, and structured messages
package main

import (
	"fmt"
	"log"
	"os"

	zalobot "github.com/vkhangstack/go-zalo-bot"
	"github.com/vkhangstack/go-zalo-bot/types"
)

func main() {
	// Get bot token from environment variable
	botToken := os.Getenv("ZALO_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("ZALO_BOT_TOKEN environment variable is required")
	}

	// Create bot instance
	bot, err := zalobot.New(botToken, types.WithDebug())
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	defer bot.Close()

	// Example chat ID (replace with actual user ID)
	chatID := "example_user_id"

	// Demonstrate different types of rich media messages
	demonstrateImageMessage(bot, chatID)
	demonstrateFileMessage(bot, chatID)
	demonstrateVideoMessage(bot, chatID)
	demonstrateAudioMessage(bot, chatID)
	demonstrateStructuredMessage(bot, chatID)
	demonstrateButtonMessage(bot, chatID)
	demonstrateQuickReplyMessage(bot, chatID)
}

// demonstrateImageMessage shows how to send an image message
func demonstrateImageMessage(bot *zalobot.BotAPI, chatID string) {
	log.Println("Sending image message...")

	config := types.ImageMessageConfig{
		ChatID:   chatID,
		ImageURL: "https://example.com/image.jpg",
		Caption:  "Check out this image! 📸",
		MimeType: "image/jpeg",
	}

	message, err := bot.SendImage(config)
	if err != nil {
		log.Printf("Failed to send image: %v", err)
		return
	}

	log.Printf("Image message sent successfully: %s", message.MessageID)
}

// demonstrateFileMessage shows how to send a file message
func demonstrateFileMessage(bot *zalobot.BotAPI, chatID string) {
	log.Println("Sending file message...")

	config := types.FileMessageConfig{
		ChatID:   chatID,
		FileURL:  "https://example.com/document.pdf",
		FileName: "document.pdf",
		MimeType: "application/pdf",
		Size:     1024 * 1024, // 1MB
	}

	message, err := bot.SendFile(config)
	if err != nil {
		log.Printf("Failed to send file: %v", err)
		return
	}

	log.Printf("File message sent successfully: %s", message.MessageID)
}

// demonstrateVideoMessage shows how to send a video message
func demonstrateVideoMessage(bot *zalobot.BotAPI, chatID string) {
	log.Println("Sending video message...")

	message, err := bot.SendVideo(chatID, "https://example.com/video.mp4", "video/mp4")
	if err != nil {
		log.Printf("Failed to send video: %v", err)
		return
	}

	log.Printf("Video message sent successfully: %s", message.MessageID)
}

// demonstrateAudioMessage shows how to send an audio message
func demonstrateAudioMessage(bot *zalobot.BotAPI, chatID string) {
	log.Println("Sending audio message...")

	message, err := bot.SendAudio(chatID, "https://example.com/audio.mp3", "audio/mpeg")
	if err != nil {
		log.Printf("Failed to send audio: %v", err)
		return
	}

	log.Printf("Audio message sent successfully: %s", message.MessageID)
}

// demonstrateStructuredMessage shows how to send a structured message with elements
func demonstrateStructuredMessage(bot *zalobot.BotAPI, chatID string) {
	log.Println("Sending structured message...")

	config := types.StructuredMessageConfig{
		ChatID: chatID,
		StructuredMessage: types.StructuredMessage{
			Type: types.StructuredMessageTypeTemplate,
			Elements: []types.MessageElement{
				{
					Title:    "Product 1",
					Subtitle: "This is an amazing product",
					ImageURL: "https://example.com/product1.jpg",
					Buttons: []types.Button{
						{
							Type:    types.ButtonTypeWebURL,
							Title:   "View Details",
							URL:     "https://example.com/product1",
						},
						{
							Type:    types.ButtonTypePostback,
							Title:   "Buy Now",
							Payload: "BUY_PRODUCT_1",
						},
					},
				},
				{
					Title:    "Product 2",
					Subtitle: "Another great product",
					ImageURL: "https://example.com/product2.jpg",
					Buttons: []types.Button{
						{
							Type:    types.ButtonTypeWebURL,
							Title:   "View Details",
							URL:     "https://example.com/product2",
						},
						{
							Type:    types.ButtonTypePostback,
							Title:   "Buy Now",
							Payload: "BUY_PRODUCT_2",
						},
					},
				},
			},
		},
	}

	message, err := bot.SendStructuredMessage(config)
	if err != nil {
		log.Printf("Failed to send structured message: %v", err)
		return
	}

	log.Printf("Structured message sent successfully: %s", message.MessageID)
}

// demonstrateButtonMessage shows how to send a message with buttons
func demonstrateButtonMessage(bot *zalobot.BotAPI, chatID string) {
	log.Println("Sending button message...")

	config := types.StructuredMessageConfig{
		ChatID: chatID,
		StructuredMessage: types.StructuredMessage{
			Type: types.StructuredMessageTypeButton,
			Elements: []types.MessageElement{
				{
					Title:    "Choose an option",
					Subtitle: "Please select one of the following:",
					Buttons: []types.Button{
						{
							Type:    types.ButtonTypePostback,
							Title:   "Option 1",
							Payload: "OPTION_1",
						},
						{
							Type:    types.ButtonTypePostback,
							Title:   "Option 2",
							Payload: "OPTION_2",
						},
						{
							Type:    types.ButtonTypeWebURL,
							Title:   "Learn More",
							URL:     "https://example.com/learn-more",
						},
					},
				},
			},
		},
	}

	message, err := bot.SendTemplate(config)
	if err != nil {
		log.Printf("Failed to send button message: %v", err)
		return
	}

	log.Printf("Button message sent successfully: %s", message.MessageID)
}

// demonstrateQuickReplyMessage shows how to send a message with quick replies
func demonstrateQuickReplyMessage(bot *zalobot.BotAPI, chatID string) {
	log.Println("Sending quick reply message...")

	config := types.StructuredMessageConfig{
		ChatID: chatID,
		StructuredMessage: types.StructuredMessage{
			Type: types.StructuredMessageTypeTemplate,
			QuickReplies: []types.QuickReply{
				{
					ContentType: types.QuickReplyTypeText,
					Title:       "Yes",
					Payload:     "QUICK_REPLY_YES",
				},
				{
					ContentType: types.QuickReplyTypeText,
					Title:       "No",
					Payload:     "QUICK_REPLY_NO",
				},
				{
					ContentType: types.QuickReplyTypeLocation,
					Title:       "Share Location",
					Payload:     "QUICK_REPLY_LOCATION",
				},
			},
		},
	}

	message, err := bot.SendStructuredMessage(config)
	if err != nil {
		log.Printf("Failed to send quick reply message: %v", err)
		return
	}

	log.Printf("Quick reply message sent successfully: %s", message.MessageID)
}

// Example: Sending multiple media types in sequence
func sendMediaSequence(bot *zalobot.BotAPI, chatID string) {
	fmt.Println("\n=== Sending Media Sequence ===")

	// Send text introduction
	bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   "Here's a collection of media for you:",
	})

	// Send image
	bot.SendImage(types.ImageMessageConfig{
		ChatID:   chatID,
		ImageURL: "https://example.com/image1.jpg",
		Caption:  "Image 1",
		MimeType: "image/jpeg",
	})

	// Send another image
	bot.SendImage(types.ImageMessageConfig{
		ChatID:   chatID,
		ImageURL: "https://example.com/image2.jpg",
		Caption:  "Image 2",
		MimeType: "image/jpeg",
	})

	// Send video
	bot.SendVideo(chatID, "https://example.com/video.mp4", "video/mp4")

	// Send closing message
	bot.SendMessage(types.MessageConfig{
		ChatID: chatID,
		Text:   "That's all! Hope you enjoyed the media.",
	})

	fmt.Println("Media sequence sent successfully")
}
