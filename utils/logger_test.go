package utils

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		name  string
		level LogLevel
		want  string
	}{
		{
			name:  "debug level",
			level: LogLevelDebug,
			want:  "DEBUG",
		},
		{
			name:  "info level",
			level: LogLevelInfo,
			want:  "INFO",
		},
		{
			name:  "warn level",
			level: LogLevelWarn,
			want:  "WARN",
		},
		{
			name:  "error level",
			level: LogLevelError,
			want:  "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.want {
				t.Errorf("LogLevel.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	config := LogConfig{
		Level:  LogLevelInfo,
		Output: buf,
		Format: LogFormatText,
	}

	logger := NewLogger(config)
	if logger == nil {
		t.Error("NewLogger() returned nil")
	}

	if !logger.IsEnabled(LogLevelInfo) {
		t.Error("Logger should be enabled for Info level")
	}
}

func TestNewDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()
	if logger == nil {
		t.Error("NewDefaultLogger() returned nil")
	}

	if !logger.IsEnabled(LogLevelInfo) {
		t.Error("Default logger should be enabled for Info level")
	}

	if logger.IsEnabled(LogLevelDebug) {
		t.Error("Default logger should not be enabled for Debug level")
	}
}

func TestDefaultLogger_SetLevel(t *testing.T) {
	logger := NewDefaultLogger()

	// Initially at Info level
	if logger.IsEnabled(LogLevelDebug) {
		t.Error("Logger should not be enabled for Debug initially")
	}

	// Set to Debug level
	logger.SetLevel(LogLevelDebug)
	if !logger.IsEnabled(LogLevelDebug) {
		t.Error("Logger should be enabled for Debug after SetLevel")
	}

	// Set to Error level
	logger.SetLevel(LogLevelError)
	if logger.IsEnabled(LogLevelInfo) {
		t.Error("Logger should not be enabled for Info after setting to Error")
	}
	if !logger.IsEnabled(LogLevelError) {
		t.Error("Logger should be enabled for Error")
	}
}

func TestDefaultLogger_IsEnabled(t *testing.T) {
	tests := []struct {
		name        string
		configLevel LogLevel
		checkLevel  LogLevel
		want        bool
	}{
		{
			name:        "debug enabled at debug level",
			configLevel: LogLevelDebug,
			checkLevel:  LogLevelDebug,
			want:        true,
		},
		{
			name:        "info enabled at debug level",
			configLevel: LogLevelDebug,
			checkLevel:  LogLevelInfo,
			want:        true,
		},
		{
			name:        "debug not enabled at info level",
			configLevel: LogLevelInfo,
			checkLevel:  LogLevelDebug,
			want:        false,
		},
		{
			name:        "error enabled at info level",
			configLevel: LogLevelInfo,
			checkLevel:  LogLevelError,
			want:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(LogConfig{
				Level:  tt.configLevel,
				Output: &bytes.Buffer{},
				Format: LogFormatText,
			})

			got := logger.IsEnabled(tt.checkLevel)
			if got != tt.want {
				t.Errorf("IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultLogger_TextFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelDebug,
		Output: buf,
		Format: LogFormatText,
	})

	logger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "[INFO]") {
		t.Error("Output should contain [INFO] level")
	}
	if !strings.Contains(output, "test message") {
		t.Error("Output should contain the message")
	}
}

func TestDefaultLogger_TextFormatWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelDebug,
		Output: buf,
		Format: LogFormatText,
	})

	logger.Info("test message", Field{Key: "user_id", Value: "123"}, Field{Key: "action", Value: "login"})
	output := buf.String()

	if !strings.Contains(output, "user_id=123") {
		t.Error("Output should contain user_id field")
	}
	if !strings.Contains(output, "action=login") {
		t.Error("Output should contain action field")
	}
}

func TestDefaultLogger_JSONFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelDebug,
		Output: buf,
		Format: LogFormatJSON,
	})

	logger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, `"level":"INFO"`) {
		t.Error("JSON output should contain level field")
	}
	if !strings.Contains(output, `"message":"test message"`) {
		t.Error("JSON output should contain message field")
	}
	if !strings.Contains(output, `"timestamp"`) {
		t.Error("JSON output should contain timestamp field")
	}
}

func TestDefaultLogger_JSONFormatWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelDebug,
		Output: buf,
		Format: LogFormatJSON,
	})

	logger.Info("test message", Field{Key: "user_id", Value: "123"})
	output := buf.String()

	if !strings.Contains(output, `"fields"`) {
		t.Error("JSON output should contain fields object")
	}
	if !strings.Contains(output, `"user_id":"123"`) {
		t.Error("JSON output should contain user_id field")
	}
}

func TestDefaultLogger_Debug(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelDebug,
		Output: buf,
		Format: LogFormatText,
	})

	logger.Debug("debug message")
	output := buf.String()

	if !strings.Contains(output, "[DEBUG]") {
		t.Error("Output should contain [DEBUG] level")
	}
	if !strings.Contains(output, "debug message") {
		t.Error("Output should contain the debug message")
	}
}

func TestDefaultLogger_Warn(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelDebug,
		Output: buf,
		Format: LogFormatText,
	})

	logger.Warn("warning message")
	output := buf.String()

	if !strings.Contains(output, "[WARN]") {
		t.Error("Output should contain [WARN] level")
	}
	if !strings.Contains(output, "warning message") {
		t.Error("Output should contain the warning message")
	}
}

func TestDefaultLogger_Error(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelDebug,
		Output: buf,
		Format: LogFormatText,
	})

	logger.Error("error message")
	output := buf.String()

	if !strings.Contains(output, "[ERROR]") {
		t.Error("Output should contain [ERROR] level")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Output should contain the error message")
	}
}

func TestDefaultLogger_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelWarn,
		Output: buf,
		Format: LogFormatText,
	})

	// These should not be logged
	logger.Debug("debug message")
	logger.Info("info message")

	// These should be logged
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	if strings.Contains(output, "debug message") {
		t.Error("Debug message should not be logged at Warn level")
	}
	if strings.Contains(output, "info message") {
		t.Error("Info message should not be logged at Warn level")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should be logged")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Error message should be logged")
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := NewNoOpLogger()

	// Should not panic
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	logger.SetLevel(LogLevelDebug)

	// IsEnabled should always return false
	if logger.IsEnabled(LogLevelDebug) {
		t.Error("NoOpLogger.IsEnabled() should always return false")
	}
	if logger.IsEnabled(LogLevelError) {
		t.Error("NoOpLogger.IsEnabled() should always return false")
	}
}

func TestLogger_ConcurrentAccess(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(LogConfig{
		Level:  LogLevelDebug,
		Output: buf,
		Format: LogFormatText,
	})

	// Test concurrent logging
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			logger.Info("concurrent message", Field{Key: "id", Value: id})
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic and should have logged messages
	if buf.Len() == 0 {
		t.Error("Logger should have logged messages")
	}
}

func TestLogger_ConcurrentSetLevel(t *testing.T) {
	logger := NewDefaultLogger()

	// Test concurrent SetLevel and IsEnabled calls
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			if id%2 == 0 {
				logger.SetLevel(LogLevelDebug)
			} else {
				logger.IsEnabled(LogLevelInfo)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should not panic
}
