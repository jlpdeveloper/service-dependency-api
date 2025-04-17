package tests

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/services"
	"strconv"
	"strings"
	"testing"
)

func TestSuccessCreate(t *testing.T) {

	handler := services.POSTServicesHandler{

		Path: "/services",
		Repository: MockServiceRepository{
			Data: map[string]any{
				"id": "1",
			},
			Err: nil,
		},
	}

	service := services.Service{
		Name:        "MockService",
		ServiceType: "Internal",
		Description: "Unit test service",
	}
	serviceJson, err := json.Marshal(&service)
	req, err := http.NewRequest("POST", "/hello/world?name=test", io.NopCloser(strings.NewReader(string(serviceJson))))

	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)
	if err != nil {
		t.Errorf("Service POST errored with %s", err.Error())
	}
	if rw.Code != http.StatusCreated {
		t.Errorf("Service POST errored with %s", strconv.Itoa(rw.Code))
	}
	returnedService := &services.Service{}
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

func TestErrorCreate(t *testing.T) {
	handler := services.POSTServicesHandler{

		Path: "/services",
		Repository: MockServiceRepository{
			Data: map[string]any{
				"id": "1",
			},
			Err: errors.New("test error"),
		},
	}

	service := services.Service{
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
	handler.ServeHTTP(rw, req)
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Service POST errored with %s", strconv.Itoa(rw.Code))
	}
}

func TestInvalidBody(t *testing.T) {
	handler := services.POSTServicesHandler{

		Path: "/services",
		Repository: MockServiceRepository{
			Data: map[string]any{
				"id": "1",
			},
			Err: nil,
		},
	}
	req, err := http.NewRequest("POST", "/hello/world?name=test", io.NopCloser(strings.NewReader("some text")))

	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)
	if err != nil {
		t.Errorf("Service POST errored with %s", err.Error())
	}

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Service POST errored with %s", strconv.Itoa(rw.Code))
	}
}
