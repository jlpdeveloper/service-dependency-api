package reportrepository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-atlas/databaseadapter"
)

type Neo4jReportRepository struct {
	manager databaseadapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jReportRepository {
	return &Neo4jReportRepository{manager: databaseadapter.NewDriverManager(driver)}
}
