package services

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/neo4jRepositories/serviceRepository"
	"service-dependency-api/repositories"
)

type ServiceCallsHandler struct {
	Repository repositories.ServiceRepository
}

func (u *ServiceCallsHandler) Register(mux *http.ServeMux) {
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

func Register(mux *http.ServeMux, driver *neo4j.DriverWithContext) {
	callsHandler := ServiceCallsHandler{
		Repository: serviceRepository.New(*driver),
	}
	callsHandler.Register(mux)
}
