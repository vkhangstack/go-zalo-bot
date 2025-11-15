package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// LogLevelDebug is for detailed debugging information
	LogLevelDebug LogLevel = iota
	// LogLevelInfo is for general informational messages
	LogLevelInfo
	// LogLevelWarn is for warning messages
	LogLevelWarn
	// LogLevelError is for error messages
	LogLevelError
)

// String returns the string representation of LogLevel
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogFormat represents the log output format
type LogFormat int

const (
	// LogFormatText is plain text format
	LogFormatText LogFormat = iota
	// LogFormatJSON is JSON format
	LogFormatJSON
)

// Field represents a structured log field
type Field struct {
	Key   string
	Value interface{}
}

// Logger interface defines logging methods
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	SetLevel(level LogLevel)
	IsEnabled(level LogLevel) bool
}

// LogConfig represents logger configuration
type LogConfig struct {
	Level  LogLevel
	Output io.Writer
	Format LogFormat
}

// DefaultLogger is a simple logger implementation
type DefaultLogger struct {
	config LogConfig
	mu     sync.RWMutex
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config LogConfig) *DefaultLogger {
	if config.Output == nil {
		config.Output = os.Stdout
	}
	
	return &DefaultLogger{
		config: config,
	}
}

// NewDefaultLogger creates a logger with default settings (Info level, text format)
func NewDefaultLogger() *DefaultLogger {
	return NewLogger(LogConfig{
		Level:  LogLevelInfo,
		Output: os.Stdout,
		Format: LogFormatText,
	})
}

// SetLevel sets the logging level
func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.config.Level = level
}

// IsEnabled checks if a log level is enabled
func (l *DefaultLogger) IsEnabled(level LogLevel) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return level >= l.config.Level
}

// Debug logs a debug message
func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	l.log(LogLevelDebug, msg, fields...)
}

// Info logs an info message
func (l *DefaultLogger) Info(msg string, fields ...Field) {
	l.log(LogLevelInfo, msg, fields...)
}

// Warn logs a warning message
func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	l.log(LogLevelWarn, msg, fields...)
}

// Error logs an error message
func (l *DefaultLogger) Error(msg string, fields ...Field) {
	l.log(LogLevelError, msg, fields...)
}

// log is the internal logging method
func (l *DefaultLogger) log(level LogLevel, msg string, fields ...Field) {
	if !l.IsEnabled(level) {
		return
	}
	
	l.mu.RLock()
	format := l.config.Format
	output := l.config.Output
	l.mu.RUnlock()
	
	timestamp := time.Now()
	
	var logLine string
	if format == LogFormatJSON {
		logLine = l.formatJSON(timestamp, level, msg, fields)
	} else {
		logLine = l.formatText(timestamp, level, msg, fields)
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(output, logLine)
}

// formatText formats log entry as plain text
func (l *DefaultLogger) formatText(timestamp time.Time, level LogLevel, msg string, fields []Field) string {
	// Format: 2006-01-02 15:04:05 [LEVEL] message key=value key=value
	line := fmt.Sprintf("%s [%s] %s", 
		timestamp.Format("2006-01-02 15:04:05"), 
		level.String(), 
		msg)
	
	if len(fields) > 0 {
		for _, field := range fields {
			line += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}
	
	return line
}

// formatJSON formats log entry as JSON
func (l *DefaultLogger) formatJSON(timestamp time.Time, level LogLevel, msg string, fields []Field) string {
	entry := map[string]interface{}{
		"timestamp": timestamp.Format(time.RFC3339),
		"level":     level.String(),
		"message":   msg,
	}
	
	if len(fields) > 0 {
		fieldsMap := make(map[string]interface{})
		for _, field := range fields {
			fieldsMap[field.Key] = field.Value
		}
		entry["fields"] = fieldsMap
	}
	
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple format if JSON marshaling fails
		return fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s","error":"failed to marshal fields"}`,
			timestamp.Format(time.RFC3339), level.String(), msg)
	}
	
	return string(jsonBytes)
}

// NoOpLogger is a logger that does nothing (for disabling logging)
type NoOpLogger struct{}

// NewNoOpLogger creates a no-op logger
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

// Debug does nothing
func (l *NoOpLogger) Debug(msg string, fields ...Field) {}

// Info does nothing
func (l *NoOpLogger) Info(msg string, fields ...Field) {}

// Warn does nothing
func (l *NoOpLogger) Warn(msg string, fields ...Field) {}

// Error does nothing
func (l *NoOpLogger) Error(msg string, fields ...Field) {}

// SetLevel does nothing
func (l *NoOpLogger) SetLevel(level LogLevel) {}

// IsEnabled always returns false
func (l *NoOpLogger) IsEnabled(level LogLevel) bool {
	return false
}
