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
	Photo       *Photo       `json:"photo,omitempty"`
	Sticker     *Sticker     `json:"sticker,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a file attachment in a message
type Attachment struct {
	Type     AttachmentType `json:"type"`
	URL      string         `json:"url,omitempty"`
	FileID   string         `json:"file_id,omitempty"`
	MimeType string         `json:"mime_type,omitempty"`
	Size     int64          `json:"size,omitempty"`
}

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
