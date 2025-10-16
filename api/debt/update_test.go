package debt

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/internal/customErrors"
	"testing"
)

func TestUpdateDebtStatusSuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	validDebtId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a request with the path pattern and valid status
	requestBody := map[string]string{
		"status": "remediated", // Valid status
	}
	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("PATCH", "/debt/"+validDebtId, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validDebtId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.UpdateDebtStatus(rw, req)

	// Check the response
	if rw.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rw.Code)
	}
}

func TestUpdateDebtStatusInvalidId(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid debt ID
	requestBody := map[string]string{
		"status": "remediated", // Valid status
	}
	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("PATCH", "/debt/invalid-id", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "invalid-id")

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.UpdateDebtStatus(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestUpdateDebtStatusInvalidStatus(t *testing.T) {
	// Create a handler with mocked dependencies
	validDebtId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid status
	requestBody := map[string]string{
		"status": "invalid-status", // Invalid status
	}
	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("PATCH", "/debt/"+validDebtId, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validDebtId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.UpdateDebtStatus(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestUpdateDebtStatusMissingStatus(t *testing.T) {
	// Create a handler with mocked dependencies
	validDebtId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a request with missing status
	requestBody := map[string]string{} // No status field
	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("PATCH", "/debt/"+validDebtId, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validDebtId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.UpdateDebtStatus(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestUpdateDebtStatusInvalidBody(t *testing.T) {
	// Create a handler with mocked dependencies
	validDebtId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid JSON body
	req, err := http.NewRequest("PATCH", "/debt/"+validDebtId, bytes.NewReader([]byte("invalid json")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validDebtId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.UpdateDebtStatus(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestUpdateDebtStatusRepositoryError(t *testing.T) {
	// Create a handler with mocked dependencies
	validDebtId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: errors.New("repository error"), // Simulate a repository error
		},
	}

	// Create a request with the path pattern and valid status
	requestBody := map[string]string{
		"status": "remediated", // Valid status
	}
	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("PATCH", "/debt/"+validDebtId, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validDebtId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.UpdateDebtStatus(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestUpdateDebtStatusHTTPError(t *testing.T) {
	// Create a handler with mocked dependencies
	validDebtId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: &customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Debt not found",
			}, // Simulate an HTTP error
		},
	}

	// Create a request with the path pattern and valid status
	requestBody := map[string]string{
		"status": "remediated", // Valid status
	}
	bodyBytes, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("PATCH", "/debt/"+validDebtId, bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validDebtId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.UpdateDebtStatus(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}
