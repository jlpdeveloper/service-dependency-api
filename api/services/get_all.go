package services

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type GetAllServicesHandler struct {
	Path       string
	Repository ServiceRepository
}

func (u *GetAllServicesHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(u.Path, u.ServeHTTP)
}

func (u *GetAllServicesHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	pageSize, err := strconv.Atoi(req.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 10
	}
	services, err := u.Repository.GetAllServices(page, pageSize)
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
