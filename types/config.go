package types

import (
	"fmt"
	"net/http"
	"time"
)

// Config represents the configuration for the Zalo Bot SDK
type Config struct {
	BotToken    string        // Bot token for authentication
	BaseURL     string        // Base API URL (default: https://bot-api.zapps.me)
	Debug       bool
	Timeout     time.Duration
	Retries     int
	Environment Environment   // Development or Production
	HTTPClient  *http.Client
	RetryConfig *RetryConfig  // Retry configuration for error handling
}

// MessageConfig represents configuration for sending messages
type MessageConfig struct {
	ChatID      string
	Text        string
	MessageType MessageType
	Attachments []Attachment
}

// WebhookConfig represents configuration for webhook setup
type WebhookConfig struct {
	URL         string
	SecretToken string
	Certificate string
}

// UpdateConfig represents configuration for getting updates
type UpdateConfig struct {
	Offset  int
	Limit   int
	Timeout int
}

// BotOption represents a configuration option for the bot
type BotOption func(*Config)

// WithBaseURL sets the base URL for the bot API
func WithBaseURL(url string) BotOption {
	return func(c *Config) { c.BaseURL = url }
}

// WithTimeout sets the timeout for HTTP requests
func WithTimeout(timeout time.Duration) BotOption {
	return func(c *Config) { c.Timeout = timeout }
}

// WithDebug enables debug mode
func WithDebug() BotOption {
	return func(c *Config) { c.Debug = true }
}

// WithEnvironment sets the environment
func WithEnvironment(env Environment) BotOption {
	return func(c *Config) { c.Environment = env }
}

// WithRetries sets the number of retries
func WithRetries(retries int) BotOption {
	return func(c *Config) { c.Retries = retries }
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) BotOption {
	return func(c *Config) { c.HTTPClient = client }
}

// WithRetryConfig sets a custom retry configuration
func WithRetryConfig(retryConfig *RetryConfig) BotOption {
	return func(c *Config) { c.RetryConfig = retryConfig }
}

// ImageMessageConfig represents configuration for sending image messages
type ImageMessageConfig struct {
	ChatID   string
	ImageURL string
	Caption  string
	MimeType string
}

// FileMessageConfig represents configuration for sending file messages
type FileMessageConfig struct {
	ChatID   string
	FileURL  string
	FileName string
	MimeType string
	Size     int64
}

// StructuredMessageConfig represents configuration for sending structured messages
type StructuredMessageConfig struct {
	ChatID            string
	StructuredMessage StructuredMessage
}

// Environment represents the environment type
type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

// IsValid validates the environment
func (e Environment) IsValid() bool {
	switch e {
	case Development, Production:
		return true
	default:
		return false
	}
}

// String returns the string representation of Environment
func (e Environment) String() string {
	return string(e)
}

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeText       MessageType = "text"
	MessageTypeImage      MessageType = "image"
	MessageTypeFile       MessageType = "file"
	MessageTypeTemplate   MessageType = "template"
	MessageTypeInteractive MessageType = "interactive"
)

// IsValid validates the message type
func (mt MessageType) IsValid() bool {
	switch mt {
	case MessageTypeText, MessageTypeImage, MessageTypeFile, MessageTypeTemplate, MessageTypeInteractive:
		return true
	default:
		return false
	}
}

// String returns the string representation of MessageType
func (mt MessageType) String() string {
	return string(mt)
}

// GetAPIEndpoint returns the full API endpoint with bot token for a specific method
func (c *Config) GetAPIEndpoint(method string) string {
	return fmt.Sprintf("%s/bot%s/%s", c.BaseURL, c.BotToken, method)
}

// Validate validates the Config
func (c *Config) Validate() error {
	if c.BotToken == "" {
		return &ZaloBotError{
			Code:    400,
			Message: "Bot token is required",
			Type:    ErrorTypeValidation,
		}
	}
	
	// Validate bot token format (should be like "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11")
	if !isValidBotTokenFormat(c.BotToken) {
		return &ZaloBotError{
			Code:    400,
			Message: "Invalid bot token format",
			Type:    ErrorTypeValidation,
		}
	}
	
	if c.BaseURL == "" {
		c.BaseURL = "https://bot-api.zapps.me"
	}
	
	if c.Timeout == 0 {
		c.Timeout = 30 * time.Second
	}
	
	if c.Retries == 0 {
		c.Retries = 3
	}
	
	if c.Environment == "" {
		c.Environment = Production
	}
	
	if !c.Environment.IsValid() {
		return &ZaloBotError{
			Code:    400,
			Message: "Invalid environment",
			Type:    ErrorTypeValidation,
		}
	}
	
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Timeout: c.Timeout,
		}
	}
	
	if c.RetryConfig == nil {
		c.RetryConfig = DefaultRetryConfig()
	}
	
	return nil
}

