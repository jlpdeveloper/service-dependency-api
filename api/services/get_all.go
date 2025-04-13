package services

import "net/http"

type GetAllServicesHandler struct {
	Path       string
	Repository ServiceRepository
}

func (u *GetAllServicesHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc(u.Path, u.ServeHTTP)
}

func (u *GetAllServicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
