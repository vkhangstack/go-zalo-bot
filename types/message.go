package types

import (
	"encoding/json"
	"fmt"
	"time"
)

// Message represents a message in Zalo Bot
type Message struct {
	MessageID   string       `json:"message_id"`
	From        *User        `json:"from"`
	Chat        *Chat        `json:"chat"`
	Date        time.Time    `json:"date"`
	Text        string       `json:"text,omitempty"`
	Caption     string       `json:"caption,omitempty"`
	Photo       *Photo       `json:"photo,omitempty"`
	Sticker     *Sticker     `json:"sticker,omitempty"`
	VoiceURL    string       `json:"voice_url,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// rawMessage mirrors Message but leaves fields that need custom decoding as
// raw JSON, so a single Message type can decode both the SDK's own encoding
// (photo/sticker as objects, date as RFC3339) and the real incoming webhook
// wire format (photo/sticker as plain strings, date as a Unix millisecond
// timestamp, a top-level sticker "url"). See https://bot.zapps.me/docs/webhook/
type rawMessage struct {
	MessageID   string          `json:"message_id"`
	From        *User           `json:"from"`
	Chat        *Chat           `json:"chat"`
	Date        json.RawMessage `json:"date"`
	Text        string          `json:"text,omitempty"`
	Caption     string          `json:"caption,omitempty"`
	Photo       json.RawMessage `json:"photo,omitempty"`
	Sticker     json.RawMessage `json:"sticker,omitempty"`
	StickerURL  string          `json:"url,omitempty"`
	VoiceURL    string          `json:"voice_url,omitempty"`
	Attachments []Attachment    `json:"attachments,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler interface
func (m *Message) UnmarshalJSON(data []byte) error {
	var raw rawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	m.MessageID = raw.MessageID
	m.From = raw.From
	m.Chat = raw.Chat
	m.Text = raw.Text
	m.Caption = raw.Caption
	m.VoiceURL = raw.VoiceURL
	m.Attachments = raw.Attachments

	if len(raw.Date) > 0 {
		date, err := parseMessageDate(raw.Date)
		if err != nil {
			return fmt.Errorf("invalid message date: %w", err)
		}
		m.Date = date
	}

	if len(raw.Photo) > 0 {
		photo, err := parsePhoto(raw.Photo)
		if err != nil {
			return fmt.Errorf("invalid photo: %w", err)
		}
		m.Photo = photo
	}

	if len(raw.Sticker) > 0 {
		sticker, err := parseSticker(raw.Sticker)
		if err != nil {
			return fmt.Errorf("invalid sticker: %w", err)
		}
		if raw.StickerURL != "" {
			sticker.URL = raw.StickerURL
		}
		m.Sticker = sticker
	}

	return nil
}

// parseMessageDate decodes a message date, accepting either the Unix
// millisecond timestamp sent by real webhook payloads or the RFC3339 string
// produced by the SDK's own time.Time marshaling.
func parseMessageDate(raw json.RawMessage) (time.Time, error) {
	var millis int64
	if err := json.Unmarshal(raw, &millis); err == nil {
		return time.UnixMilli(millis), nil
	}

	var t time.Time
	if err := json.Unmarshal(raw, &t); err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// parsePhoto decodes a photo, accepting either the plain image URL string
// sent by real webhook payloads or a structured Photo object.
func parsePhoto(raw json.RawMessage) (*Photo, error) {
	var url string
	if err := json.Unmarshal(raw, &url); err == nil {
		return &Photo{URL: url}, nil
	}

	var photo Photo
	if err := json.Unmarshal(raw, &photo); err != nil {
		return nil, err
	}
	return &photo, nil
}

// parseSticker decodes a sticker, accepting either the plain sticker
// id/reference string sent by real webhook payloads or a structured Sticker
// object.
func parseSticker(raw json.RawMessage) (*Sticker, error) {
	var id string
	if err := json.Unmarshal(raw, &id); err == nil {
		return &Sticker{FileID: id}, nil
	}

	var sticker Sticker
	if err := json.Unmarshal(raw, &sticker); err != nil {
		return nil, err
	}
	return &sticker, nil
}

// Attachment represents a file attachment in a message
type Attachment struct {
	Type     AttachmentType `json:"type"`
	URL      string         `json:"url,omitempty"`
	FileID   string         `json:"file_id,omitempty"`
	MimeType string         `json:"mime_type,omitempty"`
	Size     int64          `json:"size,omitempty"`
}

type ChatActionType string

const (
	ChatActionTyping      ChatActionType = "typing"
	ChatActionUploadPhoto ChatActionType = "upload_photo"
)

// AttachmentType represents the type of attachment
type AttachmentType string

const (
	AttachmentTypeImage AttachmentType = "image"
	AttachmentTypeFile  AttachmentType = "file"
	AttachmentTypeVideo AttachmentType = "video"
	AttachmentTypeAudio AttachmentType = "audio"
)

// IsValid validates the attachment type
func (at AttachmentType) IsValid() bool {
	switch at {
	case AttachmentTypeImage, AttachmentTypeFile, AttachmentTypeVideo, AttachmentTypeAudio:
		return true
	default:
		return false
	}
}

// String returns the string representation of AttachmentType
func (at AttachmentType) String() string {
	return string(at)
}

// MarshalJSON implements json.Marshaler interface
func (at AttachmentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(at))
}

// UnmarshalJSON implements json.Unmarshaler interface
func (at *AttachmentType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*at = AttachmentType(s)
	return nil
}

// StructuredMessage represents a structured message with rich content
type StructuredMessage struct {
	Type         StructuredMessageType `json:"type"`
	Elements     []MessageElement      `json:"elements,omitempty"`
	QuickReplies []QuickReply          `json:"quick_replies,omitempty"`
}

// StructuredMessageType represents the type of structured message
type StructuredMessageType string

const (
	StructuredMessageTypeTemplate StructuredMessageType = "template"
	StructuredMessageTypeButton   StructuredMessageType = "button"
	StructuredMessageTypeCarousel StructuredMessageType = "carousel"
)

// IsValid validates the structured message type
func (smt StructuredMessageType) IsValid() bool {
	switch smt {
	case StructuredMessageTypeTemplate, StructuredMessageTypeButton, StructuredMessageTypeCarousel:
		return true
	default:
		return false
	}
}

// String returns the string representation of StructuredMessageType
func (smt StructuredMessageType) String() string {
	return string(smt)
}

// MarshalJSON implements json.Marshaler interface
func (smt StructuredMessageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(smt))
}

