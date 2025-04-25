package dependencies

import (
	"encoding/json"
	"errors"
	"net/http"
	"service-dependency-api/api/dependencies/internal/dependencyRepository"
	"service-dependency-api/internal/customErrors"
)

func (s *ServiceCallsHandler) createDependency(rw http.ResponseWriter, req *http.Request) {
	id, ok := s.PathValidator("id", req)
	if !ok {
		http.Error(rw, "path id not valid", http.StatusBadRequest)
		return
	}
	dep := &dependencyRepository.Dependency{}
	err := json.NewDecoder(req.Body).Decode(dep)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := dep.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.Repository.AddDependency(id, dep)
	var httpErr *customErrors.HTTPError
	if err != nil {
		if errors.As(err, &httpErr) {
			http.Error(rw, httpErr.Error(), httpErr.Status)
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
