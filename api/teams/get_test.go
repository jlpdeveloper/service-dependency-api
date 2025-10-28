package teams

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"service-dependency-api/repositories"

	"github.com/google/uuid"
)

func TestGetTeamSuccess(t *testing.T) {
	id := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{team: repositories.Team{Id: id, Name: "Platform Team"}}}

	req, err := http.NewRequest("GET", "/teams/"+id, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", id)
	rw := httptest.NewRecorder()

	h.GetTeam(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rw.Code)
	}
	if ct := rw.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	body := rw.Body.String()
	if !strings.Contains(body, id) {
		t.Errorf("response body does not contain team id: %s", body)
	}
	if !strings.Contains(body, "Platform Team") {
		t.Errorf("response body does not contain team name: %s", body)
	}
}

func TestGetTeamInvalidId(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("GET", "/teams/not-a-guid", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	// Do not set path value or set invalid one to trigger validation failure
	rw := httptest.NewRecorder()

	h.GetTeam(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if !strings.Contains(rw.Body.String(), "Invalid team ID") {
		t.Errorf("expected error message for invalid team id, got %q", rw.Body.String())
	}
}

func TestGetTeamRepositoryError(t *testing.T) {
	id := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{Err: errors.New("repo error")}}

	req, err := http.NewRequest("GET", "/teams/"+id, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", id)
	rw := httptest.NewRecorder()

	h.GetTeam(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
	if !strings.Contains(rw.Body.String(), "repo error") {
		t.Errorf("expected repo error message, got %q", rw.Body.String())
	}
}
