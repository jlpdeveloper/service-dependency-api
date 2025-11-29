package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"service-atlas/repositories"
)

func TestServiceSearch_Success(t *testing.T) {
	// Arrange
	repo := mockServiceRepository{
		Data: func() []map[string]any {
			return []map[string]any{
				{
					"id":          "123",
					"name":        "find me",
					"description": "desc",
					"type":        "api",
				},
			}
		},
	}
	h := &ServiceCallsHandler{Repository: repo}

	req := httptest.NewRequest(http.MethodGet, "/services/search?query=find", nil)
	rr := httptest.NewRecorder()

	// Act
	h.Search(rr, req)

	// Assert
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var got []repositories.Service
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 service, got %d", len(got))
	}
	if got[0].Id != "123" || got[0].Name != "find me" || got[0].ServiceType != "api" {
		t.Fatalf("unexpected service payload: %+v", got[0])
	}
}

func TestServiceSearch_MissingQuery(t *testing.T) {
	h := &ServiceCallsHandler{Repository: mockServiceRepository{Data: func() []map[string]any { return nil }}}
	req := httptest.NewRequest(http.MethodGet, "/services/search", nil)
	rr := httptest.NewRecorder()

	h.Search(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestServiceSearch_RepoError(t *testing.T) {
	repo := mockServiceRepository{Data: func() []map[string]any { return nil }, Err: errors.New("boom")}
	h := &ServiceCallsHandler{Repository: repo}
	req := httptest.NewRequest(http.MethodGet, "/services/search?query=find", nil)
	rr := httptest.NewRecorder()

	h.Search(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d (%s)", rr.Code, rr.Body.String())
	}
}
