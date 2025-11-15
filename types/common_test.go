package types

import (
	"encoding/json"
	"testing"
)

func TestLogLevel_IsValid(t *testing.T) {
	tests := []struct {
		name string
		ll   LogLevel
		want bool
	}{
		{"valid debug", LogLevelDebug, true},
		{"valid info", LogLevelInfo, true},
		{"valid warn", LogLevelWarn, true},
		{"valid error", LogLevelError, true},
		{"invalid level", LogLevel("invalid"), false},
		{"empty level", LogLevel(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ll.IsValid(); got != tt.want {
				t.Errorf("LogLevel.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		name string
		ll   LogLevel
		want string
	}{
		{"debug level", LogLevelDebug, "debug"},
		{"info level", LogLevelInfo, "info"},
		{"warn level", LogLevelWarn, "warn"},
		{"error level", LogLevelError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ll.String(); got != tt.want {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIResponse_JSON(t *testing.T) {
	tests := []struct {
		name     string
		response APIResponse
	}{
		{
			name: "successful response",
			response: APIResponse{
				OK:     true,
				Result: json.RawMessage(`{"message_id":"msg123"}`),
			},
		},
		{
			name: "error response",
			response: APIResponse{
				OK:          false,
				ErrorCode:   400,
				Description: "Bad Request",
			},
		},
		{
			name: "response with all fields",
			response: APIResponse{
				OK:          true,
				Result:      json.RawMessage(`{"data":"test"}`),
				ErrorCode:   0,
				Description: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshal
			data, err := json.Marshal(tt.response)
			if err != nil {
				t.Errorf("APIResponse.MarshalJSON() error = %v", err)
				return
			}

			// Test unmarshal
			var unmarshaled APIResponse
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Errorf("APIResponse.UnmarshalJSON() error = %v", err)
				return
			}

			if unmarshaled.OK != tt.response.OK {
				t.Errorf("APIResponse.OK = %v, want %v", unmarshaled.OK, tt.response.OK)
			}
			if unmarshaled.ErrorCode != tt.response.ErrorCode {
				t.Errorf("APIResponse.ErrorCode = %v, want %v", unmarshaled.ErrorCode, tt.response.ErrorCode)
			}
			if unmarshaled.Description != tt.response.Description {
				t.Errorf("APIResponse.Description = %v, want %v", unmarshaled.Description, tt.response.Description)
			}
		})
	}
}

func TestPaginatedResponse_JSON(t *testing.T) {
	paginatedResp := PaginatedResponse{
		Data:       json.RawMessage(`[{"id":"1"},{"id":"2"}]`),
		TotalCount: 100,
		HasMore:    true,
	}

	// Test marshal
	data, err := json.Marshal(paginatedResp)
	if err != nil {
		t.Errorf("PaginatedResponse.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled PaginatedResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("PaginatedResponse.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.TotalCount != paginatedResp.TotalCount {
		t.Errorf("PaginatedResponse.TotalCount = %v, want %v", unmarshaled.TotalCount, paginatedResp.TotalCount)
	}
	if unmarshaled.HasMore != paginatedResp.HasMore {
		t.Errorf("PaginatedResponse.HasMore = %v, want %v", unmarshaled.HasMore, paginatedResp.HasMore)
	}
	if string(unmarshaled.Data) != string(paginatedResp.Data) {
		t.Errorf("PaginatedResponse.Data = %v, want %v", string(unmarshaled.Data), string(paginatedResp.Data))
	}
}

func TestAPIResponse_ParseResult(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantOK  bool
		wantErr bool
	}{
		{
			name:    "valid success response",
			json:    `{"ok":true,"result":{"message_id":"msg123"}}`,
			wantOK:  true,
			wantErr: false,
		},
		{
			name:    "valid error response",
			json:    `{"ok":false,"error_code":400,"description":"Bad Request"}`,
			wantOK:  false,
			wantErr: false,
		},
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantOK:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response APIResponse
			err := json.Unmarshal([]byte(tt.json), &response)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr && response.OK != tt.wantOK {
				t.Errorf("APIResponse.OK = %v, want %v", response.OK, tt.wantOK)
			}
		})
	}
}

func TestPaginatedResponse_EmptyData(t *testing.T) {
	paginatedResp := PaginatedResponse{
		Data:       json.RawMessage(`[]`),
		TotalCount: 0,
		HasMore:    false,
	}

	// Test marshal
	data, err := json.Marshal(paginatedResp)
	if err != nil {
		t.Errorf("PaginatedResponse.MarshalJSON() error = %v", err)
		return
	}

	// Test unmarshal
	var unmarshaled PaginatedResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("PaginatedResponse.UnmarshalJSON() error = %v", err)
		return
	}

	if unmarshaled.TotalCount != 0 {
		t.Errorf("PaginatedResponse.TotalCount = %v, want 0", unmarshaled.TotalCount)
	}
	if unmarshaled.HasMore != false {
		t.Errorf("PaginatedResponse.HasMore = %v, want false", unmarshaled.HasMore)
	}
}
