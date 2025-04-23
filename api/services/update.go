package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"service-dependency-api/api/services/internal/serviceRepository"
	"service-dependency-api/internal"
)

func (u *ServiceCallsHandler) UpdateService(rw http.ResponseWriter, req *http.Request) {
	updateServiceRequest := &serviceRepository.Service{}
	err := json.NewDecoder(req.Body).Decode(updateServiceRequest)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if id, ok := u.IdValidator("id", req); !ok || updateServiceRequest.Id != id {
		http.Error(rw, "Service Id is not valid", http.StatusBadRequest)
		return
	}

	err = u.Repository.UpdateService(req.Context(), *updateServiceRequest)

	var httpErr *internal.HTTPError
	if err != nil {
		if errors.As(err, &httpErr) {
			http.Error(rw, httpErr.Error(), httpErr.Status)
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	rw.WriteHeader(http.StatusNoContent)

}
