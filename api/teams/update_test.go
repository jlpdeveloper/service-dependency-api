package teams

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

func TestUpdateTeamSuccess(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	team := repositories.Team{
		Id:   "be00abbc-42c6-47aa-a45a-e4e02cb6363f",
		Name: "Platform Team",
	}
	payload, err := json.Marshal(&team)
	if err != nil {
		t.Fatalf("failed to marshal team: %v", err)
	}

	req, err := http.NewRequest("PUT", "/teams/"+team.Id, io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", team.Id)
	rw := httptest.NewRecorder()

	handler.UpdateTeam(rw, req)

	if rw.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rw.Code)
	}
}

func TestUpdateTeamRepositoryError(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{Err: errors.New("repo error")}}

	team := repositories.Team{
		Id:   "be00abbc-42c6-47aa-a45a-e4e02cb6363f",
		Name: "Platform Team",
	}
	payload, _ := json.Marshal(&team)
	req, err := http.NewRequest("PUT", "/teams/"+team.Id, io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", team.Id)
	rw := httptest.NewRecorder()

	handler.UpdateTeam(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestUpdateTeamInvalidBody(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("PUT", "/teams/be00abbc-42c6-47aa-a45a-e4e02cb6363f", io.NopCloser(strings.NewReader("not json")))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", "be00abbc-42c6-47aa-a45a-e4e02cb6363f")
	rw := httptest.NewRecorder()

	handler.UpdateTeam(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestUpdateTeamInvalidPathId(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("PUT", "/teams/1", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", "1") // invalid UUID
	rw := httptest.NewRecorder()

	handler.UpdateTeam(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestUpdateTeamIdMismatch(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	team := repositories.Team{
		Id:   "be00abbc-42c6-47aa-a45a-e4e02cb6364f", // different id than path
		Name: "Platform Team",
	}
	payload, _ := json.Marshal(&team)
	req, err := http.NewRequest("PUT", "/teams/be00abbc-42c6-47aa-a45a-e4e02cb6363f", io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", "be00abbc-42c6-47aa-a45a-e4e02cb6363f")
	rw := httptest.NewRecorder()

	handler.UpdateTeam(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestUpdateTeamValidationError(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	// Missing name should cause validation error
	team := repositories.Team{
		Id: "be00abbc-42c6-47aa-a45a-e4e02cb6363f",
	}
	payload, _ := json.Marshal(&team)
	req, err := http.NewRequest("PUT", "/teams/"+team.Id, io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", team.Id)
	rw := httptest.NewRecorder()

	handler.UpdateTeam(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestUpdateTeamNotFound(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{Err: &customerrors.HTTPError{Status: http.StatusNotFound, Msg: "not found"}}}

	team := repositories.Team{
		Id:   "be00abbc-42c6-47aa-a45a-e4e02cb6363f",
		Name: "Platform Team",
	}
	payload, _ := json.Marshal(&team)
	req, err := http.NewRequest("PUT", "/teams/"+team.Id, io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("id", team.Id)
	rw := httptest.NewRecorder()

	handler.UpdateTeam(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, rw.Code)
	}
}
