package internal

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// requestKey is the slog attribute key used for the request id in logs.
const requestKey = "request_id"

// LoggerKey is an unexported type used as the context key for *slog.Logger.
type LoggerKey struct{}

// requestIDKey is an unexported type used as the context key for the request id value.
type requestIDKey struct{}

// LoggerFromContext returns the *slog.Logger stored in context, or slog.Default() if missing.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(LoggerKey{}).(*slog.Logger); ok && l != nil {
		return l
	}
	return slog.Default()
}

// RequestIDLogger is middleware that ensures each request has an id in the context
// and adds it as an attribute to the logger stored in the context.
func RequestIDLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("Request-Id")
		if id == "" {
			id = uuid.NewString()
		}
		logger := slog.Default().With(slog.String(requestKey, id))
		ctx := context.WithValue(r.Context(), LoggerKey{}, logger)
		ctx = context.WithValue(ctx, requestIDKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestIdFromContext retrieves the request id stored in the context.
// It returns an empty string if the value is missing or not a string.
func GetRequestIdFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey{}).(string); ok {
		return v
	}
	return ""
}

// GetRequestId retrieves the request id stored in the request's context.
// It returns an empty string if the value is missing or not a string.
func GetRequestId(r *http.Request) string {
	if r == nil || r.Context() == nil {
		return ""
	}
	if v, ok := r.Context().Value(requestIDKey{}).(string); ok {
		return v
	}
	return ""
}
