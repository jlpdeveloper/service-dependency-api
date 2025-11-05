package debt

import (
	"service-dependency-api/neo4jrepositories/debtrepository"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CallsHandler struct {
	Repository repositories.DebtRepository
}

func New(driver neo4j.DriverWithContext) *CallsHandler {
	return &CallsHandler{
		Repository: debtrepository.New(driver),
	}
}
