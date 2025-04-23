package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/services/serviceRepository"
	"strings"
	"testing"
)

func TestUpdateServiceSuccess(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "1", true // Return valid ID and false for no error
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id":          "1",
					"name":        "ExistingService",
					"type":        "Internal",
					"description": "Existing service description",
				})
				return m
			},
			Err: nil,
		},
	}

	// Create a service update request
	service := serviceRepository.Service{
		Id:          "1", // Must match the ID in the mock data
		Name:        "UpdatedService",
		ServiceType: "External",
		Description: "Updated service description",
	}
	serviceJson, err := json.Marshal(&service)
	req, err := http.NewRequest("PUT", "/services/1", io.NopCloser(strings.NewReader(string(serviceJson))))

	rw := httptest.NewRecorder()
	handler.UpdateService(rw, req)
	if err != nil {
		t.Errorf("Service UPDATE errored with %s", err.Error())
	}
	if rw.Code != http.StatusNoContent {
		t.Errorf("Service UPDATE returned wrong status code: got %v want %v", rw.Code, http.StatusNoContent)
	}
}

func TestUpdateServiceNotFound(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "999", true // Return non-existent ID
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id":          "1",
					"name":        "ExistingService",
					"type":        "Internal",
					"description": "Existing service description",
				})
				return m
			},
			Err: nil,
		},
	}

	// Create a service update request with non-existent ID
	service := serviceRepository.Service{
		Id:          "999",
		Name:        "UpdatedService",
		ServiceType: "External",
		Description: "Updated service description",
	}
	serviceJson, err := json.Marshal(&service)
	req, err := http.NewRequest("PUT", "/services/999", io.NopCloser(strings.NewReader(string(serviceJson))))

	rw := httptest.NewRecorder()
	handler.UpdateService(rw, req)
	if err != nil {
		t.Errorf("Service UPDATE errored with %s", err.Error())
	}
	if rw.Code != http.StatusNotFound {
		t.Errorf("Service UPDATE returned wrong status code: got %v want %v", rw.Code, http.StatusNotFound)
	}
}

func TestUpdateServiceError(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "1", true
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id": "1",
				})
				return m
			},
			Err: errors.New("test error"),
		},
	}

	service := serviceRepository.Service{
		Id:          "1",
		Name:        "UpdatedService",
		ServiceType: "External",
		Description: "Updated service description",
	}
	serviceJson, err := json.Marshal(&service)
	req, err := http.NewRequest("PUT", "/services/1", io.NopCloser(strings.NewReader(string(serviceJson))))
	if err != nil {
		panic(err)
	}
	rw := httptest.NewRecorder()
	handler.UpdateService(rw, req)
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Service UPDATE returned wrong status code: got %v want %v", rw.Code, http.StatusInternalServerError)
	}
}

func TestUpdateServiceInvalidBody(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "1", false
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id": "1",
				})
				return m
			},
			Err: nil,
		},
	}
	req, err := http.NewRequest("PUT", "/services/1", io.NopCloser(strings.NewReader("some invalid json")))

	rw := httptest.NewRecorder()
	handler.UpdateService(rw, req)
	if err != nil {
		t.Errorf("Service UPDATE errored with %s", err.Error())
	}

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Service UPDATE returned wrong status code: got %v want %v", rw.Code, http.StatusBadRequest)
	}
}

func TestUpdateServiceInvalidId(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "", false // Return error for invalid ID
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id": "1",
				})
				return m
			},
			Err: nil,
		},
	}

	service := serviceRepository.Service{
		Id:          "invalid",
		Name:        "UpdatedService",
		ServiceType: "External",
		Description: "Updated service description",
	}
	serviceJson, err := json.Marshal(&service)
	req, err := http.NewRequest("PUT", "/services/invalid", io.NopCloser(strings.NewReader(string(serviceJson))))

	rw := httptest.NewRecorder()
	handler.UpdateService(rw, req)
	if err != nil {
		t.Errorf("Service UPDATE errored with %s", err.Error())
	}

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Service UPDATE returned wrong status code: got %v want %v", rw.Code, http.StatusBadRequest)
	}
}

func TestUpdateServiceIdMismatch(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "1", true // Return valid ID from path
		},
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				m = append(m, map[string]any{
					"id":          "1",
					"name":        "ExistingService",
					"type":        "Internal",
					"description": "Existing service description",
				})
				return m
			},
			Err: nil,
		},
	}

	// Create a service update request with ID that doesn't match the path ID
	service := serviceRepository.Service{
		Id:          "2", // Different from the ID in the path (1)
		Name:        "UpdatedService",
		ServiceType: "External",
		Description: "Updated service description",
	}
	serviceJson, err := json.Marshal(&service)
	req, err := http.NewRequest("PUT", "/services/1", io.NopCloser(strings.NewReader(string(serviceJson))))

	rw := httptest.NewRecorder()
	handler.UpdateService(rw, req)
	if err != nil {
		t.Errorf("Service UPDATE errored with %s", err.Error())
	}

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Service UPDATE returned wrong status code: got %v want %v", rw.Code, http.StatusBadRequest)
	}
}
