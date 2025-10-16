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
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err:              nil,  // No error
			DependencyExists: true, // Dependency exists
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/be00abbc-42c6-47aa-a45a-e4e02cb6363f/dependencies/884447b6-7fa1-4f6c-b684-7528fe08883d", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "be00abbc-42c6-47aa-a45a-e4e02cb6363f")
	req.SetPathValue("id2", "884447b6-7fa1-4f6c-b684-7528fe08883d")
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.DeleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rw.Code)
	}
}

func TestDeleteDependencyInvalidServiceId(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err:              nil,  // No error
			DependencyExists: true, // Dependency exists
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/invalid-id/dependencies/884447b6-7fa1-4f6c-b684-7528fe08883d", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "invalid-id")
	req.SetPathValue("id2", "884447b6-7fa1-4f6c-b684-7528fe08883d")
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.DeleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestDeleteDependencyInvalidDependencyId(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err:              nil,  // No error
			DependencyExists: true, // Dependency exists
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/be00abbc-42c6-47aa-a45a-e4e02cb6363f/dependencies/invalid-id", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "be00abbc-42c6-47aa-a45a-e4e02cb6363f")
	req.SetPathValue("id2", "invalid-id")
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.DeleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestDeleteDependencyRepositoryError(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockDependencyRepository{
			Data: func() []map[string]any {
				return []map[string]any{} // Empty data, not used in this test
			},
			Err:              errors.New("repository error"), // Simulate a repository error
			DependencyExists: true,                           // Dependency exists
		},
	}

	// Create a request
	req, err := http.NewRequest("DELETE", "/services/be00abbc-42c6-47aa-a45a-e4e02cb6363f/dependencies/be00abbc-42c6-47aa-a45a-e4e02cb6363fbe00abbc-42c6-47aa-a45a-e4e02cb6363f", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "be00abbc-42c6-47aa-a45a-e4e02cb6363f")
	req.SetPathValue("id2", "884447b6-7fa1-4f6c-b684-7528fe08883d")
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.DeleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestDeleteDependencyNotFound(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
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
	req, err := http.NewRequest("DELETE", "/services/be00abbc-42c6-47aa-a45a-e4e02cb6363f/dependencies/884447b6-7fa1-4f6c-b684-7528fe08883d", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "be00abbc-42c6-47aa-a45a-e4e02cb6363f")
	req.SetPathValue("id2", "884447b6-7fa1-4f6c-b684-7528fe08883d")
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.DeleteDependency(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}
