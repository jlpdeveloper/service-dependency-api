package internal

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// captureState holds shared state for the capturingHandler so children produced by
// WithAttrs write into the same storage.
type captureState struct {
	lastAttrs []slog.Attr
}

type capturingHandler struct {
	state *captureState
	attrs []slog.Attr
}

func newCapturingHandler() *capturingHandler { return &capturingHandler{state: &captureState{}} }

func (h *capturingHandler) Enabled(ctx context.Context, level slog.Level) bool { return true }

func (h *capturingHandler) Handle(ctx context.Context, r slog.Record) error {
	merged := make([]slog.Attr, 0, len(h.attrs)+4)
	merged = append(merged, h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		merged = append(merged, a)
		return true
	})
	h.state.lastAttrs = merged
	return nil
}

func (h *capturingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a child that shares the same state but has accumulated attrs.
	acc := make([]slog.Attr, 0, len(h.attrs)+len(attrs))
	acc = append(acc, h.attrs...)
	acc = append(acc, attrs...)
	return &capturingHandler{state: h.state, attrs: acc}
}

func (h *capturingHandler) WithGroup(name string) slog.Handler { return h }

func (h *capturingHandler) lastAllAttrs() []slog.Attr { return append([]slog.Attr{}, h.state.lastAttrs...) }

func TestLoggerFromContext_Default(t *testing.T) {
	ctx := context.Background()
	l := LoggerFromContext(ctx)
	if l != slog.Default() {
		t.Fatalf("expected default logger when none set in context")
	}
}

func TestRequestIDLogger_UsesHeaderAndSetsLoggerAndContext(t *testing.T) {
	cap := newCapturingHandler()
	orig := slog.Default()
	slog.SetDefault(slog.New(cap))
	t.Cleanup(func() { slog.SetDefault(orig) })

	expectedID := "req-12345"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request id should be available via helper
		if got := GetRequestId(r); got != expectedID {
			t.Fatalf("expected GetRequestId to return %q, got %q", expectedID, got)
		}
		// Logger with request_id attribute should be in context
		LoggerFromContext(r.Context()).Info("test log")
		w.WriteHeader(http.StatusOK)
	})

	h := RequestIDLogger(next)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Request-Id", expectedID)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, r)

	if status := rec.Result().StatusCode; status != http.StatusOK {
		t.Fatalf("unexpected status: %d", status)
	}

	attrs := cap.lastAllAttrs()
	found := false
	for _, a := range attrs {
		if a.Key == requestKey && a.Value.String() == expectedID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("logger did not include attr %q=%q; got attrs: %#v", requestKey, expectedID, attrs)
	}
}

func TestRequestIDLogger_GeneratesUUIDWhenMissing(t *testing.T) {
	cap := newCapturingHandler()
	orig := slog.Default()
	slog.SetDefault(slog.New(cap))
	t.Cleanup(func() { slog.SetDefault(orig) })

	var capturedID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = GetRequestId(r)
		LoggerFromContext(r.Context()).Info("test log")
		w.WriteHeader(http.StatusOK)
	})

	h := RequestIDLogger(next)
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, r)

	if _, err := uuid.Parse(capturedID); err != nil {
		t.Fatalf("expected generated request id to be a UUID, got %q: %v", capturedID, err)
	}

	attrs := cap.lastAllAttrs()
	found := false
	for _, a := range attrs {
		if a.Key == requestKey && a.Value.String() == capturedID {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("logger did not include generated request id attr %q=%q; got attrs: %#v", requestKey, capturedID, attrs)
	}
}

func TestGetRequestId_EmptyWhenMissing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	if got := GetRequestId(r); got != "" {
		t.Fatalf("expected empty string when request id not in context, got %q", got)
	}
}
