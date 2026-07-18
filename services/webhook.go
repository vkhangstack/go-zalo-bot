package services

import (
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

// ValidateSecretToken compares an incoming X-Bot-Api-Secret-Token header value
// against the configured secret token, as instructed by the Zalo webhook docs:
// https://bot.zapps.me/docs/webhook/
func (s *WebhookService) ValidateSecretToken(token string) error {
	if s.secretToken == "" {
		return fmt.Errorf("webhook secret token is not configured")
	}
	if token != s.secretToken {
		return fmt.Errorf("invalid secret token")
	}
	return nil
}

// RejectInvalidRequest creates an error for rejecting invalid webhook requests
func (s *WebhookService) RejectInvalidRequest(reason string) error {
	return utils.RejectInvalidWebhookRequest(reason)
}

// ParseUpdate parses a webhook request body into an Update, following the
// envelope Zalo sends per https://bot.zapps.me/docs/webhook/:
// {"ok":true,"result":{"event_name":...,"message":{...}}}.
func (s *WebhookService) ParseUpdate(payload []byte) (*types.Update, error) {
	if len(payload) == 0 {
		return nil, fmt.Errorf("empty webhook payload")
	}

	webhookPayload, err := types.ParseWebhookPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webhook payload: %w", err)
	}

	if webhookPayload.Result.EventName == "" {
		return nil, fmt.Errorf("webhook payload result is empty or missing event_name")
	}

	return &types.Update{
		EventName: webhookPayload.Result.EventName,
		Message:   webhookPayload.Result.Message,
	}, nil
}

// ProcessWebhook validates the request's secret token and parses its payload
// into an Update. secretToken should be the value of the
// X-Bot-Api-Secret-Token header from the incoming request.
func (s *WebhookService) ProcessWebhook(payload []byte, secretToken string) (*types.Update, error) {
	if err := s.ValidateSecretToken(secretToken); err != nil {
		return nil, s.RejectInvalidRequest(fmt.Sprintf("secret token validation failed: %v", err))
	}

	update, err := s.ParseUpdate(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webhook update: %w", err)
	}

	return update, nil
}

// HandleWebhookEvent processes a webhook request and returns a coarse event
// type ("message" or "unknown") along with the corresponding event data.
func (s *WebhookService) HandleWebhookEvent(payload []byte, secretToken string) (string, interface{}, error) {
	update, err := s.ProcessWebhook(payload, secretToken)
	if err != nil {
		return "", nil, err
	}

	if update.Message != nil {
		return "message", update.Message, nil
	}

	return "unknown", update, nil
}
