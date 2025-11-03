package debt

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/internal/customerrors"
	"service-dependency-api/repositories"
	"strings"
	"testing"
)

func TestCreateDebtSuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a debt request (without ServiceId as it comes from the path)
	debt := &repositories.Debt{
		Type:        "code",
		Title:       "Test Debt",
		Description: "This is a test debt",
		Status:      "pending",
	}
	debtJSON, err := json.Marshal(debt)
	if err != nil {
		t.Fatalf("Failed to marshal debt: %v", err)
	}

	// Create a request with the path pattern
	req, err := http.NewRequest("POST", "/services/"+validServiceId+"/debt",
		io.NopCloser(strings.NewReader(string(debtJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.CreateDebt(rw, req)

	// Check the response
	if rw.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, rw.Code)
	}
}

func TestCreateDebtInvalidBody(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid JSON
	req, err := http.NewRequest("POST", "/services/"+validServiceId+"/debt",
		io.NopCloser(strings.NewReader("invalid json")))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.CreateDebt(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestCreateDebtInvalidPathParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a debt request
	debt := &repositories.Debt{
		Type:        "code",
		Title:       "Test Debt",
		Description: "This is a test debt",
		Status:      "pending",
	}
	debtJSON, err := json.Marshal(debt)
	if err != nil {
		t.Fatalf("Failed to marshal debt: %v", err)
	}

	// Create a request with invalid service ID
	req, err := http.NewRequest("POST", "/services/invalid-id/debt",
		io.NopCloser(strings.NewReader(string(debtJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "invalid-id")
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.CreateDebt(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestCreateDebtRepositoryError(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: errors.New("repository error"), // Simulate a repository error
		},
	}

	// Create a debt request (without ServiceId as it comes from the path)
	debt := &repositories.Debt{
		Type:        "code",
		Title:       "Test Debt",
		Description: "This is a test debt",
		Status:      "pending",
	}
	debtJSON, err := json.Marshal(debt)
	if err != nil {
		t.Fatalf("Failed to marshal debt: %v", err)
	}

	// Create a request with the path pattern
	req, err := http.NewRequest("POST", "/services/"+validServiceId+"/debt",
		io.NopCloser(strings.NewReader(string(debtJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.CreateDebt(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestCreateDebtHTTPError(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: &customerrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Service not found",
			}, // Simulate an HTTP error
		},
	}

	// Create a debt request (without ServiceId as it comes from the path)
	debt := &repositories.Debt{
		Type:        "code",
		Title:       "Test Debt",
		Description: "This is a test debt",
		Status:      "pending",
	}
	debtJSON, err := json.Marshal(debt)
	if err != nil {
		t.Fatalf("Failed to marshal debt: %v", err)
	}

	// Create a request with the path pattern
	req, err := http.NewRequest("POST", "/services/"+validServiceId+"/debt",
		io.NopCloser(strings.NewReader(string(debtJSON))))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.CreateDebt(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}
