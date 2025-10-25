package teams

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/repositories"
	"strconv"
	"strings"
	"testing"
)

func TestCreateTeamSuccess(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	team := repositories.Team{
		Name: "Platform Team",
	}
	payload, err := json.Marshal(&team)
	if err != nil {
		t.Fatalf("failed to marshal team: %v", err)
	}

	req, err := http.NewRequest("POST", "/teams", io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	handler.CreateTeam(req, rw)

	if rw.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %s", http.StatusCreated, strconv.Itoa(rw.Code))
	}
}

func TestCreateTeamRepositoryError(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{Err: errors.New("repo error")}}

	team := repositories.Team{Name: "Platform Team"}
	payload, _ := json.Marshal(&team)
	req, err := http.NewRequest("POST", "/teams", io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	handler.CreateTeam(req, rw)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %s", http.StatusInternalServerError, strconv.Itoa(rw.Code))
	}
}

func TestCreateTeamInvalidJSON(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("POST", "/teams", io.NopCloser(strings.NewReader("not json")))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	handler.CreateTeam(req, rw)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %s", http.StatusBadRequest, strconv.Itoa(rw.Code))
	}
}

func TestCreateTeamValidationError(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	// missing name field should trigger validation error
	payload := `{"created":"2024-01-01T00:00:00Z"}`
	req, err := http.NewRequest("POST", "/teams", io.NopCloser(strings.NewReader(payload)))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	handler.CreateTeam(req, rw)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %s", http.StatusBadRequest, strconv.Itoa(rw.Code))
	}
}
