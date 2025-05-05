package releases

import (
	"encoding/json"
	"net/http"
	"service-dependency-api/internal/customErrors"
	"strconv"
)

func (s *ServiceCallsHandler) getReleasesByServiceId(rw http.ResponseWriter, req *http.Request) {
	serviceId, ok := s.PathValidator("id", req)
	if !ok {
		http.Error(rw, "Invalid service ID", http.StatusBadRequest)
		return
	}

	// Parse page and page size from query parameters
	pageStr := req.URL.Query().Get("page")
	pageSizeStr := req.URL.Query().Get("page_size")

	// Default values
	page := 1
	pageSize := 25

	// Parse page if provided
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil || parsedPage <= 0 {
			http.Error(rw, "Invalid page parameter", http.StatusBadRequest)
			return
		}
		page = parsedPage
	}

	// Parse page size if provided
	if pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || parsedPageSize <= 0 {
			http.Error(rw, "Invalid page_size parameter", http.StatusBadRequest)
			return
		}
		pageSize = parsedPageSize
	}

	releases, err := s.Repository.GetReleasesByServiceId(req.Context(), serviceId, page, pageSize)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}

	// Set content type header
	rw.Header().Set("Content-Type", "application/json")

	// Encode the releases as JSON and write to the response
	err = json.NewEncoder(rw).Encode(releases)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
