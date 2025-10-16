package releases

import (
	"service-dependency-api/neo4jRepositories/releaseRepository"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ServiceCallsHandler struct {
	Repository repositories.ReleaseRepository
}

func New(driver *neo4j.DriverWithContext) *ServiceCallsHandler {
	return &ServiceCallsHandler{
		Repository: releaseRepository.New(*driver),
	}
}
