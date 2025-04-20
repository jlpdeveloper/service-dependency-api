package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteServiceSuccess(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "1", true // Return valid ID
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id":          "1",
					"name":        "ExistingService",
					"type":        "Internal",
					"description": "Existing service description",
				})
				return m
			},
			Err: nil,
		},
	}

	req, err := http.NewRequest("DELETE", "/services/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.DeleteServiceById(rw, req)

	if rw.Code != http.StatusNoContent {
		t.Errorf("DeleteServiceById returned wrong status code: got %v want %v", rw.Code, http.StatusNoContent)
	}

	// Check that the response body is empty
	if rw.Body.String() != "" {
		t.Errorf("DeleteServiceById returned unexpected body: got %v want empty string", rw.Body.String())
	}
}

func TestDeleteServiceInvalidId(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "", false // Return error for invalid ID
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id": "1",
				})
				return m
			},
			Err: nil,
		},
	}

	req, err := http.NewRequest("DELETE", "/services/invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.DeleteServiceById(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("DeleteServiceById returned wrong status code: got %v want %v", rw.Code, http.StatusBadRequest)
	}
}

func TestDeleteServiceError(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "1", true
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id": "1",
				})
				return m
			},
			Err: errors.New("test error"),
		},
	}

	req, err := http.NewRequest("DELETE", "/services/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.DeleteServiceById(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("DeleteServiceById returned wrong status code: got %v want %v", rw.Code, http.StatusInternalServerError)
	}
}

func TestDeleteNonExistentService(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "999", true // Return valid but non-existent ID
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id": "1", // Only service with ID 1 exists
				})
				return m
			},
			Err: nil,
		},
	}

	req, err := http.NewRequest("DELETE", "/services/999", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.DeleteServiceById(rw, req)

	// Delete should be idempotent, so deleting a non-existent service should still return 204
	if rw.Code != http.StatusNoContent {
		t.Errorf("DeleteServiceById returned wrong status code: got %v want %v", rw.Code, http.StatusNoContent)
	}
}
