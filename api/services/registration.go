package services

import (
	"net/http"
	"service-dependency-api/neo4jRepositories/serviceRepository"
	"service-dependency-api/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ServiceCallsHandler struct {
	Repository repositories.ServiceRepository
}

func (u *ServiceCallsHandler) Register(mux *chi.Mux) {
	paths := map[string]func(http.ResponseWriter, *http.Request){
		"GET /services/{id}":    u.getById,
		"GET /services":         u.getAllServices,
		"POST /services":        u.createService,
		"PUT /services/{id}":    u.updateService,
		"DELETE /services/{id}": u.deleteServiceById,
	}
	for path, f := range paths {
		mux.HandleFunc(path, f)
	}
}

func Register(mux *chi.Mux, driver *neo4j.DriverWithContext) {
	callsHandler := ServiceCallsHandler{
		Repository: serviceRepository.New(*driver),
	}
	callsHandler.Register(mux)
}
