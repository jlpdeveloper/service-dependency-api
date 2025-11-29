package releases

import (
	"service-atlas/neo4jrepositories/releaserepository"
	"service-atlas/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ServiceCallsHandler struct {
	Repository repositories.ReleaseRepository
}

func New(driver neo4j.DriverWithContext) *ServiceCallsHandler {
	return &ServiceCallsHandler{
		Repository: releaserepository.New(driver),
	}
}
