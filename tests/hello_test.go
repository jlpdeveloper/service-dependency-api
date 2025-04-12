package tests

import (
	"net/http"
	"net/http/httptest"
	"service-dependency-api/api/hello_world"
	"testing"
)

func TestHelloWorldEmptyName(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello/world", nil)
	rw := httptest.NewRecorder()
	hello_world.HelloWorld(rw, req)
	if err != nil {
		t.Errorf("HelloWorld errored with %s", err.Error())
	}
	if rw.Code != http.StatusOK {
		t.Errorf("HelloWorld errored with %s", string(rune(rw.Code)))
	}
	if rw.Body.String() != "hello world" {
		t.Errorf("HelloWorld errored with %s", rw.Body.String())
	}
}

func TestHelloWorldWithName(t *testing.T) {
	req, err := http.NewRequest("GET", "/hello/world?name=test", nil)

	rw := httptest.NewRecorder()
	hello_world.HelloWorld(rw, req)
	if err != nil {
		t.Errorf("HelloWorld errored with %s", err.Error())
	}
	if rw.Code != http.StatusOK {
		t.Errorf("HelloWorld errored with %s", string(rune(rw.Code)))
	}
	if rw.Body.String() != "hello test" {
		t.Errorf("HelloWorld errored with %s", rw.Body.String())
	}

}
