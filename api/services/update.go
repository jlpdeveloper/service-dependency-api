package services

import (
	"encoding/json"
	"net/http"
	"service-dependency-api/api/services/service_repository"
)

func (u *ServiceCallsHandler) UpdateService(rw http.ResponseWriter, req *http.Request) {
	updateServiceRequest := &service_repository.Service{}
	err := json.NewDecoder(req.Body).Decode(updateServiceRequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if id, ok := u.IdValidator("id", req); !ok || updateServiceRequest.Id != id {
		http.Error(rw, "Service Id is not valid", http.StatusBadRequest)
		return
	}

	found, err := u.Repository.UpdateService(req.Context(), *updateServiceRequest)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if !found {
		http.Error(rw, "Service not found", http.StatusNotFound)
		return
	}
	rw.WriteHeader(http.StatusNoContent)

}
