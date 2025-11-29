package reports

import (
	"service-atlas/neo4jrepositories/reportrepository"
	"service-atlas/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CallsHandler struct {
	repository repositories.ReportRepository
}

func New(driver neo4j.DriverWithContext) *CallsHandler {
	return &CallsHandler{
		repository: reportrepository.New(driver),
	}
}
