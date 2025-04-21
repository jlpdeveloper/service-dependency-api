package services

import (
	"errors"
	"log"
	"net/http"
	"service-dependency-api/internal"
)

func (u *ServiceCallsHandler) DeleteServiceById(rw http.ResponseWriter, req *http.Request) {
	id, ok := u.IdValidator("id", req)
	log.Println("Request received - DeleteServiceById - " + id)
	if !ok {
		http.Error(rw, "Invalid Request", http.StatusBadRequest)
		log.Println("Invalid Request - " + id)
		return
	}

	err := u.Repository.DeleteService(id)
	var httpErr *internal.HTTPError
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
