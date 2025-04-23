package services

import (
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/api/services/service_repository"
)

type ServiceCallsHandler struct {
	Repository  service_repository.ServiceRepository
	IdValidator func(string, *http.Request) (string, bool)
}

func (u *ServiceCallsHandler) Register(mux *http.ServeMux) {
	paths := map[string]func(http.ResponseWriter, *http.Request){
		"GET /services/{id}":    u.GetById,
		"GET /services":         u.GetAllServices,
		"POST /services":        u.CreateService,
		"PUT /services/{id}":    u.UpdateService,
		"DELETE /services/{id}": u.DeleteServiceById,
	}
	for path, f := range paths {
		mux.HandleFunc(path, f)
	}
}

func getGuidFromRequestPath(varName string, req *http.Request) (string, bool) {
	guidVal := req.PathValue(varName)
	err := uuid.Validate(guidVal)
	return guidVal, err == nil
}

func Register(mux *http.ServeMux, driver *neo4j.DriverWithContext) {

	serviceRepo := &service_repository.ServiceNeo4jRepository{
		Driver: *driver,
	}

	callsHandler := ServiceCallsHandler{
		Repository:  serviceRepo,
		IdValidator: getGuidFromRequestPath,
	}

	callsHandler.Register(mux)
}
