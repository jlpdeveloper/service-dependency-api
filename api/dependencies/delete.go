package dependencies

import (
	"net/http"
	"service-atlas/internal"
	"service-atlas/internal/customerrors"
)

func (s *ServiceCallsHandler) DeleteDependency(rw http.ResponseWriter, req *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", req)
	if !ok {
		http.Error(rw, "path id not valid", http.StatusBadRequest)
		return
	}
	dependsOnID, ok := internal.GetGuidFromRequestPath("id2", req)
	if !ok {
		http.Error(rw, "path id2 not valid", http.StatusBadRequest)
		return
	}
	err := s.Repository.DeleteDependency(req.Context(), id, dependsOnID)

	if err != nil {
		customerrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
