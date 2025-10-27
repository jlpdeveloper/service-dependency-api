package teams

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteTeamSuccess(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("DELETE", "/teams/be00abbc-42c6-47aa-a45a-e4e02cb6363f", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "be00abbc-42c6-47aa-a45a-e4e02cb6363f")
	rw := httptest.NewRecorder()

	handler.DeleteTeam(rw, req)

	if rw.Code != http.StatusNoContent {
		t.Errorf("DeleteTeam returned wrong status code: got %v want %v", rw.Code, http.StatusNoContent)
	}
	if rw.Body.String() != "" {
		t.Errorf("DeleteTeam returned unexpected body: got %v want empty string", rw.Body.String())
	}
}

func TestDeleteTeamInvalidId(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{}}

	req, err := http.NewRequest("DELETE", "/teams/invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	rw := httptest.NewRecorder()

	handler.DeleteTeam(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("DeleteTeam returned wrong status code: got %v want %v", rw.Code, http.StatusBadRequest)
	}
}

func TestDeleteTeamRepositoryError(t *testing.T) {
	handler := CallsHandler{Repository: mockTeamRepository{Err: errors.New("repo error")}}

	req, err := http.NewRequest("DELETE", "/teams/be00abbc-42c6-47aa-a45a-e4e02cb6363f", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", "be00abbc-42c6-47aa-a45a-e4e02cb6363f")
	rw := httptest.NewRecorder()

	handler.DeleteTeam(rw, req)

	if rw.Code != http.StatusInternalServerError {
		t.Errorf("DeleteTeam returned wrong status code: got %v want %v", rw.Code, http.StatusInternalServerError)
	}
}
