package teams

import (
	"service-atlas/neo4jrepositories/teamrepository"
	"service-atlas/repositories"

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
