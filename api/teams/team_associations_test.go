package teams

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCreateTeamAssociationSuccess(t *testing.T) {
	teamID := uuid.New().String()
	serviceID := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("POST", "/teams/"+teamID+"/services/"+serviceID, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("teamId", teamID)
	req.SetPathValue("serviceId", serviceID)
	rw := httptest.NewRecorder()

	h.CreateTeamAssociation(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rw.Code)
	}
}

func TestCreateTeamAssociationInvalidTeamID(t *testing.T) {
	serviceID := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("POST", "/teams/not-a-guid/services/"+serviceID, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	// Do not set path value or set invalid one to trigger validation failure
	req.SetPathValue("serviceId", serviceID)
	rw := httptest.NewRecorder()

	h.CreateTeamAssociation(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "Invalid team ID") {
		t.Errorf("expected invalid team id error, got %q", body)
	}
}

func TestCreateTeamAssociationInvalidServiceID(t *testing.T) {
	teamID := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("POST", "/teams/"+teamID+"/services/not-a-guid", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("teamId", teamID)
	// Do not set valid service id to trigger validation failure
	rw := httptest.NewRecorder()

	h.CreateTeamAssociation(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "Invalid service ID") {
		t.Errorf("expected invalid service id error, got %q", body)
	}
}

func TestCreateTeamAssociationRepositoryError(t *testing.T) {
	teamID := uuid.New().String()
	serviceID := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{Err: errors.New("repo error")}}

	req, err := http.NewRequest("POST", "/teams/"+teamID+"/services/"+serviceID, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("teamId", teamID)
	req.SetPathValue("serviceId", serviceID)
	rw := httptest.NewRecorder()

	h.CreateTeamAssociation(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "repo error") {
		t.Errorf("expected repo error in body, got %q", body)
	}
}

func TestDeleteTeamAssociationSuccess(t *testing.T) {
	teamID := uuid.New().String()
	serviceID := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("DELETE", "/teams/"+teamID+"/services/"+serviceID, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("teamId", teamID)
	req.SetPathValue("serviceId", serviceID)
	rw := httptest.NewRecorder()

	h.DeleteTeamAssociation(rw, req)

	if rw.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rw.Code)
	}
}

func TestDeleteTeamAssociationInvalidTeamID(t *testing.T) {
	serviceID := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("DELETE", "/teams/not-a-guid/services/"+serviceID, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	// Only set serviceId; omit teamId to trigger validation error
	req.SetPathValue("serviceId", serviceID)
	rw := httptest.NewRecorder()

	h.DeleteTeamAssociation(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "Invalid team ID") {
		t.Errorf("expected invalid team id error, got %q", body)
	}
}

func TestDeleteTeamAssociationInvalidServiceID(t *testing.T) {
	teamID := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("DELETE", "/teams/"+teamID+"/services/not-a-guid", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("teamId", teamID)
	rw := httptest.NewRecorder()

	h.DeleteTeamAssociation(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "Invalid service ID") {
		t.Errorf("expected invalid service id error, got %q", body)
	}
}

func TestDeleteTeamAssociationRepositoryError(t *testing.T) {
	teamID := uuid.New().String()
	serviceID := uuid.New().String()
	h := CallsHandler{Repository: mockTeamRepository{Err: errors.New("repo error")}}

	req, err := http.NewRequest("DELETE", "/teams/"+teamID+"/services/"+serviceID, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.SetPathValue("teamId", teamID)
	req.SetPathValue("serviceId", serviceID)
	rw := httptest.NewRecorder()

	h.DeleteTeamAssociation(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rw.Code)
	}
	if body := rw.Body.String(); !strings.Contains(body, "repo error") {
		t.Errorf("expected repo error in body, got %q", body)
	}
}
