package releases

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/releases/internal/releaseRepository"
	"service-dependency-api/internal/customErrors"
	"testing"
	"time"
)

func TestGetReleasesByServiceIdSuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	mockReleases := []*releaseRepository.Release{
		{
			ServiceId:   validServiceId,
			ReleaseDate: time.Now().UTC(),
			Url:         "https://example.com/release1",
			Version:     "1.0.0",
		},
		{
			ServiceId:   validServiceId,
			ReleaseDate: time.Now().UTC(),
			Url:         "https://example.com/release2",
			Version:     "2.0.0",
		},
	}

	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err:      nil, // No error
			Releases: mockReleases,
		},
		PathValidator: func(name string, req *http.Request) (string, bool) {
			return validServiceId, true // Mock successful path validation
		},
	}

	// Create a request with no pagination parameters (should use defaults)
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/releases", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Decode the response body
	var releases []*releaseRepository.Release
	err = json.NewDecoder(rw.Body).Decode(&releases)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check that the correct number of releases was returned
	if len(releases) != len(mockReleases) {
		t.Errorf("Expected %d releases, got %d", len(mockReleases), len(releases))
	}
}

func TestGetReleasesByServiceIdInvalidPathParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
		PathValidator: func(name string, req *http.Request) (string, bool) {
			return "invalid-id", false // Mock failed path validation
		},
	}

	// Create a request
	req, err := http.NewRequest("GET", "/services/invalid-id/releases", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesByServiceIdRepositoryError(t *testing.T) {
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

	// Create a request
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/releases", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestGetReleasesByServiceIdHTTPError(t *testing.T) {
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

	// Create a request
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/releases", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}

func TestGetReleasesByServiceIdWithPagination(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID

	// Create 30 mock releases
	mockReleases := make([]*releaseRepository.Release, 30)
	for i := 0; i < 30; i++ {
		mockReleases[i] = &releaseRepository.Release{
			ServiceId:   validServiceId,
			ReleaseDate: time.Now().UTC(),
			Url:         fmt.Sprintf("https://example.com/release%d", i+1),
			Version:     fmt.Sprintf("%d.0.0", i+1),
		}
	}

	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err:      nil, // No error
			Releases: mockReleases,
		},
		PathValidator: func(name string, req *http.Request) (string, bool) {
			return validServiceId, true // Mock successful path validation
		},
	}

	// Create a request with pagination parameters (page=1, page_size=10)
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/releases?page=2&page_size=10", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Decode the response body
	var releases []*releaseRepository.Release
	err = json.NewDecoder(rw.Body).Decode(&releases)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check that the correct number of releases was returned (should be 10)
	if len(releases) != 10 {
		t.Errorf("Expected %d releases, got %d", 10, len(releases))
	}

	// Check that the correct page of releases was returned (should be releases 10-19)
	for i, release := range releases {
		expectedVersion := fmt.Sprintf("%d.0.0", i+11) // Page 1 starts at index 10
		if release.Version != expectedVersion {
			t.Errorf("Expected version %s, got %s", expectedVersion, release.Version)
		}
	}
}

func TestGetReleasesByServiceIdInvalidPageParameter(t *testing.T) {
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

	// Create a request with invalid page parameter
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/releases?page=invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesByServiceIdInvalidPageSizeParameter(t *testing.T) {
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

	// Create a request with invalid page_size parameter
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/releases?page_size=invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesByServiceIdNegativePageParameter(t *testing.T) {
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

	// Create a request with negative page parameter
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/releases?page=-1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesByServiceIdZeroPageSizeParameter(t *testing.T) {
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

	// Create a request with zero page_size parameter
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/releases?page_size=0", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}
