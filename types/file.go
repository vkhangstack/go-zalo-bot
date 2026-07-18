package types

type Photo struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int64  `json:"file_size,omitempty"`
	// URL is the image URL as sent in an incoming webhook message.
	URL string `json:"url,omitempty"`
}

type Sticker struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	IsAnimated   bool   `json:"is_animated"`
	FileSize     int64  `json:"file_size,omitempty"`
	// URL is the sticker image URL as sent in an incoming webhook message.
	URL string `json:"url,omitempty"`
}
