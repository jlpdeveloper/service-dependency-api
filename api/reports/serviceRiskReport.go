package reports

import (
	"context"
	"encoding/json"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
	"time"
)

func (c *CallsHandler) GetServiceRiskReport(rw http.ResponseWriter, req *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", req)
	if !ok {
		http.Error(rw, "Invalid service ID", http.StatusBadRequest)
		return
	}
	contextWithTimeout, cancel := context.WithTimeout(req.Context(), 10*time.Second)
	defer cancel()
	report, err := c.repository.GetServiceRiskReport(contextWithTimeout, id)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	err = json.NewEncoder(rw).Encode(report)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Header().Set("Content-Type", "application/json")
}
