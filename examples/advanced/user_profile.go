// Package main demonstrates user profile handling with the Zalo Bot SDK
// This example shows how to retrieve and work with user profile information
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

	// Example user ID (replace with actual user ID)
	userID := "example_user_id"

	// Demonstrate user profile operations
	getUserProfile(bot, userID)
	sendPersonalizedMessage(bot, userID)
	handleMultipleUsers(bot, []string{"user1", "user2", "user3"})
}

// getUserProfile retrieves and displays user profile information
func getUserProfile(bot *zalobot.BotAPI, userID string) {
	log.Printf("Retrieving profile for user: %s", userID)

	profile, err := bot.GetUserProfile(userID)
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		return
	}

	// Display profile information
	fmt.Println("\n=== User Profile ===")
	fmt.Printf("ID: %s\n", profile.ID)
	fmt.Printf("Name: %s\n", profile.Name)
	fmt.Printf("Avatar: %s\n", profile.Avatar)
	fmt.Println("===================\n")
}

// sendPersonalizedMessage sends a personalized message using user profile data
func sendPersonalizedMessage(bot *zalobot.BotAPI, userID string) {
	log.Println("Sending personalized message...")

	// Get user profile
	profile, err := bot.GetUserProfile(userID)
	if err != nil {
		log.Printf("Failed to get user profile: %v", err)
		return
	}

	// Create personalized message
	message := fmt.Sprintf("Hello %s! 👋\n\nWelcome to our bot. How can I help you today?", profile.Name)

	// Send message
	_, err = bot.SendMessage(types.MessageConfig{
		ChatID: userID,
		Text:   message,
	})
	if err != nil {
		log.Printf("Failed to send personalized message: %v", err)
		return
	}

	log.Println("Personalized message sent successfully")
}

// handleMultipleUsers demonstrates handling multiple user profiles
func handleMultipleUsers(bot *zalobot.BotAPI, userIDs []string) {
	fmt.Println("\n=== Processing Multiple Users ===")

	for _, userID := range userIDs {
		profile, err := bot.GetUserProfile(userID)
		if err != nil {
			log.Printf("Failed to get profile for user %s: %v", userID, err)
			continue
		}

		fmt.Printf("User: %s (ID: %s)\n", profile.Name, profile.ID)

		// Send welcome message to each user
		bot.SendMessage(types.MessageConfig{
			ChatID: userID,
			Text:   fmt.Sprintf("Hi %s! Thanks for using our bot.", profile.Name),
		})
	}

	fmt.Println("=================================\n")
}

// Example: Building a user directory
type UserDirectory struct {
	users map[string]*types.UserProfile
}

func NewUserDirectory() *UserDirectory {
	return &UserDirectory{
		users: make(map[string]*types.UserProfile),
	}
}

func (ud *UserDirectory) AddUser(bot *zalobot.BotAPI, userID string) error {
	profile, err := bot.GetUserProfile(userID)
	if err != nil {
		return fmt.Errorf("failed to get user profile: %w", err)
	}

	ud.users[userID] = profile
	return nil
}

func (ud *UserDirectory) GetUser(userID string) (*types.UserProfile, bool) {
	profile, exists := ud.users[userID]
	return profile, exists
}

func (ud *UserDirectory) ListUsers() []*types.UserProfile {
	profiles := make([]*types.UserProfile, 0, len(ud.users))
	for _, profile := range ud.users {
		profiles = append(profiles, profile)
	}
	return profiles
}

// Example: User profile caching to reduce API calls
type ProfileCache struct {
	bot   *zalobot.BotAPI
	cache map[string]*types.UserProfile
}

func NewProfileCache(bot *zalobot.BotAPI) *ProfileCache {
	return &ProfileCache{
		bot:   bot,
		cache: make(map[string]*types.UserProfile),
	}
}

func (pc *ProfileCache) GetProfile(userID string) (*types.UserProfile, error) {
	// Check cache first
	if profile, exists := pc.cache[userID]; exists {
		log.Printf("Profile cache hit for user: %s", userID)
		return profile, nil
	}

	// Cache miss - fetch from API
	log.Printf("Profile cache miss for user: %s", userID)
	profile, err := pc.bot.GetUserProfile(userID)
	if err != nil {
		return nil, err
	}

	// Store in cache
	pc.cache[userID] = profile
	return profile, nil
}

func (pc *ProfileCache) InvalidateCache(userID string) {
	delete(pc.cache, userID)
}

func (pc *ProfileCache) ClearCache() {
	pc.cache = make(map[string]*types.UserProfile)
}

// Example usage of profile cache
func demonstrateProfileCache(bot *zalobot.BotAPI) {
	fmt.Println("\n=== Profile Cache Demo ===")

	cache := NewProfileCache(bot)

	userID := "example_user_id"

	// First call - cache miss
	profile1, err := cache.GetProfile(userID)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("First call: %s\n", profile1.Name)

	// Second call - cache hit
	profile2, err := cache.GetProfile(userID)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Second call: %s\n", profile2.Name)

	fmt.Println("==========================\n")
}
