package services

import (
	"service-dependency-api/neo4jrepositories/servicerepository"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ServiceCallsHandler struct {
	Repository repositories.ServiceRepository
}

func New(driver neo4j.DriverWithContext) *ServiceCallsHandler {
	return &ServiceCallsHandler{
		Repository: servicerepository.New(driver),
	}
}
