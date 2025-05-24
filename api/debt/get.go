package debt

import (
	"encoding/json"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
	"strconv"
)

func (c CallsHandler) getDebtByServiceId(rw http.ResponseWriter, r *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", r)
	if !ok {
		http.Error(rw, "service id not valid", http.StatusBadRequest)
		return
	}
	page, pageSize, err := validatePageParams(r)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	onlyResolved := r.URL.Query().Get("onlyResolved")

	debt, err := c.Repository.GetDebtByServiceId(r.Context(), id, page, pageSize, onlyResolved == "true")
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(debt)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func validatePageParams(req *http.Request) (int, int, error) {

	// Parse page and page size from query parameters
	pageStr := req.URL.Query().Get("page")
	pageSizeStr := req.URL.Query().Get("pageSize")

	// Default values
	page := 1
	pageSize := 25

	// Parse page if provided
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil || parsedPage <= 0 {
			return 0, 0, &customErrors.HTTPError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid page parameter",
			}
		}
		page = parsedPage
	}

	// Parse page size if provided
	if pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil || parsedPageSize <= 0 {
			return 0, 0, &customErrors.HTTPError{
				Status: http.StatusBadRequest,
				Msg:    "Invalid page_size parameter",
			}
		}
		pageSize = parsedPageSize
	}

	return page, pageSize, nil
}
