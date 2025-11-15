package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	// ErrInvalidSignature is returned when webhook signature validation fails
	ErrInvalidSignature = errors.New("invalid webhook signature")
	// ErrEmptySignature is returned when signature is empty
	ErrEmptySignature = errors.New("empty webhook signature")
	// ErrEmptySecret is returned when secret token is empty
	ErrEmptySecret = errors.New("empty secret token")
	// ErrEmptyPayload is returned when payload is empty
	ErrEmptyPayload = errors.New("empty webhook payload")
	// ErrInvalidUserID is returned when user ID is invalid
	ErrInvalidUserID = errors.New("invalid user ID")
	// ErrInvalidRecipientID is returned when recipient ID is invalid
	ErrInvalidRecipientID = errors.New("invalid recipient ID")
	// ErrInvalidMessageContent is returned when message content is invalid
	ErrInvalidMessageContent = errors.New("invalid message content")
	// ErrMessageTooLong is returned when message exceeds maximum length
	ErrMessageTooLong = errors.New("message content exceeds maximum length")
	// ErrInvalidFileFormat is returned when file format is not supported
	ErrInvalidFileFormat = errors.New("invalid file format")
	// ErrFileTooLarge is returned when file size exceeds limit
	ErrFileTooLarge = errors.New("file size exceeds limit")
	// ErrInvalidURL is returned when URL is invalid
	ErrInvalidURL = errors.New("invalid URL")
	// ErrInvalidMimeType is returned when MIME type is invalid
	ErrInvalidMimeType = errors.New("invalid MIME type")
)

const (
	// MaxMessageLength is the maximum length for message content
	MaxMessageLength = 5000
	// MaxImageSize is the maximum size for image attachments (10MB)
	MaxImageSize = 10 * 1024 * 1024
	// MaxFileSize is the maximum size for file attachments (25MB)
	MaxFileSize = 25 * 1024 * 1024
	// MaxVideoSize is the maximum size for video attachments (50MB)
	MaxVideoSize = 50 * 1024 * 1024
	// MaxAudioSize is the maximum size for audio attachments (10MB)
	MaxAudioSize = 10 * 1024 * 1024
)

