package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test-secret-token"
	payload := []byte(`{"update_id":1,"message":{"text":"hello"}}`)

	// Compute valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSignature := hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name      string
		payload   []byte
		signature string
		secret    string
		want      bool
	}{
		{
			name:      "valid signature",
			payload:   payload,
			signature: validSignature,
			secret:    secret,
			want:      true,
		},
		{
			name:      "invalid signature",
			payload:   payload,
			signature: "invalid-signature",
			secret:    secret,
			want:      false,
		},
		{
			name:      "empty signature",
			payload:   payload,
			signature: "",
			secret:    secret,
			want:      false,
		},
		{
			name:      "empty secret",
			payload:   payload,
			signature: validSignature,
			secret:    "",
			want:      false,
		},
		{
			name:      "empty payload",
			payload:   []byte{},
			signature: validSignature,
			secret:    secret,
			want:      false,
		},
		{
			name:      "whitespace signature",
			payload:   payload,
			signature: "   ",
			secret:    secret,
			want:      false,
		},
		{
			name:      "whitespace secret",
			payload:   payload,
			signature: validSignature,
			secret:    "   ",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VerifyWebhookSignature(tt.payload, tt.signature, tt.secret)
			if got != tt.want {
				t.Errorf("VerifyWebhookSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateWebhookSignature(t *testing.T) {
	secret := "test-secret"
	payload := []byte(`{"test":"data"}`)

	// Compute valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSignature := hex.EncodeToString(mac.Sum(nil))

	tests := []struct {
		name      string
		payload   []byte
		signature string
		secret    string
		wantErr   error
	}{
		{
			name:      "valid signature",
			payload:   payload,
			signature: validSignature,
			secret:    secret,
			wantErr:   nil,
		},
		{
			name:      "invalid signature",
			payload:   payload,
			signature: "wrong-signature",
			secret:    secret,
			wantErr:   ErrInvalidSignature,
		},
		{
			name:      "empty payload",
			payload:   []byte{},
			signature: validSignature,
			secret:    secret,
			wantErr:   ErrEmptyPayload,
		},
		{
			name:      "empty signature",
			payload:   payload,
			signature: "",
			secret:    secret,
			wantErr:   ErrEmptySignature,
		},
		{
			name:      "empty secret",
			payload:   payload,
			signature: validSignature,
			secret:    "",
			wantErr:   ErrEmptySecret,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWebhookSignature(tt.payload, tt.signature, tt.secret)
			if err != tt.wantErr {
				t.Errorf("ValidateWebhookSignature() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRejectInvalidWebhookRequest(t *testing.T) {
	tests := []struct {
		name   string
		reason string
	}{
		{
			name:   "with reason",
			reason: "invalid signature",
		},
		{
			name:   "empty reason",
			reason: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RejectInvalidWebhookRequest(tt.reason)
			if err == nil {
				t.Error("RejectInvalidWebhookRequest() should return an error")
			}
		})
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "valid struct",
			input:   map[string]string{"key": "value"},
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ToJSON(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFromJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		target  interface{}
		wantErr bool
	}{
		{
			name:    "valid json",
			data:    []byte(`{"key":"value"}`),
			target:  &map[string]string{},
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    []byte(`{invalid}`),
			target:  &map[string]string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FromJSON(tt.data, tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStringPtr(t *testing.T) {
	s := "test"
	ptr := StringPtr(s)
	if ptr == nil || *ptr != s {
		t.Errorf("StringPtr() = %v, want %v", ptr, &s)
	}
}

func TestIntPtr(t *testing.T) {
	i := 42
	ptr := IntPtr(i)
	if ptr == nil || *ptr != i {
		t.Errorf("IntPtr() = %v, want %v", ptr, &i)
	}
}

func TestBoolPtr(t *testing.T) {
	b := true
	ptr := BoolPtr(b)
	if ptr == nil || *ptr != b {
		t.Errorf("BoolPtr() = %v, want %v", ptr, &b)
	}
}

func TestValidateUserID(t *testing.T) {
	tests := []struct {
		name    string
		userID  string
		wantErr bool
	}{
		{
			name:    "valid user ID",
			userID:  "user123",
			wantErr: false,
		},
		{
			name:    "valid with underscores",
			userID:  "user_123_test",
			wantErr: false,
		},
		{
			name:    "empty user ID",
			userID:  "",
			wantErr: true,
		},
		{
			name:    "user ID with special characters",
			userID:  "user@123",
			wantErr: true,
		},
		{
			name:    "user ID too long",
			userID:  "a123456789012345678901234567890123456789012345678901234567890123456",
			wantErr: true,
		},
		{
			name:    "user ID with spaces",
			userID:  "user 123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRecipientID(t *testing.T) {
	tests := []struct {
		name        string
		recipientID string
		wantErr     bool
	}{
		{
			name:        "valid recipient ID",
			recipientID: "recipient123",
			wantErr:     false,
		},
		{
			name:        "empty recipient ID",
			recipientID: "",
			wantErr:     true,
		},
		{
			name:        "recipient ID with special characters",
			recipientID: "recipient-123",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRecipientID(tt.recipientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRecipientID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMessageContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid message",
			content: "Hello, world!",
			wantErr: false,
		},
		{
			name:    "valid Vietnamese message",
			content: "Xin chào, thế giới!",
			wantErr: false,
		},
		{
			name:    "empty message",
			content: "",
			wantErr: true,
		},
		{
			name:    "message with emojis",
			content: "Hello 👋 World 🌍",
			wantErr: false,
		},
		{
			name:    "very long message",
			content: string(make([]byte, MaxMessageLength+1)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMessageContent(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMessageContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestContainsVietnamese(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "Vietnamese text",
			text: "Xin chào",
			want: true,
		},
		{
			name: "Vietnamese with diacritics",
			text: "Tiếng Việt",
			want: true,
		},
		{
			name: "English text",
			text: "Hello World",
			want: false,
		},
		{
			name: "Mixed text",
			text: "Hello Việt Nam",
			want: true,
		},
		{
			name: "Vietnamese Đ character",
			text: "Đồng",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsVietnamese(tt.text)
			if got != tt.want {
				t.Errorf("ContainsVietnamese() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateUnicodeSupport(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "valid UTF-8",
			text:    "Hello World",
			wantErr: false,
		},
		{
			name:    "valid Vietnamese UTF-8",
			text:    "Xin chào thế giới",
			wantErr: false,
		},
		{
			name:    "valid emojis",
			text:    "Hello 👋 🌍",
			wantErr: false,
		},
		{
			name:    "text with replacement character",
			text:    "Hello \uFFFD World",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUnicodeSupport(tt.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUnicodeSupport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateImageMimeType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		wantErr  bool
	}{
		{
			name:     "valid JPEG",
			mimeType: "image/jpeg",
			wantErr:  false,
		},
		{
			name:     "valid PNG",
			mimeType: "image/png",
			wantErr:  false,
		},
		{
			name:     "valid GIF",
			mimeType: "image/gif",
			wantErr:  false,
		},
		{
			name:     "invalid MIME type",
			mimeType: "image/bmp",
			wantErr:  true,
		},
		{
			name:     "empty MIME type",
			mimeType: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateImageMimeType(tt.mimeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateImageMimeType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVideoMimeType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		wantErr  bool
	}{
		{
			name:     "valid MP4",
			mimeType: "video/mp4",
			wantErr:  false,
		},
		{
			name:     "valid MPEG",
			mimeType: "video/mpeg",
			wantErr:  false,
		},
		{
			name:     "invalid MIME type",
			mimeType: "video/avi",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVideoMimeType(tt.mimeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVideoMimeType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAudioMimeType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		wantErr  bool
	}{
		{
			name:     "valid MP3",
			mimeType: "audio/mp3",
			wantErr:  false,
		},
		{
			name:     "valid MPEG",
			mimeType: "audio/mpeg",
			wantErr:  false,
		},
		{
			name:     "invalid MIME type",
			mimeType: "audio/flac",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAudioMimeType(tt.mimeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAudioMimeType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid PDF",
			filename: "document.pdf",
			wantErr:  false,
		},
		{
			name:     "valid DOCX",
			filename: "document.docx",
			wantErr:  false,
		},
		{
			name:     "invalid extension",
			filename: "document.exe",
			wantErr:  true,
		},
		{
			name:     "no extension",
			filename: "document",
			wantErr:  true,
		},
		{
			name:     "empty filename",
			filename: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileExtension(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileExtension() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	tests := []struct {
		name     string
		size     int64
		fileType string
		wantErr  bool
	}{
		{
			name:     "valid image size",
			size:     5 * 1024 * 1024, // 5MB
			fileType: "image",
			wantErr:  false,
		},
		{
			name:     "image too large",
			size:     15 * 1024 * 1024, // 15MB
			fileType: "image",
			wantErr:  true,
		},
		{
			name:     "valid file size",
			size:     20 * 1024 * 1024, // 20MB
			fileType: "file",
			wantErr:  false,
		},
		{
			name:     "file too large",
			size:     30 * 1024 * 1024, // 30MB
			fileType: "file",
			wantErr:  true,
		},
		{
			name:     "zero size",
			size:     0,
			fileType: "image",
			wantErr:  true,
		},
		{
			name:     "negative size",
			size:     -1,
			fileType: "image",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateFileSize(tt.size, tt.fileType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateFileSize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid HTTP URL",
			url:     "http://example.com",
			wantErr: false,
		},
		{
			name:    "valid HTTPS URL",
			url:     "https://example.com/path",
			wantErr: false,
		},
		{
			name:    "invalid scheme",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "no scheme",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "invalid URL",
			url:     "ht!tp://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFormatMessage(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "trim whitespace",
			text: "  Hello World  ",
			want: "Hello World",
		},
		{
			name: "normalize line breaks",
			text: "Line1\r\nLine2\rLine3\nLine4",
			want: "Line1\nLine2\nLine3\nLine4",
		},
		{
			name: "remove excessive line breaks",
			text: "Line1\n\n\n\nLine2",
			want: "Line1\n\nLine2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMessage(tt.text)
			if got != tt.want {
				t.Errorf("FormatMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTruncateMessage(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		maxLength int
		wantLen   int
	}{
		{
			name:      "short message",
			text:      "Hello",
			maxLength: 100,
			wantLen:   5,
		},
		{
			name:      "truncate long message",
			text:      "This is a very long message that needs to be truncated",
			maxLength: 20,
			wantLen:   23, // 20 + "..."
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateMessage(tt.text, tt.maxLength)
			if len([]rune(got)) > tt.wantLen+5 { // Allow some variance for word boundaries
				t.Errorf("TruncateMessage() length = %d, want <= %d", len([]rune(got)), tt.wantLen+5)
			}
		})
	}
}

func TestSanitizeText(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "remove control characters",
			text: "Hello\x00World",
			want: "HelloWorld",
		},
		{
			name: "keep newlines and tabs",
			text: "Hello\nWorld\tTest",
			want: "Hello\nWorld\tTest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeText(tt.text)
			if got != tt.want {
				t.Errorf("SanitizeText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsEmptyOrWhitespace(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{
			name: "empty string",
			text: "",
			want: true,
		},
		{
			name: "whitespace only",
			text: "   ",
			want: true,
		},
		{
			name: "non-empty string",
			text: "Hello",
			want: false,
		},
		{
			name: "string with content and whitespace",
			text: "  Hello  ",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEmptyOrWhitespace(tt.text)
			if got != tt.want {
				t.Errorf("IsEmptyOrWhitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}
