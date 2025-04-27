package dependencies

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/dependencies/internal/dependencyRepository"
	"service-dependency-api/internal/customErrors"
	"testing"
)

func TestGetByIdSuccess(t *testing.T) {
	// Create mock dependencies
	mockDeps := []map[string]any{
		{
			"id":      "dependency-id-1",
			"name":    "Dependency 1",
			"version": "1.0.0",
		},
		{
			"id":      "dependency-id-2",
			"name":    "Dependency 2",
			"version": "2.0.0",
		},
	}

	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(_ string, _ *http.Request) (string, bool) {
			return "service-id-123", true // Return a valid service ID
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return mockDeps
			},
			Err: nil, // No error
		},
	}

	// Create a request
	req, err := http.NewRequest("GET", "/services/service-id-123/dependencies", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getById(rw, req)

	// Check the response
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Check the content type
	contentType := rw.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type %s, got %s", "application/json", contentType)
	}

	// Decode the response
	var dependencies []*dependencyRepository.Dependency
	err = json.NewDecoder(rw.Body).Decode(&dependencies)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check the number of dependencies
	if len(dependencies) != len(mockDeps) {
		t.Errorf("Expected %d dependencies, got %d", len(mockDeps), len(dependencies))
	}

	// Check the dependencies
	for i, dep := range dependencies {
		if dep.Id != mockDeps[i]["id"] {
			t.Errorf("Expected dependency ID %s, got %s", mockDeps[i]["id"], dep.Id)
		}
		if dep.Name != mockDeps[i]["name"] {
			t.Errorf("Expected dependency name %s, got %s", mockDeps[i]["name"], dep.Name)
		}
		if dep.Version != mockDeps[i]["version"] {
			t.Errorf("Expected dependency version %s, got %s", mockDeps[i]["version"], dep.Version)
		}
	}
}

func TestGetByIdInvalidPath(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(_ string, _ *http.Request) (string, bool) {
			return "", false // Return an invalid service ID
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err: nil, // No error
		},
	}

	// Create a request
	req, err := http.NewRequest("GET", "/services/invalid-id/dependencies", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getById(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetByIdRepositoryError(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(_ string, _ *http.Request) (string, bool) {
			return "service-id-123", true // Return a valid service ID
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err: errors.New("repository error"), // Simulate a repository error
		},
	}

	// Create a request
	req, err := http.NewRequest("GET", "/services/service-id-123/dependencies", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getById(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestGetByIdHTTPError(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(_ string, _ *http.Request) (string, bool) {
			return "service-id-123", true // Return a valid service ID
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err: &customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Service not found",
			}, // Simulate an HTTP error
		},
	}

	// Create a request
	req, err := http.NewRequest("GET", "/services/service-id-123/dependencies", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.getById(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}
