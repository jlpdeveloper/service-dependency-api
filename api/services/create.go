package services

import (
	"encoding/json"
	"log"
	"net/http"
	"service-dependency-api/api/services/serviceRepository"
)

func (u *ServiceCallsHandler) CreateService(w http.ResponseWriter, req *http.Request) {

	createServiceRequest := &serviceRepository.Service{}
	err := json.NewDecoder(req.Body).Decode(createServiceRequest)
	if err != nil {
		// return HTTP 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id, err := u.Repository.CreateService(req.Context(), *createServiceRequest)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	createServiceRequest.Id = id
	err = json.NewEncoder(w).Encode(createServiceRequest)
	if err != nil {
		log.Println(err)
	}
}
