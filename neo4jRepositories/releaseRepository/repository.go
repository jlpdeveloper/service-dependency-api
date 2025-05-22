package releaseRepository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/databaseAdapter"
)

type Neo4jReleaseRepository struct {
	manager databaseAdapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jReleaseRepository {
	return &Neo4jReleaseRepository{
		manager: databaseAdapter.NewDriverManager(driver),
	}
}
