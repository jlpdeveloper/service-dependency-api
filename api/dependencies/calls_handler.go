package dependencies

import (
	"service-dependency-api/neo4jrepositories/dependencyRepository"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ServiceCallsHandler struct {
	Repository repositories.DependencyRepository
}

func New(driver neo4j.DriverWithContext) *ServiceCallsHandler {
	return &ServiceCallsHandler{
		Repository: dependencyRepository.New(driver),
	}
}