var (
	// userIDPattern matches valid user IDs (alphanumeric and underscores, 1-64 chars)
	userIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_]{1,64}$`)
	
	// Supported image MIME types
	supportedImageMimeTypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	
	// Supported video MIME types
	supportedVideoMimeTypes = map[string]bool{
		"video/mp4":  true,
		"video/mpeg": true,
		"video/webm": true,
	}
	
	// Supported audio MIME types
	supportedAudioMimeTypes = map[string]bool{
		"audio/mpeg": true,
		"audio/mp3":  true,
		"audio/wav":  true,
		"audio/ogg":  true,
	}
	
	// Supported file extensions
	supportedFileExtensions = map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".ppt":  true,
		".pptx": true,
		".txt":  true,
		".zip":  true,
		".rar":  true,
	}
)

// VerifyWebhookSignature verifies webhook signature using HMAC-SHA256
// This function provides secure validation using secret tokens
// Returns true if signature is valid, false otherwise
func VerifyWebhookSignature(payload []byte, signature, secret string) bool {
	// Validate inputs
	if len(payload) == 0 || strings.TrimSpace(signature) == "" || strings.TrimSpace(secret) == "" {
		return false
	}

	// Compute HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))

	// Use constant-time comparison to prevent timing attacks
	return hmac.Equal([]byte(signature), []byte(expected))
}

// ValidateWebhookSignature validates webhook signature and returns detailed error
// This function provides enhanced validation with specific error messages
func ValidateWebhookSignature(payload []byte, signature, secret string) error {
	// Check for empty payload
	if len(payload) == 0 {
		return ErrEmptyPayload
	}

	// Check for empty signature
	if strings.TrimSpace(signature) == "" {
		return ErrEmptySignature
	}

	// Check for empty secret
	if strings.TrimSpace(secret) == "" {
		return ErrEmptySecret
	}

	// Verify signature
	if !VerifyWebhookSignature(payload, signature, secret) {
		return ErrInvalidSignature
	}

	return nil
}

// RejectInvalidWebhookRequest creates an error for rejecting invalid webhook requests
// This provides a standardized way to reject requests with invalid signatures
func RejectInvalidWebhookRequest(reason string) error {
	if reason == "" {
		reason = "invalid webhook request"
	}
	return errors.New("webhook request rejected: " + reason)
}

// ToJSON converts struct to JSON
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// FromJSON parses JSON to struct
func FromJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// StringPtr returns pointer to string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns pointer to int
func IntPtr(i int) *int {
	return &i
}

// BoolPtr returns pointer to bool
func BoolPtr(b bool) *bool {
	return &b
}

// ValidateUserID validates a user ID format
// User IDs should be alphanumeric with underscores, 1-64 characters
func ValidateUserID(userID string) error {
	if userID == "" {
		return ErrInvalidUserID
	}
	
	if !userIDPattern.MatchString(userID) {
		return fmt.Errorf("%w: must be alphanumeric with underscores, 1-64 characters", ErrInvalidUserID)
	}
	
	return nil
}

// ValidateRecipientID validates a recipient ID format
// Recipient IDs follow the same rules as user IDs
func ValidateRecipientID(recipientID string) error {
	if recipientID == "" {
		return ErrInvalidRecipientID
	}
	
	if !userIDPattern.MatchString(recipientID) {
		return fmt.Errorf("%w: must be alphanumeric with underscores, 1-64 characters", ErrInvalidRecipientID)
	}
	
	return nil
}

// ValidateMessageContent validates message content
// Checks for empty content, length limits, and Unicode support
func ValidateMessageContent(content string) error {
	if content == "" {
		return ErrInvalidMessageContent
	}
	
	// Check UTF-8 validity
	if !utf8.ValidString(content) {
		return fmt.Errorf("%w: content must be valid UTF-8", ErrInvalidMessageContent)
	}
	
	// Check length (count runes, not bytes, to properly handle Unicode)
	runeCount := utf8.RuneCountInString(content)
	if runeCount > MaxMessageLength {
		return fmt.Errorf("%w: %d characters (max %d)", ErrMessageTooLong, runeCount, MaxMessageLength)
	}
	
	return nil
}

// ContainsVietnamese checks if the text contains Vietnamese characters
// This is useful for ensuring Vietnamese character support is working
func ContainsVietnamese(text string) bool {
	// Vietnamese Unicode ranges
	// Latin Extended Additional (U+1E00-U+1EFF) contains Vietnamese diacritics
	for _, r := range text {
		if (r >= 0x1E00 && r <= 0x1EFF) || // Vietnamese diacritics
			r == 0x0110 || r == 0x0111 || // Đ, đ
			r == 0x00C0 || r == 0x00C1 || r == 0x00C2 || r == 0x00C3 || // À, Á, Â, Ã
			r == 0x00C8 || r == 0x00C9 || r == 0x00CA || // È, É, Ê
			r == 0x00CC || r == 0x00CD || // Ì, Í
			r == 0x00D2 || r == 0x00D3 || r == 0x00D4 || r == 0x00D5 || // Ò, Ó, Ô, Õ
			r == 0x00D9 || r == 0x00DA || // Ù, Ú
			r == 0x00DD || // Ý
			r == 0x00E0 || r == 0x00E1 || r == 0x00E2 || r == 0x00E3 || // à, á, â, ã
			r == 0x00E8 || r == 0x00E9 || r == 0x00EA || // è, é, ê
			r == 0x00EC || r == 0x00ED || // ì, í
			r == 0x00F2 || r == 0x00F3 || r == 0x00F4 || r == 0x00F5 || // ò, ó, ô, õ
			r == 0x00F9 || r == 0x00FA || // ù, ú
			r == 0x00FD { // ý
			return true
		}
	}
	return false
}

// ValidateUnicodeSupport validates that the text is properly encoded and supports Unicode
func ValidateUnicodeSupport(text string) error {
	if !utf8.ValidString(text) {
		return fmt.Errorf("text contains invalid UTF-8 sequences")
	}
	
	// Check for replacement characters which indicate encoding issues
	if strings.Contains(text, "\uFFFD") {
		return fmt.Errorf("text contains replacement characters, indicating encoding issues")
	}
	
	return nil
}

// ValidateImageMimeType validates image MIME type
func ValidateImageMimeType(mimeType string) error {
	if mimeType == "" {
		return ErrInvalidMimeType
	}
	
	if !supportedImageMimeTypes[strings.ToLower(mimeType)] {
		return fmt.Errorf("%w: %s is not a supported image format", ErrInvalidMimeType, mimeType)
	}
	
	return nil
}

// ValidateVideoMimeType validates video MIME type
func ValidateVideoMimeType(mimeType string) error {
	if mimeType == "" {
		return ErrInvalidMimeType
	}
	
	if !supportedVideoMimeTypes[strings.ToLower(mimeType)] {
		return fmt.Errorf("%w: %s is not a supported video format", ErrInvalidMimeType, mimeType)
	}
	
	return nil
}

// ValidateAudioMimeType validates audio MIME type
func ValidateAudioMimeType(mimeType string) error {
	if mimeType == "" {
		return ErrInvalidMimeType
	}
	
	if !supportedAudioMimeTypes[strings.ToLower(mimeType)] {
		return fmt.Errorf("%w: %s is not a supported audio format", ErrInvalidMimeType, mimeType)
	}
	
	return nil
}

// ValidateFileExtension validates file extension
func ValidateFileExtension(filename string) error {
	if filename == "" {
		return ErrInvalidFileFormat
	}
	
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		return fmt.Errorf("%w: file has no extension", ErrInvalidFileFormat)
	}
	
	if !supportedFileExtensions[ext] {
		return fmt.Errorf("%w: %s is not a supported file extension", ErrInvalidFileFormat, ext)
	}
	
	return nil
}

// ValidateFileSize validates file size against limits based on file type
func ValidateFileSize(size int64, fileType string) error {
	if size <= 0 {
		return fmt.Errorf("%w: file size must be positive", ErrFileTooLarge)
	}
	
	var maxSize int64
	switch strings.ToLower(fileType) {
	case "image":
		maxSize = MaxImageSize
	case "video":
		maxSize = MaxVideoSize
	case "audio":
		maxSize = MaxAudioSize
	case "file":
		maxSize = MaxFileSize
	default:
		maxSize = MaxFileSize
	}
	
	if size > maxSize {
		return fmt.Errorf("%w: %d bytes exceeds maximum of %d bytes for %s", 
			ErrFileTooLarge, size, maxSize, fileType)
	}
	
	return nil
}

// ValidateURL validates a URL format
func ValidateURL(urlStr string) error {
	if urlStr == "" {
		return ErrInvalidURL
	}
	
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidURL, err.Error())
	}
	
	// Check for valid scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("%w: scheme must be http or https", ErrInvalidURL)
	}
	
	// Check for valid host
	if parsedURL.Host == "" {
		return fmt.Errorf("%w: missing host", ErrInvalidURL)
	}
	
	return nil
}

// FormatMessage formats a message with proper line breaks and trimming
func FormatMessage(text string) string {
	// Trim leading and trailing whitespace
	text = strings.TrimSpace(text)
	
	// Normalize line breaks
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	
	// Remove excessive consecutive line breaks (more than 2)
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}
	
	return text
}

// TruncateMessage truncates a message to the maximum length, preserving whole words
func TruncateMessage(text string, maxLength int) string {
	if maxLength <= 0 {
		maxLength = MaxMessageLength
	}
	
	runes := []rune(text)
	if len(runes) <= maxLength {
		return text
	}
	
	// Truncate to maxLength
	truncated := runes[:maxLength]
	
	// Try to find the last space to avoid cutting words
	lastSpace := -1
	for i := len(truncated) - 1; i >= 0; i-- {
		if unicode.IsSpace(truncated[i]) {
			lastSpace = i
			break
		}
	}
	
	// If we found a space in the last 20% of the text, cut there
	if lastSpace > int(float64(maxLength)*0.8) {
		truncated = truncated[:lastSpace]
	}
	
	return string(truncated) + "..."
}

// SanitizeText removes control characters and normalizes whitespace
func SanitizeText(text string) string {
	// Remove control characters except newlines and tabs
	var builder strings.Builder
	for _, r := range text {
		if r == '\n' || r == '\t' || !unicode.IsControl(r) {
			builder.WriteRune(r)
		}
	}
	
	return builder.String()
}

// IsEmptyOrWhitespace checks if a string is empty or contains only whitespace
func IsEmptyOrWhitespace(text string) bool {
	return strings.TrimSpace(text) == ""
}
