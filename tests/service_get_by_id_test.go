package tests

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/services"
	"strings"
	"testing"
	"time"
)

func TestGetByIdSuccess(t *testing.T) {
	id := uuid.New().String()
	handler := services.ServiceCallsHandler{
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any

				m = append(m, map[string]any{
					"id":          id,
					"name":        "service",
					"description": "test desc",
					"type":        "service",
					"createdDate": time.Now(),
				})

				return m
			},
			Err: nil,
		},
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return id, true
		},
	}

	req, err := http.NewRequest("GET", "/services/"+id, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.GetById(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Check that response contains expected data
	responseBody := rw.Body.String()
	if !strings.Contains(responseBody, id) {
		t.Errorf("Response does not contain service ID: %s", responseBody)
	}
	if !strings.Contains(responseBody, "service") {
		t.Errorf("Response does not contain service name: %s", responseBody)
	}
	if !strings.Contains(responseBody, "test desc") {
		t.Errorf("Response does not contain service description: %s", responseBody)
	}
}

func TestGetByIdInvalidId(t *testing.T) {
	handler := services.ServiceCallsHandler{
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				return []map[string]any{}
			},
			Err: nil,
		},
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "", false // Invalid ID
		},
	}

	req, err := http.NewRequest("GET", "/services/invalid-id", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.GetById(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}

	expectedError := "Service id is required"
	if !strings.Contains(rw.Body.String(), expectedError) {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, rw.Body.String())
	}
}

func TestGetByIdRepositoryError(t *testing.T) {
	id := uuid.New().String()
	expectedError := "database connection error"

	handler := services.ServiceCallsHandler{
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				return []map[string]any{}
			},
			Err: errors.New(expectedError),
		},
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return id, true
		},
	}

	req, err := http.NewRequest("GET", "/services/"+id, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.GetById(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}

	if !strings.Contains(rw.Body.String(), expectedError) {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, rw.Body.String())
	}
}

func TestGetByIdServiceNotFound(t *testing.T) {
	id := uuid.New().String()
	nonExistentId := uuid.New().String()

	handler := services.ServiceCallsHandler{
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any

				m = append(m, map[string]any{
					"id":          id, // Different ID than the one we'll request
					"name":        "service",
					"description": "test desc",
					"type":        "service",
					"createdDate": time.Now(),
				})

				return m
			},
			Err: nil,
		},
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return nonExistentId, true // Valid but non-existent ID
		},
	}

	req, err := http.NewRequest("GET", "/services/"+nonExistentId, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.GetById(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}

	expectedError := "Service not found"
	if !strings.Contains(rw.Body.String(), expectedError) {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, rw.Body.String())
	}
}
