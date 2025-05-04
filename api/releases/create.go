package releases

import (
	"encoding/json"
	"net/http"
	"service-dependency-api/api/releases/internal/releaseRepository"
	"service-dependency-api/internal/customErrors"
)

func (s *ServiceCallsHandler) createRelease(rw http.ResponseWriter, req *http.Request) {
	serviceId, ok := s.PathValidator("id", req)
	if !ok {
		http.Error(rw, "Invalid service ID", http.StatusBadRequest)
		return
	}

	r := &releaseRepository.Release{}
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
