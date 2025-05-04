package dependencies

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/internal/customErrors"
	"testing"
)

func TestDeleteDependencySuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(param string, _ *http.Request) (string, bool) {
			// Return valid IDs for both path parameters
			if param == "id" {
				return "service-id-123", true
			}
			if param == "id2" {
				return "dependency-id-456", true
			}
			return "", false
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err:              nil,  // No error
			DependencyExists: true, // Dependency exists
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/service-id-123/dependencies/dependency-id-456", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.deleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rw.Code)
	}
}

func TestDeleteDependencyInvalidServiceId(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(param string, _ *http.Request) (string, bool) {
			// Return invalid service ID
			if param == "id" {
				return "", false
			}
			if param == "id2" {
				return "dependency-id-456", true
			}
			return "", false
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err:              nil,  // No error
			DependencyExists: true, // Dependency exists
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/invalid-id/dependencies/dependency-id-456", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.deleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestDeleteDependencyInvalidDependencyId(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(param string, _ *http.Request) (string, bool) {
			// Return valid service ID but invalid dependency ID
			if param == "id" {
				return "service-id-123", true
			}
			if param == "id2" {
				return "", false
			}
			return "", false
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err:              nil,  // No error
			DependencyExists: true, // Dependency exists
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/service-id-123/dependencies/invalid-id", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.deleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestDeleteDependencyRepositoryError(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(param string, _ *http.Request) (string, bool) {
			// Return valid IDs for both path parameters
			if param == "id" {
				return "service-id-123", true
			}
			if param == "id2" {
				return "dependency-id-456", true
			}
			return "", false
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err:              errors.New("repository error"), // Simulate a repository error
			DependencyExists: true,                           // Dependency exists
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/service-id-123/dependencies/dependency-id-456", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.deleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestDeleteDependencyNotFound(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		PathValidator: func(param string, _ *http.Request) (string, bool) {
			// Return valid IDs for both path parameters
			if param == "id" {
				return "service-id-123", true
			}
			if param == "id2" {
				return "dependency-id-456", true
			}
			return "", false
		},
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err: &customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Dependency relationship not found",
			}, // Simulate a 404 error
			DependencyExists: false, // Dependency doesn't exist
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/service-id-123/dependencies/dependency-id-456", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.deleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}
