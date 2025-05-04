package releases

import (
	"encoding/json"
	"net/http"
	"service-dependency-api/api/releases/internal/releaseRepository"
	"service-dependency-api/internal/customErrors"
)

func (s *ServiceCallsHandler) createRelease(rw http.ResponseWriter, req *http.Request) {
	r := &releaseRepository.Release{}
	err := json.NewDecoder(req.Body).Decode(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
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
