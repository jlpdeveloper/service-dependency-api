package releases

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
	"testing"
	"time"
)

func TestGetReleasesByServiceIdSuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	mockReleases := []*repositories.Release{
		{
			ServiceId:   validServiceId,
			ReleaseDate: time.Now().UTC(),
			Url:         "https://example.com/release1",
			Version:     "1.0.0",
		},
		{
			ServiceId:   validServiceId,
			ReleaseDate: time.Now().UTC(),
			Url:         "https://example.com/release2",
			Version:     "2.0.0",
		},
	}

	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err:      nil, // No error
			Releases: mockReleases,
		},
	}

	// Create a request with no pagination parameters (should use defaults)
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/release", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Decode the response body
	var releases []*repositories.Release
	err = json.NewDecoder(rw.Body).Decode(&releases)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check that the correct number of releases was returned
	if len(releases) != len(mockReleases) {
		t.Errorf("Expected %d releases, got %d", len(mockReleases), len(releases))
	}
}

func TestGetReleasesByServiceIdInvalidPathParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
	}

	// Create a request
	req, err := http.NewRequest("GET", "/services/invalid-id/release", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesByServiceIdRepositoryError(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: errors.New("repository error"), // Simulate a repository error
		},
	}

	// Create a request
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/release", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, rw.Code)
	}
}

func TestGetReleasesByServiceIdHTTPError(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: &customErrors.HTTPError{
				Status: http.StatusNotFound,
				Msg:    "Service not found",
			}, // Simulate an HTTP error
		},
	}

	// Create a request
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/release", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rw.Code)
	}
}

func TestGetReleasesByServiceIdWithPagination(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID

	// Create 30 mock releases
	mockReleases := make([]*repositories.Release, 30)
	for i := 0; i < 30; i++ {
		mockReleases[i] = &repositories.Release{
			ServiceId:   validServiceId,
			ReleaseDate: time.Now().UTC(),
			Url:         fmt.Sprintf("https://example.com/release%d", i+1),
			Version:     fmt.Sprintf("%d.0.0", i+1),
		}
	}

	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err:      nil, // No error
			Releases: mockReleases,
		},
	}

	// Create a request with pagination parameters (page=1, pageSize=10)
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/release?page=2&pageSize=10", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Decode the response body
	var releases []*repositories.Release
	err = json.NewDecoder(rw.Body).Decode(&releases)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check that the correct number of releases was returned (should be 10)
	if len(releases) != 10 {
		t.Errorf("Expected %d releases, got %d", 10, len(releases))
	}

	// Check that the correct page of releases was returned (should be releases 10-19)
	for i, release := range releases {
		expectedVersion := fmt.Sprintf("%d.0.0", i+11) // Page 1 starts at index 10
		if release.Version != expectedVersion {
			t.Errorf("Expected version %s, got %s", expectedVersion, release.Version)
		}
	}
}

func TestGetReleasesByServiceIdInvalidPageParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid page parameter
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/release?page=invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesByServiceIdInvalidPageSizeParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid pageSize parameter
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/release?pageSize=invalid", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesByServiceIdNegativePageParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
	}

	// Create a request with negative page parameter
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/release?page=-1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesByServiceIdZeroPageSizeParameter(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
	}

	// Create a request with zero pageSize parameter
	req, err := http.NewRequest("GET", "/services/"+validServiceId+"/release?pageSize=0", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesByServiceId(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesInDateRangeSuccess(t *testing.T) {
	// Create a handler with mocked dependencies
	validServiceId := "123e4567-e89b-12d3-a456-426614174000" // Valid GUID
	mockServiceInfo := []*repositories.ServiceReleaseInfo{
		{
			ServiceType: "service-type-1",
			ServiceName: "service-name-1",
			Release: repositories.Release{
				ServiceId:   validServiceId,
				ReleaseDate: time.Now().UTC(),
				Url:         "https://example.com/release1",
				Version:     "1.0.0",
			},
		},
		{
			ServiceType: "service-type-2",
			ServiceName: "service-name-2",
			Release: repositories.Release{
				ServiceId:   validServiceId,
				ReleaseDate: time.Now().UTC(),
				Url:         "https://example.com/release2",
				Version:     "2.0.0",
			},
		},
	}

	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err:         nil, // No error
			ServiceInfo: mockServiceInfo,
		},
	}
	// Create a request with no pagination parameters (should use defaults)
	req, err := http.NewRequest("GET", "/releases/2025-01-01/2025-02-02", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set path values for the request
	req = req.WithContext(req.Context())
	req.SetPathValue("startDate", "2025-01-01")
	req.SetPathValue("endDate", "2025-02-02")
	req.SetPathValue("id", validServiceId)
	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesInDateRange(rw, req)

	// Check the response
	if rw.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rw.Code)
	}

	// Decode the response body
	var releases []*repositories.ServiceReleaseInfo
	err = json.NewDecoder(rw.Body).Decode(&releases)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	// Check that the correct number of releases was returned
	if len(releases) != len(mockServiceInfo) {
		t.Errorf("Expected %d releases, got %d", len(mockServiceInfo), len(releases))
	}
}

func TestGetReleasesInDateRangeInvalidStartDate(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid start date
	req, err := http.NewRequest("GET", "/releases/invalid-date/2025-02-02", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set path values for the request
	req = req.WithContext(req.Context())
	req.SetPathValue("start_date", "invalid-date")
	req.SetPathValue("end_date", "2025-02-02")

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesInDateRange(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesInDateRangeInvalidEndDate(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
	}

	// Create a request with invalid end date
	req, err := http.NewRequest("GET", "/releases/2025-01-01/invalid-date", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set path values for the request
	req = req.WithContext(req.Context())
	req.SetPathValue("start_date", "2025-01-01")
	req.SetPathValue("end_date", "invalid-date")

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesInDateRange(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}

func TestGetReleasesInDateRangeEndDateBeforeStartDate(t *testing.T) {
	// Create a handler with mocked dependencies
	handler := ServiceCallsHandler{
		Repository: mockReleaseRepository{
			Err: nil, // No error
		},
	}

	// Create a request with end date before start date
	req, err := http.NewRequest("GET", "/releases/2025-02-02/2025-01-01", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set path values for the request
	req = req.WithContext(req.Context())
	req.SetPathValue("start_date", "2025-02-02")
	req.SetPathValue("end_date", "2025-01-01")

	// Create a response recorder
	rw := httptest.NewRecorder()

	// Call the handler
	handler.GetReleasesInDateRange(rw, req)

	// Check the response
	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rw.Code)
	}
}
