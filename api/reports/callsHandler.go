package reports

import (
	"service-dependency-api/neo4jRepositories/reportRepository"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CallsHandler struct {
	repository repositories.ReportRepository
}

func New(driver *neo4j.DriverWithContext) *CallsHandler {
	return &CallsHandler{
		repository: reportRepository.New(*driver),
	}
}
