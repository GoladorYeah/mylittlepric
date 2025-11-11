package utils

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// ContextKey represents a key for storing values in context
type ContextKey string

const (
	// RequestIDKey is the context key for request ID
	RequestIDKey ContextKey = "request_id"
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// SessionIDKey is the context key for session ID
	SessionIDKey ContextKey = "session_id"
)

var (
	logger     *slog.Logger
	lokiWriter *LokiWriter
)

// InitLogger initializes the global logger with the specified level and format
func InitLogger(level string, format string, lokiEnabled bool, lokiURL string, serviceName string) {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	// Create output writer
	var output io.Writer = os.Stdout

	// If Loki is enabled, create a MultiWriter
	if lokiEnabled && lokiURL != "" {
		labels := map[string]string{
			"service": serviceName,
			"job":     serviceName,
			"level":   level,
		}
		lokiWriter = NewLokiWriter(lokiURL, labels)
		output = io.MultiWriter(os.Stdout, lokiWriter)
	}

	if format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
}

// CloseLoki closes the Loki writer and flushes remaining logs
func CloseLoki() error {
	if lokiWriter != nil {
		return lokiWriter.Close()
	}
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *slog.Logger {
	if logger == nil {
		// Initialize with defaults if not initialized
		InitLogger("info", "json", false, "", "")
	}
	return logger
}

// LogDebug logs a debug message with context
func LogDebug(ctx context.Context, msg string, args ...any) {
	logger.DebugContext(ctx, msg, extractContextAttrs(ctx, args)...)
}

// LogInfo logs an info message with context
func LogInfo(ctx context.Context, msg string, args ...any) {
	logger.InfoContext(ctx, msg, extractContextAttrs(ctx, args)...)
}

// LogWarn logs a warning message with context
func LogWarn(ctx context.Context, msg string, args ...any) {
	logger.WarnContext(ctx, msg, extractContextAttrs(ctx, args)...)
}

// LogError logs an error message with context
func LogError(ctx context.Context, msg string, err error, args ...any) {
	attrs := extractContextAttrs(ctx, args)
	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
	}
	logger.ErrorContext(ctx, msg, attrs...)
}

// extractContextAttrs extracts relevant attributes from context
func extractContextAttrs(ctx context.Context, args []any) []any {
	if ctx == nil {
		return args
	}

	attrs := make([]any, 0, len(args)+3)

	// Extract request_id from context
	if reqID := ctx.Value(RequestIDKey); reqID != nil {
		if reqIDStr, ok := reqID.(string); ok {
			attrs = append(attrs, slog.String("request_id", reqIDStr))
		}
	}

	// Extract user_id from context
	if userID := ctx.Value(UserIDKey); userID != nil {
		if userIDStr, ok := userID.(string); ok {
			attrs = append(attrs, slog.String("user_id", userIDStr))
		}
	}

	// Extract session_id from context
	if sessionID := ctx.Value(SessionIDKey); sessionID != nil {
		if sessionIDStr, ok := sessionID.(string); ok {
			attrs = append(attrs, slog.String("session_id", sessionIDStr))
		}
	}

	return append(attrs, args...)
}

// WithRequestID adds request ID to context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithSessionID adds session ID to context
func WithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, SessionIDKey, sessionID)
}
