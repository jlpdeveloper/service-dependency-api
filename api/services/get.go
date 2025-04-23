package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (u *ServiceCallsHandler) GetAllServices(rw http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate that page is positive
	if page < 1 {
		http.Error(rw, "page must be positive", http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(req.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	// Validate that pageSize is between 1 and 100
	if pageSize < 1 || pageSize > 100 {
		http.Error(rw, "pageSize must be between 1 and 100", http.StatusBadRequest)
		return
	}
	services, err := u.Repository.GetAllServices(req.Context(), page, pageSize)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(services)
	if err != nil {
		log.Println(err)
	}
}

func (u *ServiceCallsHandler) GetById(rw http.ResponseWriter, req *http.Request) {
	id, ok := u.IdValidator("id", req)

	if !ok {
		http.Error(rw, "Service id is required", http.StatusBadRequest)
		return
	}

	service, err := u.Repository.GetServiceById(req.Context(), id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if service was found (Id will be empty if not found)
	if service.Id == "" {
		http.Error(rw, "Service not found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(service)
	if err != nil {
		log.Println(err)
	}
}
