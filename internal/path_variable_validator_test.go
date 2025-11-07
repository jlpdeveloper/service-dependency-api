package internal

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetGuidFromRequestPath(t *testing.T) {
	testCases := []struct {
		name     string
		reqPath  string
		expected string
		ok       bool
	}{
		{"Valid GUID", "/service/85622399-b2b7-4e98-9a8d-628e28b9aeb4", "85622399-b2b7-4e98-9a8d-628e28b9aeb4", true},
		{"Invalid GUID", "/service/invalid-guid", "invalid-guid", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.reqPath, nil)
			id := strings.Split(tc.reqPath, "/")
			req.SetPathValue("id", id[2])
			guidVal, ok := GetGuidFromRequestPath("id", req)
			if guidVal != tc.expected || ok != tc.ok {
				t.Errorf("GetGuidFromRequestPath(%q) = (%q, %v), want (%q, %v)", tc.reqPath, guidVal, ok, tc.expected, tc.ok)
			}
		})
	}
}

func TestGetDateFromRequestPath(t *testing.T) {
	testCases := []struct {
		name     string
		reqPath  string
		expected time.Time
		ok       bool
	}{
		{"Valid Date", "/service/2023-10-05", time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC), true},
		{"Invalid Date", "/service/2023-10-32", time.Time{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.reqPath, nil)
			pathVars := strings.Split(tc.reqPath, "/")
			req.SetPathValue("startDate", pathVars[2])
			dateVal, ok := GetDateFromRequestPath("startDate", req)
			if dateVal != tc.expected || ok != tc.ok {
				t.Errorf("GetDateFromRequestPath(%q) = (%v, %v), want (%v, %v)", tc.reqPath, dateVal, ok, tc.expected, tc.ok)
			}
		})
	}
}
