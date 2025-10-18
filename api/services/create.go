package services

import (
	"encoding/json"
	"log"
	"net/http"
	"service-dependency-api/repositories"
)

func (u *ServiceCallsHandler) CreateService(w http.ResponseWriter, req *http.Request) {

	createServiceRequest := &repositories.Service{}
	err := json.NewDecoder(req.Body).Decode(createServiceRequest)
	if err != nil {
		// return HTTP 400 bad request
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = createServiceRequest.Validate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
