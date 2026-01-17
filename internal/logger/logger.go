package logger

import (
	"encoding/json" // Added for JSON marshaling
	"fmt"
	"os"
	"time" // Added for timestamps
)

// Logger defines the interface for logging.
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
	Fatal(code int, format string, args ...interface{})
	Debug(format string, args ...interface{})
}

// HumanReadableLogger implements Logger for human-readable output.
type HumanReadableLogger struct{}

// NewHumanReadableLogger creates a new HumanReadableLogger.
func NewHumanReadableLogger() *HumanReadableLogger {
	return &HumanReadableLogger{}
}

// Info prints informational messages to stdout.
func (l *HumanReadableLogger) Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Error prints error messages to stderr.
func (l *HumanReadableLogger) Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

// Fatal prints a fatal error message to stderr and exits with the given code.
func (l *HumanReadableLogger) Fatal(code int, format string, args ...interface{}) {
	if format != "" {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
	os.Exit(code)
}

// Debug prints debug messages (currently to stdout, can be conditional).
func (l *HumanReadableLogger) Debug(format string, args ...interface{}) {
	// For now, debug messages are just info messages.
	// This can be made conditional based on a verbose flag later.
	fmt.Printf("DEBUG: "+format+"\n", args...)
}

// LogEntry defines the schema for JSON log output.
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Code      int    `json:"code,omitempty"` // For fatal errors
}

// JSONLogger implements Logger for JSON output.
type JSONLogger struct{}

// NewJSONLogger creates a new JSONLogger.
func NewJSONLogger() *JSONLogger {
	return &JSONLogger{}
}

func (l *JSONLogger) logJSON(level string, format string, args ...interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   fmt.Sprintf(format, args...),
	}
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		// Fallback to human-readable error if JSON marshaling fails
		fmt.Fprintf(os.Stderr, "ERROR: Failed to marshal log entry to JSON: %v - %s\n", err, fmt.Sprintf(format, args...))
		return
	}
	fmt.Println(string(jsonBytes))
}

// Info prints informational messages to stdout in JSON format.
func (l *JSONLogger) Info(format string, args ...interface{}) {
	l.logJSON("info", format, args...)
}

// Error prints error messages to stderr in JSON format.
func (l *JSONLogger) Error(format string, args ...interface{}) {
	l.logJSON("error", format, args...)
}

// Fatal prints a fatal error message to stderr in JSON format and exits with the given code.
func (l *JSONLogger) Fatal(code int, format string, args ...interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "fatal",
		Message:   fmt.Sprintf(format, args...),
		Code:      code,
	}
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to marshal fatal log entry to JSON: %v - %s\n", err, fmt.Sprintf(format, args...))
	} else {
		fmt.Println(string(jsonBytes))
	}
	os.Exit(code)
}

// Debug prints debug messages to stdout in JSON format.
func (l *JSONLogger) Debug(format string, args ...interface{}) {
	l.logJSON("debug", format, args...)
}

