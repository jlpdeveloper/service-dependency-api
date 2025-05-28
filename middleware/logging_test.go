package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	// Save the original log output and restore it after the test
	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	// Create a test handler that will be wrapped by the logging middleware
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Test different status codes
		if r.URL.Path == "/ok" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/not-found" {
			w.WriteHeader(http.StatusNotFound)
		} else if r.URL.Path == "/server-error" {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	// Wrap the test handler with the logging middleware
	loggingHandler := Logging(testHandler)

	// Test cases
	testCases := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{"OK Response", "GET", "/ok", http.StatusOK},
		{"Not Found Response", "POST", "/not-found", http.StatusNotFound},
		{"Server Error Response", "PUT", "/server-error", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear the buffer before each test case
			buf.Reset()

			// Create a test request
			req := httptest.NewRequest(tc.method, tc.path, nil)

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Serve the request using our logging handler
			loggingHandler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tc.expectedStatus)
			}

			// Check that the log contains the expected information
			logOutput := buf.String()
			if !strings.Contains(logOutput, tc.method) {
				t.Errorf("log does not contain request method: %s", tc.method)
			}
			if !strings.Contains(logOutput, tc.path) {
				t.Errorf("log does not contain request path: %s", tc.path)
			}
			if !strings.Contains(logOutput, strconv.Itoa(tc.expectedStatus)) {
				t.Errorf("log does not contain status code: %d", tc.expectedStatus)
			}
		})
	}
}

func TestWrappedWriter(t *testing.T) {
	// Create a wrapped writer
	recorder := httptest.NewRecorder()
	wrapped := &wrappedWriter{
		ResponseWriter: recorder,
		statusCode:     http.StatusOK,
	}

	// Test WriteHeader method
	wrapped.WriteHeader(http.StatusNotFound)

	// Check that the status code was updated in the wrapped writer
	if wrapped.statusCode != http.StatusNotFound {
		t.Errorf("wrapped writer did not update status code: got %v want %v",
			wrapped.statusCode, http.StatusNotFound)
	}

	// Check that the status code was passed to the underlying ResponseWriter
	if recorder.Code != http.StatusNotFound {
		t.Errorf("status code not passed to underlying ResponseWriter: got %v want %v",
			recorder.Code, http.StatusNotFound)
	}
}
