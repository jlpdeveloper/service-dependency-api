package helloworld

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHelloWorldEmptyName(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello/world", nil)
	rw := httptest.NewRecorder()
	HelloWorld(rw, req)
	if err != nil {
		t.Errorf("HelloWorld errored with %s", err.Error())
	}
	if rw.Code != http.StatusOK {
		t.Errorf("HelloWorld errored with %s", strconv.Itoa(rw.Code))
	}
	if rw.Body.String() != "hello world" {
		t.Errorf("HelloWorld errored with %s", rw.Body.String())
	}
}

func TestHelloWorldWithName(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello/world?name=test", nil)

	rw := httptest.NewRecorder()
	HelloWorld(rw, req)
	if err != nil {
		t.Errorf("HelloWorld errored with %s", err.Error())
	}
	if rw.Code != http.StatusOK {
		t.Errorf("HelloWorld errored with %s", strconv.Itoa(rw.Code))
	}
	if rw.Body.String() != "hello test" {
		t.Errorf("HelloWorld errored with %s", rw.Body.String())
	}
}
