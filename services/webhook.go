package services

import (
	"encoding/json"
	"fmt"

	"github.com/vkhangstack/go-zalo-bot/types"
	"github.com/vkhangstack/go-zalo-bot/utils"
)

// WebhookService handles webhook-related operations
type WebhookService struct {
	*BaseService
	secretToken string
}

// NewWebhookService creates a new webhook service
func NewWebhookService(base *BaseService, secretToken string) *WebhookService {
	return &WebhookService{
		BaseService: base,
		secretToken: secretToken,
	}
}

// SetSecretToken sets the webhook secret token
func (s *WebhookService) SetSecretToken(token string) {
	s.secretToken = token
}

// GetSecretToken returns the webhook secret token
func (s *WebhookService) GetSecretToken() string {
	return s.secretToken
}

// ValidateSignature validates the webhook signature
func (s *WebhookService) ValidateSignature(payload []byte, signature string) error {
	return utils.ValidateWebhookSignature(payload, signature, s.secretToken)
}

// RejectInvalidRequest creates an error for rejecting invalid webhook requests
func (s *WebhookService) RejectInvalidRequest(reason string) error {
	return utils.RejectInvalidWebhookRequest(reason)
}

// ParseUpdate parses a webhook payload into an Update structure
// Supports multiple event types including text messages, attachments, postback events, and user actions
func (s *WebhookService) ParseUpdate(payload []byte) (*types.Update, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("empty webhook payload")
	}

	// First, try to parse as a standard Update
	var update types.Update
	updateErr := json.Unmarshal(payload, &update)

	// Check if it's a valid Update (has update_id field)
	if updateErr == nil && update.UpdateID != 0 {
		// Successfully parsed as Update
		return &update, nil
	}

	// If that fails, try to parse as WebhookEvent
	var webhookEvent types.WebhookEvent
	if err := json.Unmarshal(payload, &webhookEvent); err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	// Convert WebhookEvent to Update based on event type
	return s.convertWebhookEventToUpdate(&webhookEvent)
}

// convertWebhookEventToUpdate converts a WebhookEvent to an Update
func (s *WebhookService) convertWebhookEventToUpdate(event *types.WebhookEvent) (*types.Update, error) {
	// Use timestamp as update ID, ensure it's non-zero
	updateID := int(event.Timestamp)
	if updateID == 0 {
		updateID = 1 // Default to 1 if timestamp is 0
	}

	update := &types.Update{
		UpdateID: updateID,
	}

	// Parse event based on event name
	switch event.EventName {
	case "message", "text_message", "message_received":
		// Parse as message event
		msgEvent, err := event.ParseMessageEvent()
		if err != nil {
			return nil, fmt.Errorf("failed to parse message event: %w", err)
		}
		update.Message = s.convertMessageEventToMessage(msgEvent)

	case "postback", "button_click":
		// Parse as postback event
		var postback types.PostbackEvent
		if err := json.Unmarshal(event.Data, &postback); err != nil {
			return nil, fmt.Errorf("failed to parse postback event: %w", err)
		}
		update.PostbackEvent = &postback

	case "user_action", "user_join", "user_leave", "user_block":
		// Parse as user action event
		actionEvent, err := event.ParseUserActionEvent()
		if err != nil {
			return nil, fmt.Errorf("failed to parse user action event: %w", err)
		}
		update.UserAction = s.convertUserActionEventToUserAction(actionEvent)

	default:
		return nil, fmt.Errorf("unsupported event type: %s", event.EventName)
	}

	return update, nil
}

// convertMessageEventToMessage converts a MessageEvent to a Message
func (s *WebhookService) convertMessageEventToMessage(event *types.MessageEvent) *types.Message {
	return &types.Message{
		MessageID: event.MessageID,
		From: &types.User{
			ID: event.UserID,
		},
		Text:        event.Text,
		Attachments: event.Attachments,
	}
}

// convertUserActionEventToUserAction converts a UserActionEvent to a UserAction
func (s *WebhookService) convertUserActionEventToUserAction(event *types.UserActionEvent) *types.UserAction {
	return &types.UserAction{
		Type:   types.UserActionType(event.Action),
		UserID: event.UserID,
	}
}

// ProcessWebhook processes a webhook request with signature validation
// Returns the parsed Update if validation succeeds, or an error if validation fails
func (s *WebhookService) ProcessWebhook(payload []byte, signature string) (*types.Update, error) {
	// Validate signature
	if err := s.ValidateSignature(payload, signature); err != nil {
		return nil, s.RejectInvalidRequest(fmt.Sprintf("signature validation failed: %v", err))
	}

	// Parse update
	update, err := s.ParseUpdate(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webhook update: %w", err)
	}

	return update, nil
}

// HandleWebhookEvent handles different types of webhook events
// Returns the event type and parsed data
func (s *WebhookService) HandleWebhookEvent(payload []byte, signature string) (string, interface{}, error) {
	// Process webhook with validation
	update, err := s.ProcessWebhook(payload, signature)
	if err != nil {
		return "", nil, err
	}

	// Determine event type and return appropriate data
	if update.Message != nil {
		return "message", update.Message, nil
	}

	if update.PostbackEvent != nil {
		return "postback", update.PostbackEvent, nil
	}

	if update.UserAction != nil {
		return "user_action", update.UserAction, nil
	}

	return "unknown", update, nil
}
