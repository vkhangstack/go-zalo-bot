// Package zalobot provides a comprehensive Go SDK for building Zalo Bot applications.
//
// The SDK supports bot token authentication, message handling, webhooks, polling,
// rich media content, and user profile management.
//
// # Quick Start
//
// Create a bot instance with bot token authentication:
//
//	bot, err := zalobot.New("123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer bot.Close()
//
// Send a text message:
//
//	message, err := bot.SendMessage(types.MessageConfig{
//	    ChatID: "user123",
//	    Text:   "Hello, World!",
//	})
//
// # Polling for Updates
//
// Use polling to receive updates:
//
//	updateConfig := types.UpdateConfig{
//	    Limit:   100,
//	    Timeout: 30, // Long polling timeout in seconds
//	}
//	updates := bot.GetUpdatesChan(updateConfig)
//
//	for update := range updates {
//	    if update.Message != nil {
//	        // Handle message
//	    }
//	}
//
// # Webhooks
//
// Set up webhooks for real-time updates:
//
//	bot.SetWebhookSecretToken("your-webhook-secret")
//	err := bot.SetWebhook(types.WebhookConfig{
//	    URL:         "https://your-domain.com/webhook",
//	    SecretToken: "your-webhook-secret",
//	})
//
// Process webhook requests:
//
//	update, err := bot.ProcessWebhook(payload, signature)
//	if err != nil {
//	    // Invalid signature
//	    return
//	}
//	// Handle update
//
// # Rich Media
//
// Send images, files, videos, and structured messages:
//
//	// Send image
//	bot.SendImage(types.ImageMessageConfig{
//	    ChatID:   "user123",
//	    ImageURL: "https://example.com/image.jpg",
//	    Caption:  "Check this out!",
//	})
//
//	// Send structured message with buttons
//	bot.SendTemplate(types.StructuredMessageConfig{
//	    ChatID: "user123",
//	    StructuredMessage: types.StructuredMessage{
//	        Type: types.StructuredMessageTypeButton,
//	        Elements: []types.MessageElement{...},
//	    },
//	})
//
// # Error Handling
//
// The SDK provides typed errors for better error handling:
//
//	if err != nil {
//	    if zaloBotErr, ok := err.(*types.ZaloBotError); ok {
//	        switch zaloBotErr.Type {
//	        case types.ErrorTypeRateLimit:
//	            // Handle rate limiting
//	        case types.ErrorTypeValidation:
//	            // Handle validation errors
//	        }
//	    }
//	}
//
// # Configuration
//
// Customize bot behavior with options:
//
//	bot, err := zalobot.New(botToken,
//	    zalobot.WithDebug(),                    // Enable debug logging
//	    zalobot.WithTimeout(30*time.Second),    // Set request timeout
//	    zalobot.WithRetries(5),                 // Set max retries
//	    zalobot.WithEnvironment(types.Production),
//	)
//
// For more examples and detailed documentation, see the examples directory
// and visit https://pkg.go.dev/github.com/vkhangstack/go-zalo-bot
package zalobot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/vkhangstack/go-zalo-bot/auth"
	"github.com/vkhangstack/go-zalo-bot/services"
	"github.com/vkhangstack/go-zalo-bot/types"
)

// BotAPI represents the main Zalo Bot API client
type BotAPI struct {
	// Configuration
	config *types.Config

	// Core components
	client      *http.Client
	authService *auth.AuthService

	// Services
	messageService *services.MessageService
	userService    *services.UserService
	webhookService *services.WebhookService

	// Internal state
	mu sync.RWMutex

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// Polling state
	isPolling     bool
	pollingMu     sync.RWMutex
	stopPollingCh chan struct{}
	updatesChan   chan types.Update
	pollingWg     sync.WaitGroup
}

