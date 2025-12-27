// Package logger provides structured logging functionality with multiple levels and output formats.
package logger

import "context"

// Logger defines the interface for structured logging.
// Implementations should be safe for concurrent use.
type Logger interface {
	// Debug logs a debug message with optional key-value pairs
	Debug(msg string, keysAndValues ...any)

	// Info logs an informational message with optional key-value pairs
	Info(msg string, keysAndValues ...any)

	// Warn logs a warning message with optional key-value pairs
	Warn(msg string, keysAndValues ...any)

	// Error logs an error message with optional key-value pairs
	Error(msg string, keysAndValues ...any)

	// With returns a new Logger with additional context fields
	// This allows creating child loggers with pre-set fields
	With(keysAndValues ...any) Logger

	// WithContext returns a new Logger with context
	WithContext(ctx context.Context) Logger
}

// Level represents the logging level.
type Level int

const (
	// LevelDebug is the lowest logging level.
	LevelDebug Level = iota
	// LevelInfo is the default logging level.
	LevelInfo
	// LevelWarn is for warning messages.
	LevelWarn
	// LevelError is for error messages.
	LevelError
)

// Format represents the output format.
type Format string

const (
	// FormatJSON represents JSON output format.
	FormatJSON Format = "json"
	// FormatText represents plain text output format.
	FormatText Format = "text"
)

// Config holds logger configuration.
type Config struct {
	Level  Level  // Minimum level to log
	Format Format // Output format (json or text)
}