// UnmarshalJSON implements json.Unmarshaler interface
func (smt *StructuredMessageType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*smt = StructuredMessageType(s)
	return nil
}

// MessageElement represents an element in a structured message
type MessageElement struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle,omitempty"`
	ImageURL string   `json:"image_url,omitempty"`
	Buttons  []Button `json:"buttons,omitempty"`
}

// Button represents a button in a message element
type Button struct {
	Type    ButtonType `json:"type"`
	Title   string     `json:"title"`
	Payload string     `json:"payload,omitempty"`
	URL     string     `json:"url,omitempty"`
}

// ButtonType represents the type of button
type ButtonType string

const (
	ButtonTypePostback    ButtonType = "postback"
	ButtonTypeWebURL      ButtonType = "web_url"
	ButtonTypePhoneNumber ButtonType = "phone_number"
)

// IsValid validates the button type
func (bt ButtonType) IsValid() bool {
	switch bt {
	case ButtonTypePostback, ButtonTypeWebURL, ButtonTypePhoneNumber:
		return true
	default:
		return false
	}
}

// String returns the string representation of ButtonType
func (bt ButtonType) String() string {
	return string(bt)
}

// MarshalJSON implements json.Marshaler interface
func (bt ButtonType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(bt))
}

// UnmarshalJSON implements json.Unmarshaler interface
func (bt *ButtonType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*bt = ButtonType(s)
	return nil
}

// QuickReply represents a quick reply option
type QuickReply struct {
	ContentType QuickReplyType `json:"content_type"`
	Title       string         `json:"title,omitempty"`
	Payload     string         `json:"payload,omitempty"`
}

// QuickReplyType represents the type of quick reply
type QuickReplyType string

const (
	QuickReplyTypeText     QuickReplyType = "text"
	QuickReplyTypeLocation QuickReplyType = "location"
)

// IsValid validates the quick reply type
func (qrt QuickReplyType) IsValid() bool {
	switch qrt {
	case QuickReplyTypeText, QuickReplyTypeLocation:
		return true
	default:
		return false
	}
}

// String returns the string representation of QuickReplyType
func (qrt QuickReplyType) String() string {
	return string(qrt)
}

// MarshalJSON implements json.Marshaler interface
func (qrt QuickReplyType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(qrt))
}

// UnmarshalJSON implements json.Unmarshaler interface
func (qrt *QuickReplyType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*qrt = QuickReplyType(s)
	return nil
}

// Validate validates the StructuredMessage
func (sm *StructuredMessage) Validate() error {
	if !sm.Type.IsValid() {
		return NewValidationError("Invalid structured message type")
	}

	// Validate elements
	for i, element := range sm.Elements {
		if err := element.Validate(); err != nil {
			return NewValidationError(fmt.Sprintf("Invalid element at index %d: %s", i, err.Error()))
		}
	}

	// Validate quick replies
	for i, quickReply := range sm.QuickReplies {
		if err := quickReply.Validate(); err != nil {
			return NewValidationError(fmt.Sprintf("Invalid quick reply at index %d: %s", i, err.Error()))
		}
	}

	return nil
}

// Validate validates the MessageElement
func (me *MessageElement) Validate() error {
	if me.Title == "" {
		return NewValidationError("Title is required for message element")
	}

	// Validate buttons
	for i, button := range me.Buttons {
		if err := button.Validate(); err != nil {
			return NewValidationError(fmt.Sprintf("Invalid button at index %d: %s", i, err.Error()))
		}
	}

	return nil
}

// Validate validates the Button
func (b *Button) Validate() error {
	if !b.Type.IsValid() {
		return NewValidationError("Invalid button type")
	}

	if b.Title == "" {
		return NewValidationError("Title is required for button")
	}

	switch b.Type {
	case ButtonTypePostback:
		if b.Payload == "" {
			return NewValidationError("Payload is required for postback button")
		}
	case ButtonTypeWebURL:
		if b.URL == "" {
			return NewValidationError("URL is required for web URL button")
		}
	case ButtonTypePhoneNumber:
		if b.Payload == "" {
			return NewValidationError("Payload (phone number) is required for phone number button")
		}
	}

	return nil
}

// Validate validates the QuickReply
func (qr *QuickReply) Validate() error {
	if !qr.ContentType.IsValid() {
		return NewValidationError("Invalid quick reply content type")
	}

	if qr.ContentType == QuickReplyTypeText && qr.Title == "" {
		return NewValidationError("Title is required for text quick reply")
	}

	return nil
}