// New creates a new BotAPI instance with bot token authentication
// The bot token should be in format: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
func New(botToken string, options ...types.BotOption) (*BotAPI, error) {
	// Create default config
	config := &types.Config{
		BotToken:    botToken,
		BaseURL:     "https://bot-api.zapps.me",
		Timeout:     30 * time.Second,
		Retries:     3,
		Environment: types.Production,
	}

	// Apply options
	for _, opt := range options {
		opt(config)
	}

	// Validate config (includes bot token validation)
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Create auth service
	authService, err := auth.NewAuthService(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	bot := &BotAPI{
		config:      config,
		client:      config.HTTPClient,
		authService: authService,
		ctx:         ctx,
		cancel:      cancel,
	}

	// Initialize services
	bot.messageService = services.NewMessageService(authService, config.HTTPClient, config)
	bot.userService = services.NewUserService(authService, config.HTTPClient, config)

	// Initialize webhook service with empty secret token (can be set later)
	baseService := services.NewBaseService(authService, config.HTTPClient, config)
	bot.webhookService = services.NewWebhookService(baseService, "")

	return bot, nil
}

// GetBotToken returns the bot token
func (b *BotAPI) GetBotToken() string {
	return b.authService.GetToken()
}

// GetAPIEndpoint returns the full API endpoint URL for a specific method
// Pattern: https://bot-api.zapps.me/bot${BOT_TOKEN}/method
func (b *BotAPI) GetAPIEndpoint(method string) string {
	return b.authService.GetAPIEndpoint(method)
}

// GetConfig returns the bot configuration
func (b *BotAPI) GetConfig() *types.Config {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.config
}

// GetAuthService returns the authentication service
func (b *BotAPI) GetAuthService() *auth.AuthService {
	return b.authService
}

// GetHTTPClient returns the HTTP client
func (b *BotAPI) GetHTTPClient() *http.Client {
	return b.client
}

// GetContext returns the bot context
func (b *BotAPI) GetContext() context.Context {
	return b.ctx
}

// GetMessageService returns the message service
func (b *BotAPI) GetMessageService() *services.MessageService {
	return b.messageService
}

// GetUserService returns the user service
func (b *BotAPI) GetUserService() *services.UserService {
	return b.userService
}

// GetWebhookService returns the webhook service
func (b *BotAPI) GetWebhookService() *services.WebhookService {
	return b.webhookService
}

// SetWebhookSecretToken sets the webhook secret token for signature validation
func (b *BotAPI) SetWebhookSecretToken(token string) {
	if b.webhookService != nil {
		b.webhookService.SetSecretToken(token)
	}
}

// Close closes the bot and releases resources
func (b *BotAPI) Close() {
	// Stop polling if active
	b.StopPolling()

	if b.cancel != nil {
		b.cancel()
	}
}

// SendMessage sends a text message to a chat
// Delegates to the message service
func (b *BotAPI) SendMessage(config types.MessageConfig) (*types.Message, error) {
	return b.messageService.Send(b.ctx, config)
}

// SendImage sends an image message
// Delegates to the message service
func (b *BotAPI) SendImage(config types.ImageMessageConfig) (*types.Message, error) {
	return b.messageService.SendImage(b.ctx, config)
}

// SendFile sends a file message
// Delegates to the message service
func (b *BotAPI) SendFile(config types.FileMessageConfig) (*types.Message, error) {
	return b.messageService.SendFile(b.ctx, config)
}

// SendVideo sends a video message
// Delegates to the message service
func (b *BotAPI) SendVideo(chatID, videoURL, mimeType string) (*types.Message, error) {
	return b.messageService.SendVideo(b.ctx, chatID, videoURL, mimeType)
}

// SendAudio sends an audio message
// Delegates to the message service
func (b *BotAPI) SendAudio(chatID, audioURL, mimeType string) (*types.Message, error) {
	return b.messageService.SendAudio(b.ctx, chatID, audioURL, mimeType)
}

// SendTemplate sends a structured message with buttons and quick replies
// Delegates to the message service
func (b *BotAPI) SendTemplate(config types.StructuredMessageConfig) (*types.Message, error) {
	return b.messageService.SendTemplate(b.ctx, config)
}

// SendStructuredMessage sends a structured message (alias for SendTemplate)
// Delegates to the message service
func (b *BotAPI) SendStructuredMessage(config types.StructuredMessageConfig) (*types.Message, error) {
	return b.messageService.SendStructuredMessage(b.ctx, config)
}

// GetUserProfile retrieves user profile information
// Delegates to the user service
func (b *BotAPI) GetUserProfile(userID string) (*types.UserProfile, error) {
	return b.userService.GetUserProfile(b.ctx, userID)
}

// ProcessWebhook processes a webhook request with signature validation
// Delegates to the webhook service
func (b *BotAPI) ProcessWebhook(payload []byte, signature string) (*types.Update, error) {
	return b.webhookService.ProcessWebhook(payload, signature)
}

// ValidateWebhookSignature validates a webhook signature
// Delegates to the webhook service
func (b *BotAPI) ValidateWebhookSignature(payload []byte, signature string) error {
	return b.webhookService.ValidateSignature(payload, signature)
}

// ParseWebhookUpdate parses a webhook payload into an Update structure
// Delegates to the webhook service
func (b *BotAPI) ParseWebhookUpdate(payload []byte) (*types.Update, error) {
	return b.webhookService.ParseUpdate(payload)
}

// SetWebhook sets the webhook URL for receiving updates
// The webhook URL must be a valid HTTPS URL
func (b *BotAPI) SetWebhook(config types.WebhookConfig) error {
	// Validate webhook URL
	if err := validateWebhookURL(config.URL); err != nil {
		return types.NewValidationError(err.Error())
	}

	// Prepare request body
	requestBody := map[string]interface{}{
		"url":          config.URL,
		"secret_token": config.SecretToken,
	}

	// Construct URL with bot token embedded
	url := b.authService.GetAPIEndpoint("setWebhook")

	// Prepare request body
	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return types.NewValidationError("failed to marshal request body")
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(b.ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return types.NewNetworkError("failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", services.UserAgent())

	// Execute request
	resp, err := b.client.Do(req)
	if err != nil {
		return types.NewNetworkError("request failed")
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.NewNetworkError("failed to read response")
	}

	// Parse response
	var apiResp struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code,omitempty"`
		Description string `json:"description,omitempty"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return types.NewAPIError(resp.StatusCode, "failed to parse response", err.Error())
	}

	if !apiResp.OK {
		return types.NewAPIError(apiResp.ErrorCode, apiResp.Description, "setWebhook failed")
	}

	// Set the secret token in the webhook service for signature validation
	if config.SecretToken != "" {
		b.SetWebhookSecretToken(config.SecretToken)
	}

	return nil
}

// DeleteWebhook removes the webhook configuration
func (b *BotAPI) DeleteWebhook() error {
	// Construct URL with bot token embedded
	url := b.authService.GetAPIEndpoint("deleteWebhook")

	// Create HTTP request
	req, err := http.NewRequestWithContext(b.ctx, "POST", url, nil)
	if err != nil {
		return types.NewNetworkError("failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", services.UserAgent())

	// Execute request
	resp, err := b.client.Do(req)
	if err != nil {
		return types.NewNetworkError("request failed")
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.NewNetworkError("failed to read response")
	}

	// Parse response
	var apiResp struct {
		OK          bool   `json:"ok"`
		ErrorCode   int    `json:"error_code,omitempty"`
		Description string `json:"description,omitempty"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return types.NewAPIError(resp.StatusCode, "failed to parse response", err.Error())
	}

	if !apiResp.OK {
		return types.NewAPIError(apiResp.ErrorCode, apiResp.Description, "deleteWebhook failed")
	}

	return nil
}

// GetWebhookInfo retrieves information about the current webhook configuration
func (b *BotAPI) GetWebhookInfo() (*types.WebhookInfo, error) {
	// Construct URL with bot token embedded
	url := b.authService.GetAPIEndpoint("getWebhookInfo")

	// Create HTTP request
	req, err := http.NewRequestWithContext(b.ctx, "GET", url, nil)
	if err != nil {
		return nil, types.NewNetworkError("failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", services.UserAgent())

	// Execute request
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, types.NewNetworkError("request failed")
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.NewNetworkError("failed to read response")
	}

	// Parse response
	var apiResp struct {
		OK          bool               `json:"ok"`
		Result      *types.WebhookInfo `json:"result,omitempty"`
		ErrorCode   int                `json:"error_code,omitempty"`
		Description string             `json:"description,omitempty"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, types.NewAPIError(resp.StatusCode, "failed to parse response", err.Error())
	}

	if !apiResp.OK {
		return nil, types.NewAPIError(apiResp.ErrorCode, apiResp.Description, "getWebhookInfo failed")
	}

	return apiResp.Result, nil
}

func (b *BotAPI) GetFieldSignature() string {
	return b.userService.GetFieldSignature()
}

// validateWebhookURL validates that the webhook URL is properly formatted and uses HTTPS
func validateWebhookURL(webhookURL string) error {
	if strings.TrimSpace(webhookURL) == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	// Parse URL
	parsedURL, err := url.Parse(webhookURL)
	if err != nil {
		return fmt.Errorf("invalid webhook URL format: %w", err)
	}

	// Check scheme
	if parsedURL.Scheme != "https" {
		return fmt.Errorf("webhook URL must use HTTPS protocol")
	}

	// Check host
	if parsedURL.Host == "" {
		return fmt.Errorf("webhook URL must have a valid host")
	}

	return nil
}

// GetUpdates retrieves updates from the Zalo Bot API using polling
// Implements Requirements 3.1, 5.3
func (b *BotAPI) GetUpdates(config types.UpdateConfig) ([]types.Update, error) {
	return b.GetUpdatesWithContext(b.ctx, config)
}

// GetUpdatesWithContext retrieves updates with a custom context
func (b *BotAPI) GetUpdatesWithContext(ctx context.Context, config types.UpdateConfig) ([]types.Update, error) {
	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Construct URL with bot token embedded
	endpoint := b.authService.GetAPIEndpoint("getUpdates")

	// Add query parameters
	params := url.Values{}
	if config.Offset > 0 {
		params.Add("offset", fmt.Sprintf("%d", config.Offset))
	}
	if config.Limit > 0 {
		params.Add("limit", fmt.Sprintf("%d", config.Limit))
	}
	if config.Timeout > 0 {
		params.Add("timeout", fmt.Sprintf("%d", config.Timeout))
	}

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	// Create HTTP request with context
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, types.NewNetworkError("failed to create request")
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", services.UserAgent())

	// Execute request
	resp, err := b.client.Do(req)
	if err != nil {
		return nil, types.NewNetworkError("request failed")
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.NewNetworkError("failed to read response")
	}

	// Parse response
	var apiResp struct {
		OK          bool           `json:"ok"`
		Result      []types.Update `json:"result,omitempty"`
		ErrorCode   int            `json:"error_code,omitempty"`
		Description string         `json:"description,omitempty"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, types.NewAPIError(resp.StatusCode, "failed to parse response", err.Error())
	}

	if !apiResp.OK {
		return nil, types.NewAPIError(apiResp.ErrorCode, apiResp.Description, "getUpdates failed")
	}

	return apiResp.Result, nil
}

// GetUpdatesChan returns a channel that receives updates via polling
// The polling loop runs in a background goroutine with configurable intervals and error handling
// Implements Requirements 3.1, 5.3
func (b *BotAPI) GetUpdatesChan(config types.UpdateConfig) <-chan types.Update {
	// Validate config
	if err := config.Validate(); err != nil {
		// Return a closed channel if config is invalid
		ch := make(chan types.Update)
		close(ch)
		return ch
	}

	b.pollingMu.Lock()
	defer b.pollingMu.Unlock()

	// If already polling, return existing channel
	if b.isPolling && b.updatesChan != nil {
		return b.updatesChan
	}

	// Create updates channel
	if b.updatesChan == nil {
		b.updatesChan = make(chan types.Update, 100)
	}
	if b.stopPollingCh == nil {
		b.stopPollingCh = make(chan struct{})
	}
	b.isPolling = true

	// Start polling in background goroutine
	b.pollingWg.Add(1)
	go b.pollUpdates(config)

	return b.updatesChan
}

// pollUpdates runs the polling loop
func (b *BotAPI) pollUpdates(config types.UpdateConfig) {
	defer b.pollingWg.Done()

	offset := config.Offset
	pollInterval := 1 * time.Second // Default polling interval

	// Use long polling if timeout is set
	if config.Timeout > 0 {
		pollInterval = time.Duration(config.Timeout) * time.Second
	}

	// Create a context that can be cancelled when stopping
	pollCtx, pollCancel := context.WithCancel(b.ctx)
	defer pollCancel()

	// Monitor stop signal in a separate goroutine
	go func() {
		select {
		case <-b.stopPollingCh:
			pollCancel()
		case <-b.ctx.Done():
			pollCancel()
		}
	}()

	for {
		select {
		case <-pollCtx.Done():
			// Stop polling
			b.pollingMu.Lock()
			if b.updatesChan != nil {
				close(b.updatesChan)
				b.updatesChan = nil
			}
			b.isPolling = false
			b.pollingMu.Unlock()
			return

		default:
			// Get updates with cancellable context
			updateConfig := types.UpdateConfig{
				Offset:  offset,
				Limit:   config.Limit,
				Timeout: config.Timeout,
			}

			updates, err := b.GetUpdatesWithContext(pollCtx, updateConfig)
			if err != nil {
				// Check if context was cancelled
				if pollCtx.Err() != nil {
					b.pollingMu.Lock()
					if b.updatesChan != nil {
						close(b.updatesChan)
						b.updatesChan = nil
					}
					b.isPolling = false
					b.pollingMu.Unlock()
					return
				}

				// Log error if debug mode is enabled
				if b.config.Debug {
					fmt.Printf("Error getting updates: %v\n", err)
				}

				// Wait before retrying
				select {
				case <-pollCtx.Done():
					b.pollingMu.Lock()
					if b.updatesChan != nil {
						close(b.updatesChan)
						b.updatesChan = nil
					}
					b.isPolling = false
					b.pollingMu.Unlock()
					return
				case <-time.After(pollInterval):
					continue
				}
			}

			// Send updates to channel
			for _, update := range updates {
				select {
				case b.updatesChan <- update:
					// Update sent successfully
					// Update offset to acknowledge this update
					if update.UpdateID >= offset {
						offset = update.UpdateID + 1
					}
				case <-pollCtx.Done():
					// Stop polling
					b.pollingMu.Lock()
					if b.updatesChan != nil {
						close(b.updatesChan)
						b.updatesChan = nil
					}
					b.isPolling = false
					b.pollingMu.Unlock()
					return
				}
			}

			// If no updates received and not using long polling, wait before next poll
			if len(updates) == 0 && config.Timeout == 0 {
				select {
				case <-pollCtx.Done():
					b.pollingMu.Lock()
					if b.updatesChan != nil {
						close(b.updatesChan)
						b.updatesChan = nil
					}
					b.isPolling = false
					b.pollingMu.Unlock()
					return
				case <-time.After(pollInterval):
					// Continue polling
				}
			}
		}
	}
}

// IsPolling returns true if the bot is currently polling for updates
func (b *BotAPI) IsPolling() bool {
	b.pollingMu.RLock()
	defer b.pollingMu.RUnlock()
	return b.isPolling
}

// StopPolling stops the polling loop
func (b *BotAPI) StopPolling() {
	b.pollingMu.Lock()
	if b.isPolling && b.stopPollingCh != nil {
		close(b.stopPollingCh)
		b.stopPollingCh = nil
	}
	b.pollingMu.Unlock()

	// Wait for polling goroutine to finish
	b.pollingWg.Wait()
}
