package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestPOSTSuccess(t *testing.T) {

	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "", false
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

	service := Service{
		Name:        "MockService",
		ServiceType: "Internal",
		Description: "Unit test service",
	}
	serviceJson, err := json.Marshal(&service)
	req, err := http.NewRequest("POST", "/hello/world?name=test", io.NopCloser(strings.NewReader(string(serviceJson))))

	rw := httptest.NewRecorder()
	handler.CreateService(rw, req)
	if err != nil {
		t.Errorf("Service POST errored with %s", err.Error())
	}
	if rw.Code != http.StatusCreated {
		t.Errorf("Service POST errored with %s", strconv.Itoa(rw.Code))
	}
	returnedService := &Service{}
	err = json.Unmarshal(rw.Body.Bytes(), &returnedService)
	switch {
	case err != nil:
		t.Errorf("Service POST errored parsing return body with %s", err.Error())
	case returnedService.Id != "1":
		t.Errorf("Service POST errored with %s", returnedService.Id)
	case returnedService.Name != "MockService":
		t.Errorf("Service POST errored with Name inconsistency %s", returnedService.Name)
	case returnedService.Description != "Unit test service":
		t.Errorf("Service POST errored with Description inconsistency %s", returnedService.Description)
	case returnedService.ServiceType != "Internal":
		t.Errorf("Service POST errored with ServiceType inconsistency %s", returnedService.ServiceType)

	}

}

func TestPOSTError(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "", false
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

	service := Service{
		Name:        "MockService",
		ServiceType: "Internal",
		Description: "Unit test service",
	}
	serviceJson, err := json.Marshal(&service)
	req, err := http.NewRequest("POST", "/hello/world?name=test", io.NopCloser(strings.NewReader(string(serviceJson))))
	if err != nil {
		panic(err)
	}
	rw := httptest.NewRecorder()
	handler.CreateService(rw, req)
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Service POST errored with %s", strconv.Itoa(rw.Code))
	}
}

func TestPOSTInvalidBody(t *testing.T) {
	handler := ServiceCallsHandler{
		IdValidator: func(_ string, _ *http.Request) (string, bool) {
			return "", false
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
	req, err := http.NewRequest("POST", "/hello/world?name=test", io.NopCloser(strings.NewReader("some text")))

	rw := httptest.NewRecorder()
	handler.CreateService(rw, req)
	if err != nil {
		t.Errorf("Service POST errored with %s", err.Error())
	}

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Service POST errored with %s", strconv.Itoa(rw.Code))
	}
}
