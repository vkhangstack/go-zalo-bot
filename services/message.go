package services

import (
	"context"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/vkhangstack/go-zalo-bot/auth"
	"github.com/vkhangstack/go-zalo-bot/types"
)

// MessageService handles message-related operations
type MessageService struct {
	*BaseService
}

// NewMessageService creates a new message service
func NewMessageService(authService *auth.AuthService, client *http.Client, config *types.Config) *MessageService {
	return &MessageService{
		BaseService: NewBaseService(authService, client, config),
	}
}

// Send sends a message with proper payload validation and bot token URL construction
// Supports Unicode text including Vietnamese characters
func (s *MessageService) Send(ctx context.Context, config types.MessageConfig) (*types.Message, error) {
	// Validate message config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Validate text content for Unicode support (including Vietnamese)
	if config.Text != "" {
		if !utf8.ValidString(config.Text) {
			return nil, types.NewValidationError("Text contains invalid UTF-8 characters")
		}
	}

	// Validate recipient ID format
	if err := validateRecipientID(config.ChatID); err != nil {
		return nil, err
	}

	// Prepare request payload
	payload := map[string]interface{}{
		"chat_id": config.ChatID,
	}

	// Add text if present
	if config.Text != "" {
		payload["text"] = config.Text
	}

	// Add attachments if present
	if len(config.Attachments) > 0 {
		attachments := make([]map[string]interface{}, len(config.Attachments))
		for i, att := range config.Attachments {
			if !att.Type.IsValid() {
				return nil, types.NewValidationError(fmt.Sprintf("Invalid attachment type at index %d", i))
			}
			attachments[i] = map[string]interface{}{
				"type":      att.Type,
				"url":       att.URL,
				"file_id":   att.FileID,
				"mime_type": att.MimeType,
				"size":      att.Size,
			}
		}
		payload["attachments"] = attachments
	}

	// Execute request with bot token URL construction
	apiReq := &APIRequest{
		Method:    http.MethodPost,
		APIMethod: "sendMessage",
		Body:      payload,
	}

	resp, err := s.DoRequest(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var message types.Message
	if err := parseResult(resp.Result, &message); err != nil {
		return nil, types.NewAPIError(0, "Failed to parse message response", err.Error())
	}

	return &message, nil
}

// SendImage sends an image message
func (s *MessageService) SendImage(ctx context.Context, config types.ImageMessageConfig) (*types.Message, error) {
	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Validate recipient ID
	if err := validateRecipientID(config.ChatID); err != nil {
		return nil, err
	}

	// Prepare message config with image attachment
	messageConfig := types.MessageConfig{
		ChatID:      config.ChatID,
		Text:        config.Caption,
		MessageType: types.MessageTypeImage,
		Attachments: []types.Attachment{
			{
				Type:     types.AttachmentTypeImage,
				URL:      config.ImageURL,
				MimeType: config.MimeType,
			},
		},
	}

	return s.Send(ctx, messageConfig)
}

// SendFile sends a file message
func (s *MessageService) SendFile(ctx context.Context, config types.FileMessageConfig) (*types.Message, error) {
	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Validate recipient ID
	if err := validateRecipientID(config.ChatID); err != nil {
		return nil, err
	}

	// Prepare message config with file attachment
	messageConfig := types.MessageConfig{
		ChatID:      config.ChatID,
		MessageType: types.MessageTypeFile,
		Attachments: []types.Attachment{
			{
				Type:     types.AttachmentTypeFile,
				URL:      config.FileURL,
				MimeType: config.MimeType,
				Size:     config.Size,
			},
		},
	}

	return s.Send(ctx, messageConfig)
}

// SendVideo sends a video message
func (s *MessageService) SendVideo(ctx context.Context, chatID, videoURL, mimeType string) (*types.Message, error) {
	// Validate recipient ID
	if err := validateRecipientID(chatID); err != nil {
		return nil, err
	}

	if videoURL == "" {
		return nil, types.NewValidationError("Video URL is required")
	}

	// Validate MIME type if provided
	if mimeType != "" {
		validMimeTypes := []string{"video/mp4", "video/mpeg", "video/quicktime", "video/webm"}
		isValid := false
		for _, validType := range validMimeTypes {
			if mimeType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, types.NewValidationError("Invalid MIME type for video message")
		}
	}

	// Prepare message config with video attachment
	messageConfig := types.MessageConfig{
		ChatID:      chatID,
		MessageType: types.MessageTypeFile,
		Attachments: []types.Attachment{
			{
				Type:     types.AttachmentTypeVideo,
				URL:      videoURL,
				MimeType: mimeType,
			},
		},
	}

	return s.Send(ctx, messageConfig)
}

// SendAudio sends an audio message
func (s *MessageService) SendAudio(ctx context.Context, chatID, audioURL, mimeType string) (*types.Message, error) {
	// Validate recipient ID
	if err := validateRecipientID(chatID); err != nil {
		return nil, err
	}

	if audioURL == "" {
		return nil, types.NewValidationError("Audio URL is required")
	}

	// Validate MIME type if provided
	if mimeType != "" {
		validMimeTypes := []string{"audio/mpeg", "audio/mp4", "audio/ogg", "audio/wav"}
		isValid := false
		for _, validType := range validMimeTypes {
			if mimeType == validType {
				isValid = true
				break
			}
		}
		if !isValid {
			return nil, types.NewValidationError("Invalid MIME type for audio message")
		}
	}

	// Prepare message config with audio attachment
	messageConfig := types.MessageConfig{
		ChatID:      chatID,
		MessageType: types.MessageTypeFile,
		Attachments: []types.Attachment{
			{
				Type:     types.AttachmentTypeAudio,
				URL:      audioURL,
				MimeType: mimeType,
			},
		},
	}

	return s.Send(ctx, messageConfig)
}

// SendTemplate sends a structured message with buttons, quick replies, and interactive elements
func (s *MessageService) SendTemplate(ctx context.Context, config types.StructuredMessageConfig) (*types.Message, error) {
	// Validate config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Validate recipient ID
	if err := validateRecipientID(config.ChatID); err != nil {
		return nil, err
	}

	// Prepare request payload
	payload := map[string]interface{}{
		"chat_id":            config.ChatID,
		"structured_message": config.StructuredMessage,
	}

	// Execute request with bot token URL construction
	apiReq := &APIRequest{
		Method:    http.MethodPost,
		APIMethod: "sendTemplate",
		Body:      payload,
	}

	resp, err := s.DoRequest(ctx, apiReq)
	if err != nil {
		return nil, err
	}

	// Parse response
	var message types.Message
	if err := parseResult(resp.Result, &message); err != nil {
		return nil, types.NewAPIError(0, "Failed to parse message response", err.Error())
	}

	return &message, nil
}

// SendStructuredMessage sends a structured message (alias for SendTemplate)
func (s *MessageService) SendStructuredMessage(ctx context.Context, config types.StructuredMessageConfig) (*types.Message, error) {
	return s.SendTemplate(ctx, config)
}

// SendImageMessage sends an image message (alias for SendImage)
func (s *MessageService) SendImageMessage(ctx context.Context, config types.ImageMessageConfig) (*types.Message, error) {
	return s.SendImage(ctx, config)
}

// SendFileMessage sends a file message (alias for SendFile)
func (s *MessageService) SendFileMessage(ctx context.Context, config types.FileMessageConfig) (*types.Message, error) {
	return s.SendFile(ctx, config)
}

// validateRecipientID validates the recipient ID format
func validateRecipientID(recipientID string) error {
	if recipientID == "" {
		return types.NewValidationError("Recipient ID is required")
	}

	// Recipient ID should be non-empty and reasonable length
	if len(recipientID) < 1 || len(recipientID) > 100 {
		return types.NewValidationError("Invalid recipient ID length")
	}

	return nil
}
