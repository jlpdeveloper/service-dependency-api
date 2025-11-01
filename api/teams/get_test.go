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

func TestGetTeamsSuccess(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{teams: []repositories.Team{{Id: "1", Name: "Platform"}, {Id: "2", Name: "Payments"}}}}

	req, err := http.NewRequest("GET", "/teams?page=1&pageSize=2", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rw.Code)
	}
	if ct := rw.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	body := rw.Body.String()
	if !strings.Contains(body, "Platform") || !strings.Contains(body, "Payments") {
		t.Errorf("response body missing expected team names: %s", body)
	}
}

func TestGetTeamsInvalidPage(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("GET", "/teams?page=abc&pageSize=5", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "invalid syntax") {
		t.Errorf("expected invalid syntax error in body, got %q", body)
	}
}

func TestGetTeamsRepositoryError(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{Err: errors.New("repo error")}}

	req, err := http.NewRequest("GET", "/teams?page=1&pageSize=10", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
	if !strings.Contains(rw.Body.String(), "repo error") {
		t.Errorf("expected repo error message, got %q", rw.Body.String())
	}
}

func TestGetTeamsDefaultPageSize(t *testing.T) {
	// pageSize is non-numeric -> defaults to 10; handler should still succeed
	h := CallsHandler{Repository: mockTeamRepository{teams: []repositories.Team{{Id: "1", Name: "A"}}}}

	req, err := http.NewRequest("GET", "/teams?page=1&pageSize=not-a-number", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rw.Code)
	}
	if ct := rw.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", ct)
	}
	if body := rw.Body.String(); !strings.Contains(body, "\"A\"") {
		t.Errorf("expected team in response, got %q", body)
	}
}

func TestGetTeamsPageSizeTooSmall(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("GET", "/teams?page=1&pageSize=0", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "pageSize must be between 1 and 100") {
		t.Errorf("expected pageSize bounds error, got %q", body)
	}
}

func TestGetTeamsPageSizeTooLarge(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("GET", "/teams?page=1&pageSize=101", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "pageSize must be between 1 and 100") {
		t.Errorf("expected pageSize bounds error, got %q", body)
	}
}

func TestGetTeamsPageZero(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("GET", "/teams?page=0&pageSize=10", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "page must be positive") {
		t.Errorf("expected page positivity error, got %q", body)
	}
}

func TestGetTeamsPageNegative(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("GET", "/teams?page=-1&pageSize=10", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "page must be positive") {
		t.Errorf("expected page positivity error, got %q", body)
	}
}

func TestGetTeamsMissingPageParameter(t *testing.T) {
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("GET", "/teams?pageSize=10", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	h.GetTeams(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "invalid syntax") {
		t.Errorf("expected invalid syntax error in body for missing page parameter, got %q", body)
	}
}
