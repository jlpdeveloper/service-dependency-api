package dependencies

import (
	"encoding/json"
	"log"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
)

func (s *ServiceCallsHandler) getDependencies(rw http.ResponseWriter, req *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", req)
	if !ok {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	dep, err := s.Repository.GetDependencies(req.Context(), id)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(dep)
	if err != nil {
		log.Println(err)
	}
}

func (s *ServiceCallsHandler) getDependents(rw http.ResponseWriter, req *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", req)
	if !ok {
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	deps, err := s.Repository.GetDependents(req.Context(), id)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	ver := req.URL.Query().Get("version")
	returnObj := make([]*repositories.Dependency, 0)
	if ver != "" {
		for _, dep := range deps {
			if ver == dep.Version {
				returnObj = append(returnObj, dep)
			}
		}
	} else {
		returnObj = append(returnObj, deps...)
	}

	err = json.NewEncoder(rw).Encode(returnObj)
	if err != nil {
		log.Println(err)
	}
}
