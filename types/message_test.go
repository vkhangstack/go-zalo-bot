package types

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAttachmentType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		at   AttachmentType
		want bool
	}{
		{"valid image", AttachmentTypeImage, true},
		{"valid file", AttachmentTypeFile, true},
		{"valid video", AttachmentTypeVideo, true},
		{"valid audio", AttachmentTypeAudio, true},
		{"invalid type", AttachmentType("invalid"), false},
		{"empty type", AttachmentType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.IsValid(); got != tt.want {
				t.Errorf("AttachmentType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentType_String(t *testing.T) {
	tests := []struct {
		name string
		at   AttachmentType
		want string
	}{
		{"image type", AttachmentTypeImage, "image"},
		{"file type", AttachmentTypeFile, "file"},
		{"video type", AttachmentTypeVideo, "video"},
		{"audio type", AttachmentTypeAudio, "audio"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.at.String(); got != tt.want {
				t.Errorf("AttachmentType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAttachmentType_JSON(t *testing.T) {
	tests := []struct {
		name string
		at   AttachmentType
		want string
	}{
		{"marshal image", AttachmentTypeImage, `"image"`},
		{"marshal file", AttachmentTypeFile, `"file"`},
		{"marshal video", AttachmentTypeVideo, `"video"`},
		{"marshal audio", AttachmentTypeAudio, `"audio"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.at)
			if err != nil {
				t.Errorf("AttachmentType.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("AttachmentType.MarshalJSON() = %v, want %v", string(got), tt.want)
			}

			// Test unmarshal
			var at AttachmentType
			if err := json.Unmarshal(got, &at); err != nil {
				t.Errorf("AttachmentType.UnmarshalJSON() error = %v", err)
				return
			}
			if at != tt.at {
				t.Errorf("AttachmentType.UnmarshalJSON() = %v, want %v", at, tt.at)
			}
		})
	}
}

func TestMessage_JSON(t *testing.T) {
	now := time.Now()
	message := Message{
		MessageID: "msg123",
		From: &User{
			ID:   "user123",
			Name: "Test User",
		},
		Chat: &Chat{
			ID:   "chat123",
			Type: ChatTypePrivate,
		},
		Date: now,
		Text: "Hello, World!",
		Attachments: []Attachment{
			{
				Type:     AttachmentTypeImage,
				URL:      "https://example.com/image.jpg",
				FileID:   "file123",
				MimeType: "image/jpeg",
				Size:     1024,
			},
		},
	}

	// Test marshal
	data, err := json.Marshal(message)
	if err != nil {
		t.Errorf("Message.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled Message
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Message.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.MessageID != message.MessageID {
		t.Errorf("Message.MessageID = %v, want %v", unmarshaled.MessageID, message.MessageID)
	}
	if unmarshaled.Text != message.Text {
		t.Errorf("Message.Text = %v, want %v", unmarshaled.Text, message.Text)
	}
	if len(unmarshaled.Attachments) != len(message.Attachments) {
		t.Errorf("Message.Attachments length = %v, want %v", len(unmarshaled.Attachments), len(message.Attachments))
	}
}

func TestStructuredMessageType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		smt  StructuredMessageType
		want bool
	}{
		{"valid template", StructuredMessageTypeTemplate, true},
		{"valid button", StructuredMessageTypeButton, true},
		{"valid carousel", StructuredMessageTypeCarousel, true},
		{"invalid type", StructuredMessageType("invalid"), false},
		{"empty type", StructuredMessageType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.smt.IsValid(); got != tt.want {
				t.Errorf("StructuredMessageType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructuredMessageType_String(t *testing.T) {
	tests := []struct {
		name string
		smt  StructuredMessageType
		want string
	}{
		{"template type", StructuredMessageTypeTemplate, "template"},
		{"button type", StructuredMessageTypeButton, "button"},
		{"carousel type", StructuredMessageTypeCarousel, "carousel"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.smt.String(); got != tt.want {
				t.Errorf("StructuredMessageType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructuredMessageType_JSON(t *testing.T) {
	tests := []struct {
		name string
		smt  StructuredMessageType
		want string
	}{
		{"marshal template", StructuredMessageTypeTemplate, `"template"`},
		{"marshal button", StructuredMessageTypeButton, `"button"`},
		{"marshal carousel", StructuredMessageTypeCarousel, `"carousel"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.smt)
			if err != nil {
				t.Errorf("StructuredMessageType.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("StructuredMessageType.MarshalJSON() = %v, want %v", string(got), tt.want)
			}

			// Test unmarshal
			var smt StructuredMessageType
			if err := json.Unmarshal(got, &smt); err != nil {
				t.Errorf("StructuredMessageType.UnmarshalJSON() error = %v", err)
				return
			}
			if smt != tt.smt {
				t.Errorf("StructuredMessageType.UnmarshalJSON() = %v, want %v", smt, tt.smt)
			}
		})
	}
}

func TestButtonType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		bt   ButtonType
		want bool
	}{
		{"valid postback", ButtonTypePostback, true},
		{"valid web_url", ButtonTypeWebURL, true},
		{"valid phone_number", ButtonTypePhoneNumber, true},
		{"invalid type", ButtonType("invalid"), false},
		{"empty type", ButtonType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bt.IsValid(); got != tt.want {
				t.Errorf("ButtonType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButtonType_String(t *testing.T) {
	tests := []struct {
		name string
		bt   ButtonType
		want string
	}{
		{"postback type", ButtonTypePostback, "postback"},
		{"web_url type", ButtonTypeWebURL, "web_url"},
		{"phone_number type", ButtonTypePhoneNumber, "phone_number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.bt.String(); got != tt.want {
				t.Errorf("ButtonType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestButtonType_JSON(t *testing.T) {
	tests := []struct {
		name string
		bt   ButtonType
		want string
	}{
		{"marshal postback", ButtonTypePostback, `"postback"`},
		{"marshal web_url", ButtonTypeWebURL, `"web_url"`},
		{"marshal phone_number", ButtonTypePhoneNumber, `"phone_number"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.bt)
			if err != nil {
				t.Errorf("ButtonType.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("ButtonType.MarshalJSON() = %v, want %v", string(got), tt.want)
			}

			// Test unmarshal
			var bt ButtonType
			if err := json.Unmarshal(got, &bt); err != nil {
				t.Errorf("ButtonType.UnmarshalJSON() error = %v", err)
				return
			}
			if bt != tt.bt {
				t.Errorf("ButtonType.UnmarshalJSON() = %v, want %v", bt, tt.bt)
			}
		})
	}
}

func TestQuickReplyType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		qrt  QuickReplyType
		want bool
	}{
		{"valid text", QuickReplyTypeText, true},
		{"valid location", QuickReplyTypeLocation, true},
		{"invalid type", QuickReplyType("invalid"), false},
		{"empty type", QuickReplyType(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.qrt.IsValid(); got != tt.want {
				t.Errorf("QuickReplyType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuickReplyType_String(t *testing.T) {
	tests := []struct {
		name string
		qrt  QuickReplyType
		want string
	}{
		{"text type", QuickReplyTypeText, "text"},
		{"location type", QuickReplyTypeLocation, "location"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.qrt.String(); got != tt.want {
				t.Errorf("QuickReplyType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuickReplyType_JSON(t *testing.T) {
	tests := []struct {
		name string
		qrt  QuickReplyType
		want string
	}{
		{"marshal text", QuickReplyTypeText, `"text"`},
		{"marshal location", QuickReplyTypeLocation, `"location"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.qrt)
			if err != nil {
				t.Errorf("QuickReplyType.MarshalJSON() error = %v", err)
				return
			}
			if string(got) != tt.want {
				t.Errorf("QuickReplyType.MarshalJSON() = %v, want %v", string(got), tt.want)
			}

			// Test unmarshal
			var qrt QuickReplyType
			if err := json.Unmarshal(got, &qrt); err != nil {
				t.Errorf("QuickReplyType.UnmarshalJSON() error = %v", err)
				return
			}
			if qrt != tt.qrt {
				t.Errorf("QuickReplyType.UnmarshalJSON() = %v, want %v", qrt, tt.qrt)
			}
		})
	}
}

func TestStructuredMessage_JSON(t *testing.T) {
	structuredMsg := StructuredMessage{
		Type: StructuredMessageTypeTemplate,
		Elements: []MessageElement{
			{
				Title:    "Test Element",
				Subtitle: "Test Subtitle",
				ImageURL: "https://example.com/image.jpg",
				Buttons: []Button{
					{
						Type:    ButtonTypePostback,
						Title:   "Click Me",
						Payload: "test_payload",
					},
					{
						Type:  ButtonTypeWebURL,
						Title: "Visit Website",
						URL:   "https://example.com",
					},
				},
			},
		},
		QuickReplies: []QuickReply{
			{
				ContentType: QuickReplyTypeText,
				Title:       "Quick Reply 1",
				Payload:     "qr_payload_1",
			},
			{
				ContentType: QuickReplyTypeLocation,
				Title:       "Share Location",
			},
		},
	}

	// Test marshal
	data, err := json.Marshal(structuredMsg)
	if err != nil {
		t.Errorf("StructuredMessage.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled StructuredMessage
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("StructuredMessage.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.Type != structuredMsg.Type {
		t.Errorf("StructuredMessage.Type = %v, want %v", unmarshaled.Type, structuredMsg.Type)
	}
	if len(unmarshaled.Elements) != len(structuredMsg.Elements) {
		t.Errorf("StructuredMessage.Elements length = %v, want %v", len(unmarshaled.Elements), len(structuredMsg.Elements))
	}
	if len(unmarshaled.QuickReplies) != len(structuredMsg.QuickReplies) {
		t.Errorf("StructuredMessage.QuickReplies length = %v, want %v", len(unmarshaled.QuickReplies), len(structuredMsg.QuickReplies))
	}

	// Verify element details
	if len(unmarshaled.Elements) > 0 {
		elem := unmarshaled.Elements[0]
		if elem.Title != "Test Element" {
			t.Errorf("MessageElement.Title = %v, want %v", elem.Title, "Test Element")
		}
		if elem.Subtitle != "Test Subtitle" {
			t.Errorf("MessageElement.Subtitle = %v, want %v", elem.Subtitle, "Test Subtitle")
		}
		if len(elem.Buttons) != 2 {
			t.Errorf("MessageElement.Buttons length = %v, want %v", len(elem.Buttons), 2)
		}
	}

	// Verify quick reply details
	if len(unmarshaled.QuickReplies) > 0 {
		qr := unmarshaled.QuickReplies[0]
		if qr.ContentType != QuickReplyTypeText {
			t.Errorf("QuickReply.ContentType = %v, want %v", qr.ContentType, QuickReplyTypeText)
		}
		if qr.Title != "Quick Reply 1" {
			t.Errorf("QuickReply.Title = %v, want %v", qr.Title, "Quick Reply 1")
		}
	}
}

func TestButton_JSON(t *testing.T) {
	button := Button{
		Type:    ButtonTypePostback,
		Title:   "Test Button",
		Payload: "test_payload",
		URL:     "https://example.com",
	}

	// Test marshal
	data, err := json.Marshal(button)
	if err != nil {
		t.Errorf("Button.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled Button
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Button.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.Type != button.Type {
		t.Errorf("Button.Type = %v, want %v", unmarshaled.Type, button.Type)
	}
	if unmarshaled.Title != button.Title {
		t.Errorf("Button.Title = %v, want %v", unmarshaled.Title, button.Title)
	}
	if unmarshaled.Payload != button.Payload {
		t.Errorf("Button.Payload = %v, want %v", unmarshaled.Payload, button.Payload)
	}
	if unmarshaled.URL != button.URL {
		t.Errorf("Button.URL = %v, want %v", unmarshaled.URL, button.URL)
	}
}

func TestMessageElement_JSON(t *testing.T) {
	element := MessageElement{
		Title:    "Test Element",
		Subtitle: "Test Subtitle",
		ImageURL: "https://example.com/image.jpg",
		Buttons: []Button{
			{
				Type:    ButtonTypeWebURL,
				Title:   "Visit",
				URL:     "https://example.com",
			},
		},
	}

	// Test marshal
	data, err := json.Marshal(element)
	if err != nil {
		t.Errorf("MessageElement.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled MessageElement
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("MessageElement.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.Title != element.Title {
		t.Errorf("MessageElement.Title = %v, want %v", unmarshaled.Title, element.Title)
	}
	if unmarshaled.Subtitle != element.Subtitle {
		t.Errorf("MessageElement.Subtitle = %v, want %v", unmarshaled.Subtitle, element.Subtitle)
	}
	if unmarshaled.ImageURL != element.ImageURL {
		t.Errorf("MessageElement.ImageURL = %v, want %v", unmarshaled.ImageURL, element.ImageURL)
	}
	if len(unmarshaled.Buttons) != len(element.Buttons) {
		t.Errorf("MessageElement.Buttons length = %v, want %v", len(unmarshaled.Buttons), len(element.Buttons))
	}
}

func TestQuickReply_JSON(t *testing.T) {
	quickReply := QuickReply{
		ContentType: QuickReplyTypeText,
		Title:       "Quick Reply",
		Payload:     "qr_payload",
	}

	// Test marshal
	data, err := json.Marshal(quickReply)
	if err != nil {
		t.Errorf("QuickReply.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled QuickReply
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("QuickReply.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.ContentType != quickReply.ContentType {
		t.Errorf("QuickReply.ContentType = %v, want %v", unmarshaled.ContentType, quickReply.ContentType)
	}
	if unmarshaled.Title != quickReply.Title {
		t.Errorf("QuickReply.Title = %v, want %v", unmarshaled.Title, quickReply.Title)
	}
	if unmarshaled.Payload != quickReply.Payload {
		t.Errorf("QuickReply.Payload = %v, want %v", unmarshaled.Payload, quickReply.Payload)
	}
}

func TestStructuredMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		sm      *StructuredMessage
		wantErr bool
	}{
		{
			name: "valid structured message",
			sm: &StructuredMessage{
				Type: StructuredMessageTypeTemplate,
				Elements: []MessageElement{
					{
						Title: "Test Element",
						Buttons: []Button{
							{
								Type:    ButtonTypePostback,
								Title:   "Click Me",
								Payload: "test_payload",
							},
						},
					},
				},
				QuickReplies: []QuickReply{
					{
						ContentType: QuickReplyTypeText,
						Title:       "Quick Reply",
						Payload:     "qr_payload",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid structured message type",
			sm: &StructuredMessage{
				Type: StructuredMessageType("invalid"),
			},
			wantErr: true,
		},
		{
			name: "invalid element",
			sm: &StructuredMessage{
				Type: StructuredMessageTypeTemplate,
				Elements: []MessageElement{
					{
						Title: "", // Empty title
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid quick reply",
			sm: &StructuredMessage{
				Type: StructuredMessageTypeTemplate,
				QuickReplies: []QuickReply{
					{
						ContentType: QuickReplyType("invalid"),
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sm.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("StructuredMessage.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMessageElement_Validate(t *testing.T) {
	tests := []struct {
		name    string
		me      *MessageElement
		wantErr bool
	}{
		{
			name: "valid message element",
			me: &MessageElement{
				Title:    "Test Element",
				Subtitle: "Test Subtitle",
				ImageURL: "https://example.com/image.jpg",
				Buttons: []Button{
					{
						Type:    ButtonTypeWebURL,
						Title:   "Visit",
						URL:     "https://example.com",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing title",
			me: &MessageElement{
				Title: "",
			},
			wantErr: true,
		},
		{
			name: "invalid button",
			me: &MessageElement{
				Title: "Test Element",
				Buttons: []Button{
					{
						Type:  ButtonTypePostback,
						Title: "Click Me",
						// Missing required Payload
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.me.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("MessageElement.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestButton_Validate(t *testing.T) {
	tests := []struct {
		name    string
		button  *Button
		wantErr bool
	}{
		{
			name: "valid postback button",
			button: &Button{
				Type:    ButtonTypePostback,
				Title:   "Click Me",
				Payload: "test_payload",
			},
			wantErr: false,
		},
		{
			name: "valid web URL button",
			button: &Button{
				Type:  ButtonTypeWebURL,
				Title: "Visit",
				URL:   "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "valid phone number button",
			button: &Button{
				Type:    ButtonTypePhoneNumber,
				Title:   "Call",
				Payload: "+1234567890",
			},
			wantErr: false,
		},
		{
			name: "invalid button type",
			button: &Button{
				Type:  ButtonType("invalid"),
				Title: "Test",
			},
			wantErr: true,
		},
		{
			name: "missing title",
			button: &Button{
				Type:    ButtonTypePostback,
				Payload: "test_payload",
			},
			wantErr: true,
		},
		{
			name: "postback button missing payload",
			button: &Button{
				Type:  ButtonTypePostback,
				Title: "Click Me",
			},
			wantErr: true,
		},
		{
			name: "web URL button missing URL",
			button: &Button{
				Type:  ButtonTypeWebURL,
				Title: "Visit",
			},
			wantErr: true,
		},
		{
			name: "phone number button missing payload",
			button: &Button{
				Type:  ButtonTypePhoneNumber,
				Title: "Call",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.button.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Button.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuickReply_Validate(t *testing.T) {
	tests := []struct {
		name    string
		qr      *QuickReply
		wantErr bool
	}{
		{
			name: "valid text quick reply",
			qr: &QuickReply{
				ContentType: QuickReplyTypeText,
				Title:       "Quick Reply",
				Payload:     "qr_payload",
			},
			wantErr: false,
		},
		{
			name: "valid location quick reply",
			qr: &QuickReply{
				ContentType: QuickReplyTypeLocation,
			},
			wantErr: false,
		},
		{
			name: "invalid content type",
			qr: &QuickReply{
				ContentType: QuickReplyType("invalid"),
			},
			wantErr: true,
		},
		{
			name: "text quick reply missing title",
			qr: &QuickReply{
				ContentType: QuickReplyTypeText,
				Payload:     "qr_payload",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.qr.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("QuickReply.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAttachment_JSON(t *testing.T) {
	attachment := Attachment{
		Type:     AttachmentTypeImage,
		URL:      "https://example.com/image.jpg",
		FileID:   "file123",
		MimeType: "image/jpeg",
		Size:     1024,
	}

	// Test marshal
	data, err := json.Marshal(attachment)
	if err != nil {
		t.Errorf("Attachment.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled Attachment
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Attachment.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.Type != attachment.Type {
		t.Errorf("Attachment.Type = %v, want %v", unmarshaled.Type, attachment.Type)
	}
	if unmarshaled.URL != attachment.URL {
		t.Errorf("Attachment.URL = %v, want %v", unmarshaled.URL, attachment.URL)
	}
	if unmarshaled.FileID != attachment.FileID {
		t.Errorf("Attachment.FileID = %v, want %v", unmarshaled.FileID, attachment.FileID)
	}
	if unmarshaled.MimeType != attachment.MimeType {
		t.Errorf("Attachment.MimeType = %v, want %v", unmarshaled.MimeType, attachment.MimeType)
	}
	if unmarshaled.Size != attachment.Size {
		t.Errorf("Attachment.Size = %v, want %v", unmarshaled.Size, attachment.Size)
	}
}
