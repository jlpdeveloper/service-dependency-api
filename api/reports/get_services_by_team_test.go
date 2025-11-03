package reports

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
	"testing"
)

func TestGetServicesByTeamSuccess(t *testing.T) {
	validTeamId := "123e4567-e89b-12d3-a456-426614174000"

	services := []repositories.Service{
		{Id: "1", Name: "svc-a", ServiceType: "service", Description: "desc", Url: "https://a.example.com"},
		{Id: "2", Name: "svc-b", ServiceType: "service", Description: "desc", Url: "https://b.example.com"},
	}

	h := CallsHandler{repository: mockReportRepository{Services: services}}

	req, err := http.NewRequest("GET", "/reports/teams/"+validTeamId+"/services", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("teamId", validTeamId)

	rw := httptest.NewRecorder()
	h.GetServicesByTeam(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rw.Code)
	}

	var got []repositories.Service
	if err := json.NewDecoder(rw.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if len(got) != len(services) {
		t.Fatalf("expected %d services, got %d", len(services), len(got))
	}
	for i := range services {
		if got[i].Name != services[i].Name {
			t.Errorf("service[%d].Name = %q, want %q", i, got[i].Name, services[i].Name)
		}
	}
}

func TestGetServicesByTeamInvalidTeamId(t *testing.T) {
	h := CallsHandler{repository: mockReportRepository{}}

	req, err := http.NewRequest("GET", "/reports/teams/invalid-id/services", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("teamId", "invalid-id")

	rw := httptest.NewRecorder()
	h.GetServicesByTeam(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetServicesByTeamRepositoryError(t *testing.T) {
	validTeamId := "123e4567-e89b-12d3-a456-426614174000"
	h := CallsHandler{repository: mockReportRepository{Err: errors.New("repo error")}}

	req, err := http.NewRequest("GET", "/reports/teams/"+validTeamId+"/services", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("teamId", validTeamId)

	rw := httptest.NewRecorder()
	h.GetServicesByTeam(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestGetServicesByTeamHTTPError(t *testing.T) {
	validTeamId := "123e4567-e89b-12d3-a456-426614174000"
	h := CallsHandler{repository: mockReportRepository{Err: &customErrors.HTTPError{Status: http.StatusNotFound, Msg: "team not found"}}}

	req, err := http.NewRequest("GET", "/reports/teams/"+validTeamId+"/services", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("teamId", validTeamId)

	rw := httptest.NewRecorder()
	h.GetServicesByTeam(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rw.Code)
	}
}
