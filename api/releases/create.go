package releases

import (
	"encoding/json"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
)

func (s *ServiceCallsHandler) createRelease(rw http.ResponseWriter, req *http.Request) {
	serviceId, ok := internal.GetGuidFromRequestPath("id", req)
	if !ok {
		http.Error(rw, "Invalid service ID", http.StatusBadRequest)
		return
	}

	r := &repositories.Release{}
	const maxBodySize = 1 << 20 // 1 MB
	req.Body = http.MaxBytesReader(rw, req.Body, maxBodySize)
	err := json.NewDecoder(req.Body).Decode(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the service ID from the path parameter
	r.ServiceId = serviceId

	if err = r.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.Repository.CreateRelease(req.Context(), *r)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
