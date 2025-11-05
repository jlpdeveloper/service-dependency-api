package teams

import (
	"service-dependency-api/neo4jrepositories/teamrepository"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CallsHandler struct {
	Repository repositories.TeamRepository
}

func New(driver neo4j.DriverWithContext) *CallsHandler {
	return &CallsHandler{
		Repository: teamrepository.New(driver),
	}
}
