package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

type ServiceCallsHandler struct {
	Repository  ServiceRepository
	IdValidator func(string, *http.Request) (string, bool)
}

func (u *ServiceCallsHandler) Register(mux *http.ServeMux) {
	paths := map[string]func(http.ResponseWriter, *http.Request){
		"GET /services/{id}": u.GetById,
		"GET /services":      u.GetAllServices,
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

func Register(mux *http.ServeMux, driver *neo4j.DriverWithContext, ctx *context.Context) {

	serviceRepo := &ServiceNeo4jRepository{
		Ctx:    *ctx,
		Driver: *driver,
	}

	callsHandler := ServiceCallsHandler{
		Repository:  serviceRepo,
		IdValidator: getGuidFromRequestPath,
	}

	postHandler := POSTServicesHandler{
		Path:       "POST /services",
		Repository: serviceRepo,
	}
	postHandler.Register(mux)

	callsHandler.Register(mux)
}
