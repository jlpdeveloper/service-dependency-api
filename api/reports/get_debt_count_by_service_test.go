package reports

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"service-atlas/repositories"
	"testing"
)

func TestGetServiceDebtReportSuccess(t *testing.T) {
	// Arrange
	expected := []repositories.ServiceDebtReport{
		{Name: "svc-a", Id: "11111111-1111-1111-1111-111111111111", Count: 3},
		{Name: "svc-b", Id: "22222222-2222-2222-2222-222222222222", Count: 2},
	}

	handler := CallsHandler{
		repository: mockReportRepository{
			Debt: expected,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/reports/debt/services", nil)
	rw := httptest.NewRecorder()

	// Act
	handler.GetServiceDebtReport(rw, req)

	// Assert
	if rw.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rw.Code)
	}
	if ct := rw.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	var got []repositories.ServiceDebtReport
	if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
		t.Fatalf("failed decoding response: %v", err)
	}
	if len(got) != len(expected) {
		t.Fatalf("expected %d items, got %d", len(expected), len(got))
	}
	// simple element-wise comparison (order preserved by mock)
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("at %d: expected %+v, got %+v", i, expected[i], got[i])
		}
	}
}

func TestGetServiceDebtReportRepositoryError(t *testing.T) {
	// Arrange
	handler := CallsHandler{
		repository: mockReportRepository{Err: errors.New("boom")},
	}
	req := httptest.NewRequest(http.MethodGet, "/reports/debt/services", nil)
	rw := httptest.NewRecorder()

	// Act
	handler.GetServiceDebtReport(rw, req)

	// Assert
	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestGetServiceDebtReportEmptyResult(t *testing.T) {
	// Arrange: return empty slice (should encode as [])
	handler := CallsHandler{
		repository: mockReportRepository{Debt: []repositories.ServiceDebtReport{}},
	}
	req := httptest.NewRequest(http.MethodGet, "/reports/debt/services", nil)
	rw := httptest.NewRecorder()

	// Act
	handler.GetServiceDebtReport(rw, req)

	// Assert
	if rw.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rw.Code)
	}
	var got []repositories.ServiceDebtReport
	if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
		t.Fatalf("failed decoding response: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty array, got %v", got)
	}
}
