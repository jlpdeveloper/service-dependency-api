package services

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
)

func Register(mux *http.ServeMux, driver *neo4j.DriverWithContext, ctx *context.Context) {

	serviceRepo := ServiceNeo4jRepository{
		Ctx:    *ctx,
		Driver: *driver,
	}
	//getAllHandler := GetAllServicesHandler{
	//	Path:       "GET /services",
	//	Repository: &serviceRepo,
	//}
	//
	//getAllHandler.Register(mux)

	postHandler := POSTServicesHandler{
		Path:       "POST /services",
		Repository: &serviceRepo,
	}
	postHandler.Register(mux)
}
