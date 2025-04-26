package dependencies

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/dependencies/internal/dependencyRepository"
	"service-dependency-api/internal/customErrors"
	"strings"
	"testing"
)

func TestCreateDependencySuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(_ string, _ *http.Request) (string, bool) {
			return "service-id-123", true // Return a valid service ID
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err: nil, // No error
		},
	}

	// Create a dependency request
	dependency := &dependencyRepository.Dependency{
		Id:      "dependency-id-456",
		Version: "1.0.0",
	}
	dependencyJSON, err := json.Marshal(dependency)
	if err != nil {
		t.Fatalf("Failed to marshal dependency: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/services/service-id-123/dependency",
		io.NopCloser(strings.NewReader(string(dependencyJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rw.Code)
	}
}

func TestCreateDependencyInvalidPath(t *testing.T) {
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

	// Create a dependency request
	dependency := &dependencyRepository.Dependency{
		Id:      "dependency-id-456",
		Version: "1.0.0",
	}
	dependencyJSON, err := json.Marshal(dependency)
	if err != nil {
		t.Fatalf("Failed to marshal dependency: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/services/invalid-id/dependency",
		io.NopCloser(strings.NewReader(string(dependencyJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestCreateDependencyInvalidBody(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(_ string, _ *http.Request) (string, bool) {
			return "service-id-123", true // Return a valid service ID
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err: nil, // No error
		},
	}

	// Create a request with invalid JSON
	req, err := http.NewRequest("POST", "/services/service-id-123/dependency",
		io.NopCloser(strings.NewReader("invalid json")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestCreateDependencyInvalidDependency(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(_ string, _ *http.Request) (string, bool) {
			return "service-id-123", true // Return a valid service ID
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err: nil, // No error
		},
	}

	// Create an invalid dependency (missing ID)
	dependency := &dependencyRepository.Dependency{
		Version: "1.0.0",
	}
	dependencyJSON, err := json.Marshal(dependency)
	if err != nil {
		t.Fatalf("Failed to marshal dependency: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/services/service-id-123/dependency",
		io.NopCloser(strings.NewReader(string(dependencyJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestCreateDependencyRepositoryError(t *testing.T) {
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

	// Create a dependency request
	dependency := &dependencyRepository.Dependency{
		Id:      "dependency-id-456",
		Version: "1.0.0",
	}
	dependencyJSON, err := json.Marshal(dependency)
	if err != nil {
		t.Fatalf("Failed to marshal dependency: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/services/service-id-123/dependency",
		io.NopCloser(strings.NewReader(string(dependencyJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestCreateDependencyHTTPError(t *testing.T) {
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

	// Create a dependency request
	dependency := &dependencyRepository.Dependency{
		Id:      "dependency-id-456",
		Version: "1.0.0",
	}
	dependencyJSON, err := json.Marshal(dependency)
	if err != nil {
		t.Fatalf("Failed to marshal dependency: %v", err)
	}

	// Create a request
	req, err := http.NewRequest("POST", "/services/service-id-123/dependency",
		io.NopCloser(strings.NewReader(string(dependencyJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.createDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}
