package debt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
	"testing"
)

func TestGetDebtByServiceIdSuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	mockDebts := []repositories.Debt{
		{
			ServiceId:   validServiceId,
			Type:        "code",
			Title:       "Test Debt 1",
			Description: "This is a test debt 1",
			Status:      "pending",
		},
		{
			ServiceId:   validServiceId,
			Type:        "documentation",
			Title:       "Test Debt 2",
			Description: "This is a test debt 2",
			Status:      "resolved",
		},
	}

	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err:   nil, // No error
			Debts: mockDebts,
		},
	}

	// Create a request with the path pattern
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/debt", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetDebtByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Decode the response body
	var responseDebts []repositories.Debt
	err = json.NewDecoder(rw.Body).Decode(&responseDebts)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check the response body
	if len(responseDebts) != len(mockDebts) {
		t.Errorf("Expected %d debts, got %d", len(mockDebts), len(responseDebts))
	}

	// Check each debt
	for i, debt := range responseDebts {
		if debt.ServiceId != mockDebts[i].ServiceId {
			t.Errorf("Expected ServiceId %s, got %s", mockDebts[i].ServiceId, debt.ServiceId)
		}
		if debt.Type != mockDebts[i].Type {
			t.Errorf("Expected Type %s, got %s", mockDebts[i].Type, debt.Type)
		}
		if debt.Title != mockDebts[i].Title {
			t.Errorf("Expected Title %s, got %s", mockDebts[i].Title, debt.Title)
		}
		if debt.Description != mockDebts[i].Description {
			t.Errorf("Expected Description %s, got %s", mockDebts[i].Description, debt.Description)
		}
		if debt.Status != mockDebts[i].Status {
			t.Errorf("Expected Status %s, got %s", mockDebts[i].Status, debt.Status)
		}
	}
}

func TestGetDebtByServiceIdEmptyResult(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err:   nil,                   // No error
			Debts: []repositories.Debt{}, // Empty slice
		},
	}

	// Create a request with the path pattern
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/debt", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetDebtByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Decode the response body
	var responseDebts []repositories.Debt
	err = json.NewDecoder(rw.Body).Decode(&responseDebts)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check the response body
	if len(responseDebts) != 0 {
		t.Errorf("Expected 0 debts, got %d", len(responseDebts))
	}
}

func TestGetDebtByServiceIdInvalidPathParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid service ID
	req, err := http.NewRequest("GET", "/services/invalid-id/debt", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "invalid-id")

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetDebtByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetDebtByServiceIdRepositoryError(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: errors.New("repository error"), // Simulate a repository error
		},
	}

	// Create a request with the path pattern
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/debt", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetDebtByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestGetDebtByServiceIdHTTPError(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := CallsHandler{
		Repository: mockDebtRepository{
			Err: &customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Service not found",
			}, // Simulate an HTTP error
		},
	}

	// Create a request with the path pattern
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/debt", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetDebtByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}

func TestGetDebtByServiceId_OnlyResolved(t *testing.T) {
	ctx := context.Background()
	mock := &mockDebtRepository{
		Debts: []repositories.Debt{
			{ServiceId: "svc1", Title: "A", Status: "remediated"},
			{ServiceId: "svc1", Title: "B", Status: "pending"},
			{ServiceId: "svc1", Title: "C", Status: "remediated"},
			{ServiceId: "svc2", Title: "D", Status: "in_progress"},
		},
	}

	type args struct {
		serviceId    string
		onlyResolved bool
		wantCount    int
		wantTitles   []string
	}
	tests := []args{
		{serviceId: "svc1", onlyResolved: false, wantCount: 3, wantTitles: []string{"A", "B", "C"}},
		{serviceId: "svc1", onlyResolved: true, wantCount: 2, wantTitles: []string{"A", "C"}},
		{serviceId: "svc2", onlyResolved: true, wantCount: 0, wantTitles: nil},
	}

	for _, tt := range tests {
		got, err := mock.GetDebtByServiceId(ctx, tt.serviceId, 0, 10, tt.onlyResolved)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != tt.wantCount {
			t.Errorf("GetDebtByServiceId(%q, onlyResolved=%v) got %d debts, want %d",
				tt.serviceId, tt.onlyResolved, len(got), tt.wantCount)
		}
		for i, d := range got {
			if i < len(tt.wantTitles) && d.Title != tt.wantTitles[i] {
				t.Errorf("debt[%d].Title = %q, want %q", i, d.Title, tt.wantTitles[i])
			}
		}
	}
}
