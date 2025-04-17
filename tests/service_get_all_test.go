package tests

import (
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/services"
	"strconv"
	"testing"
)

func TestGetAllSuccess(t *testing.T) {
	handler := services.GetAllServicesHandler{

		Path: "/services",
		Repository: MockServiceRepository{
			Data: func() []map[string]any {
				var m []map[string]any
				for i := 0; i < 10; i++ {
					m = append(m, map[string]any{
						"id":          strconv.Itoa(i),
						"name":        "service" + strconv.Itoa(i),
						"description": "test desc",
						"type":        "service",
					})
				}
				return m
			},
			Err: nil,
		},
	}
	req, err := http.NewRequest("GET", "/hello/world?page=1&pageSize=5", nil)

	rw := httptest.NewRecorder()
	handler.ServeHTTP(rw, req)
	if err != nil {
		t.Errorf("Service GET get all errored with %s", err.Error())
	}
}
