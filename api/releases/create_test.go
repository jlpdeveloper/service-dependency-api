package releases

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/releases/internal/releaseRepository"
	"service-dependency-api/internal/customErrors"
	"strings"
	"testing"
)

func TestCreateReleaseSuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
		PathValidator: func(name string, req *http.Request) (string, bool) {
			return validServiceId, true // Mock successful path validation
		},
	}

	// Create a release request (without ServiceId as it comes from the path)
	release := &releaseRepository.Release{
		Url: "https://example.com/release",
	}
	releaseJSON, err := json.Marshal(release)
	if err != nil {
		t.Fatalf("Failed to marshal release: %v", err)
	}

	// Create a request with the new path pattern
	req, err := http.NewRequest("POST", "/services/"+validServiceId+"/release",
		io.NopCloser(strings.NewReader(string(releaseJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createRelease(rw, req)

	// Check the response
	if rw.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rw.Code)
	}
}

func TestCreateReleaseInvalidBody(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
		PathValidator: func(name string, req *http.Request) (string, bool) {
			return validServiceId, true // Mock successful path validation
		},
	}

	// Create a request with invalid JSON
	req, err := http.NewRequest("POST", "/services/"+validServiceId+"/release",
		io.NopCloser(strings.NewReader("invalid json")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createRelease(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestCreateReleaseInvalidPathParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
		PathValidator: func(name string, req *http.Request) (string, bool) {
			return "invalid-id", false // Mock failed path validation
		},
	}

	// Create a release request
	release := &releaseRepository.Release{
		Url: "https://example.com/release",
	}
	releaseJSON, err := json.Marshal(release)
	if err != nil {
		t.Fatalf("Failed to marshal release: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/services/invalid-id/release",
		io.NopCloser(strings.NewReader(string(releaseJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createRelease(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestCreateReleaseRepositoryError(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: errors.New("repository error"), // Simulate a repository error
		},
		PathValidator: func(name string, req *http.Request) (string, bool) {
			return validServiceId, true // Mock successful path validation
		},
	}

	// Create a release request (without ServiceId as it comes from the path)
	release := &releaseRepository.Release{
		Url: "https://example.com/release",
	}
	releaseJSON, err := json.Marshal(release)
	if err != nil {
		t.Fatalf("Failed to marshal release: %v", err)
	}

	// Create a request with the new path pattern
	req, err := http.NewRequest("POST", "/services/"+validServiceId+"/release",
		io.NopCloser(strings.NewReader(string(releaseJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createRelease(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestCreateReleaseHTTPError(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: &customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Service not found",
			}, // Simulate an HTTP error
		},
		PathValidator: func(name string, req *http.Request) (string, bool) {
			return validServiceId, true // Mock successful path validation
		},
	}

	// Create a release request (without ServiceId as it comes from the path)
	release := &releaseRepository.Release{
		Url: "https://example.com/release",
	}
	releaseJSON, err := json.Marshal(release)
	if err != nil {
		t.Fatalf("Failed to marshal release: %v", err)
	}

	// Create a request with the new path pattern
	req, err := http.NewRequest("POST", "/services/"+validServiceId+"/release",
		io.NopCloser(strings.NewReader(string(releaseJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createRelease(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}
