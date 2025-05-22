package services

import (
	"errors"
	"log"
	"net/http"
	"service-dependency-api/internal"
	errors2 "service-dependency-api/internal/customErrors"
)

func (u *ServiceCallsHandler) deleteServiceById(rw http.ResponseWriter, req *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", req)
	log.Println("Request received - deleteServiceById - " + id)
	if !ok {
		http.Error(rw, "Invalid Request", http.StatusBadRequest)
		log.Println("Invalid Request - " + id)
		return
	}

	err := u.Repository.DeleteService(req.Context(), id)
	var httpErr *errors2.HTTPError
	if err != nil {
		if errors.As(err, &httpErr) {
			http.Error(rw, httpErr.Error(), httpErr.Status)
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			log.Println("Error deleting service: " + err.Error())
		}
		return
	}
	rw.WriteHeader(http.StatusNoContent)

}
