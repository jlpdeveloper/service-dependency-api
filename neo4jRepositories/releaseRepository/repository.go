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

// releaseTimeFormat the default time format for formatting the release date
const releaseTimeFormat = "2006-01-02T15:04:05Z"
