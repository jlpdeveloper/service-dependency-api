package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"service-atlas/repositories"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// Tests for GetTeamsByServiceId handler

func TestGetTeamsByServiceId_Success(t *testing.T) {
	serviceID := uuid.New().String()

	handler := ServiceCallsHandler{
		Repository: mockServiceRepository{
			Data: func() []map[string]any {
				return []map[string]any{
					{
						// mock maps are interpreted by mockServiceRepository.GetTeamsByServiceId
						// It uses teamId and teamName keys
						"teamId":    "team123",
						"teamName":  "Team One",
						"serviceId": serviceID,
					},
				}
			},
			Err: nil,
		},
	}

	req, err := http.NewRequest("GET", "/services/"+serviceID+"/teams", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", serviceID)

	rw := httptest.NewRecorder()
	handler.GetTeamsByServiceId(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rw.Code)
	}

	if ct := rw.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("expected Content-Type application/json, got %s", ct)
	}

	var teams []repositories.Team
	if err := json.NewDecoder(rw.Body).Decode(&teams); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(teams) != 1 {
		t.Fatalf("expected 1 team, got %d", len(teams))
	}
	if teams[0].Id != "team123" {
		t.Fatalf("unexpected team id: %s", teams[0].Id)
	}
	if teams[0].Name != "Team One" {
		t.Fatalf("unexpected team name: %s", teams[0].Name)
	}
}

func TestGetTeamsByServiceId_BadRequest(t *testing.T) {
	handler := ServiceCallsHandler{Repository: mockServiceRepository{Data: func() []map[string]any { return nil }}}

	req, err := http.NewRequest("GET", "/services/invalid/teams", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rw := httptest.NewRecorder()
	handler.GetTeamsByServiceId(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if !strings.Contains(rw.Body.String(), "Invalid service ID") {
		t.Fatalf("expected error message for invalid id, got: %s", rw.Body.String())
	}
}

func TestGetTeamsByServiceId_RepoError(t *testing.T) {
	serviceID := uuid.New().String()
	expectedError := "repository failure"

	handler := ServiceCallsHandler{
		Repository: mockServiceRepository{
			Data: func() []map[string]any { return nil },
			Err:  errors.New(expectedError),
		},
	}

	req, err := http.NewRequest("GET", "/services/"+serviceID+"/teams", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", serviceID)

	rw := httptest.NewRecorder()
	handler.GetTeamsByServiceId(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
	if !strings.Contains(rw.Body.String(), expectedError) {
		t.Fatalf("expected error message '%s', got '%s'", expectedError, rw.Body.String())
	}
}

func TestGetTeamsByServiceId_EmptyList(t *testing.T) {
	serviceID := uuid.New().String()

	handler := ServiceCallsHandler{
		Repository: mockServiceRepository{
			Data: func() []map[string]any { return []map[string]any{} },
			Err:  nil,
		},
	}

	req, err := http.NewRequest("GET", "/services/"+serviceID+"/teams", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", serviceID)

	rw := httptest.NewRecorder()
	handler.GetTeamsByServiceId(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rw.Code)
	}

	// Expect JSON array (possibly empty)
	var teams []repositories.Team
	if err := json.NewDecoder(rw.Body).Decode(&teams); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(teams) != 0 {
		t.Fatalf("expected 0 teams, got %d", len(teams))
	}
}
