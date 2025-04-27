package dependencies

import (
	"encoding/json"
	"log"
	"net/http"
	"service-dependency-api/internal/customErrors"
)

func (s *ServiceCallsHandler) getById(rw http.ResponseWriter, req *http.Request) {
	id, ok := s.PathValidator("id", req)
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
