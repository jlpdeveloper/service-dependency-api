package dependencies

import (
	"service-atlas/neo4jrepositories/dependencyrepository"
	"service-atlas/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ServiceCallsHandler struct {
	Repository repositories.DependencyRepository
}

func New(driver neo4j.DriverWithContext) *ServiceCallsHandler {
	return &ServiceCallsHandler{
		Repository: dependencyrepository.New(driver),
	}
}
