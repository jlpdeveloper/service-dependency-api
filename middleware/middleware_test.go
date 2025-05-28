package middleware

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateStackOrder(t *testing.T) {
	// Create a slice to record the order of execution
	var executionOrder []int

	// Create test middleware functions that record when they're executed
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, 1)
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, 2)
			next.ServeHTTP(w, r)
		})
	}

	middleware3 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, 3)
			next.ServeHTTP(w, r)
		})
	}

	// Create the final handler
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, 4)
	})

	// Create the middleware stack
	stack := CreateStack(middleware1, middleware2, middleware3)

	// Apply the stack to the final handler
	handler := stack(finalHandler)

	// Execute the handler with a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the execution order (should be 1, 2, 3, 4 since middleware are applied in reverse)
	expected := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(executionOrder, expected) {
		t.Errorf("Middleware executed in wrong order: got %v want %v", executionOrder, expected)
	}
}

func TestCreateStackModification(t *testing.T) {
	// Create middleware that adds a header to the request
	addRequestHeader := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Header.Add("X-Test-Header", "test-value")
			next.ServeHTTP(w, r)
		})
	}

	// Create middleware that adds a header to the response
	addResponseHeader := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Response-Header", "response-value")
			next.ServeHTTP(w, r)
		})
	}

	// Create the final handler that checks the request header
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Header") != "test-value" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Create the middleware stack
	stack := CreateStack(addRequestHeader, addResponseHeader)

	// Apply the stack to the final handler
	handler := stack(finalHandler)

	// Execute the handler with a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check that the response status is OK (meaning the request header was properly set)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that the response header was properly set
	if rr.Header().Get("X-Response-Header") != "response-value" {
		t.Errorf("Response header not set correctly: got %v want %v",
			rr.Header().Get("X-Response-Header"), "response-value")
	}
}

func TestCreateStackEmpty(t *testing.T) {
	// Test with no middleware
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create an empty middleware stack
	stack := CreateStack()

	// Apply the stack to the final handler
	handler := stack(finalHandler)

	// Execute the handler with a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check that the response status is OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