// isValidBotTokenFormat validates the bot token format
func isValidBotTokenFormat(token string) bool {
	// Bot token should be in format: "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	// At minimum, it should contain a colon separating bot ID and token
	if len(token) < 10 {
		return false
	}
	
	colonIndex := -1
	for i, char := range token {
		if char == ':' {
			colonIndex = i
			break
		}
	}
	
	// Must have a colon and both parts should be non-empty
	if colonIndex <= 0 || colonIndex >= len(token)-1 {
		return false
	}
	
	// Bot ID part should be numeric
	botIDPart := token[:colonIndex]
	for _, char := range botIDPart {
		if char < '0' || char > '9' {
			return false
		}
	}
	
	// Token part should be at least 10 characters
	tokenPart := token[colonIndex+1:]
	if len(tokenPart) < 10 {
		return false
	}
	
	return true
}

// Validate validates the MessageConfig
func (mc *MessageConfig) Validate() error {
	if mc.ChatID == "" {
		return &ZaloBotError{
			Code:    400,
			Message: "ChatID is required",
			Type:    ErrorTypeValidation,
		}
	}
	
	if mc.MessageType == "" {
		mc.MessageType = MessageTypeText
	}
	
	if !mc.MessageType.IsValid() {
		return &ZaloBotError{
			Code:    400,
			Message: "Invalid message type",
			Type:    ErrorTypeValidation,
		}
	}
	
	if mc.MessageType == MessageTypeText && mc.Text == "" {
		return &ZaloBotError{
			Code:    400,
			Message: "Text is required for text messages",
			Type:    ErrorTypeValidation,
		}
	}
	
	return nil
}

// Validate validates the WebhookConfig
func (wc *WebhookConfig) Validate() error {
	if wc.URL == "" {
		return &ZaloBotError{
			Code:    400,
			Message: "Webhook URL is required",
			Type:    ErrorTypeValidation,
		}
	}
	
	return nil
}

// Validate validates the UpdateConfig
func (uc *UpdateConfig) Validate() error {
	if uc.Limit <= 0 {
		uc.Limit = 100
	}
	
	if uc.Limit > 100 {
		uc.Limit = 100
	}
	
	if uc.Timeout < 0 {
		uc.Timeout = 0
	}
	
	return nil
}

// Validate validates the ImageMessageConfig
func (imc *ImageMessageConfig) Validate() error {
	if imc.ChatID == "" {
		return NewValidationError("ChatID is required for image messages")
	}
	
	if imc.ImageURL == "" {
		return NewValidationError("ImageURL is required for image messages")
	}
	
	// Validate MIME type if provided
	if imc.MimeType != "" {
		validMimeTypes := []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
		isValid := false
		for _, validType := range validMimeTypes {
			if imc.MimeType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return NewValidationError("Invalid MIME type for image message")
		}
	}
	
	return nil
}

// Validate validates the FileMessageConfig
func (fmc *FileMessageConfig) Validate() error {
	if fmc.ChatID == "" {
		return NewValidationError("ChatID is required for file messages")
	}
	
	if fmc.FileURL == "" {
		return NewValidationError("FileURL is required for file messages")
	}
	
	if fmc.FileName == "" {
		return NewValidationError("FileName is required for file messages")
	}
	
	// Validate file size (max 50MB)
	const maxFileSize = 50 * 1024 * 1024 // 50MB
	if fmc.Size > maxFileSize {
		return NewValidationError("File size exceeds maximum limit of 50MB")
	}
	
	return nil
}

// Validate validates the StructuredMessageConfig
func (smc *StructuredMessageConfig) Validate() error {
	if smc.ChatID == "" {
		return NewValidationError("ChatID is required for structured messages")
	}
	
	// Validate structured message content
	if err := smc.StructuredMessage.Validate(); err != nil {
		return err
	}
	
	return nil
}