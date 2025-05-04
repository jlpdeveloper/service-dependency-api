package dependencies

import (
	"net/http"
	"service-dependency-api/internal/customErrors"
)

func (s *ServiceCallsHandler) deleteDependency(rw http.ResponseWriter, req *http.Request) {
	id, ok := s.PathValidator("id", req)
	if !ok {
		http.Error(rw, "path id not valid", http.StatusBadRequest)
		return
	}
	dependsOnID, ok := s.PathValidator("id2", req)
	if !ok {
		http.Error(rw, "path id2 not valid", http.StatusBadRequest)
		return
	}
	err := s.Repository.DeleteDependency(req.Context(), id, dependsOnID)

	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
