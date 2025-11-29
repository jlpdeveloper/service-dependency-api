package releaserepository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-atlas/databaseadapter"
)

type Neo4jReleaseRepository struct {
	manager databaseadapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jReleaseRepository {
	return &Neo4jReleaseRepository{
		manager: databaseadapter.NewDriverManager(driver),
	}
}

// releaseTimeFormat the default time format for formatting the release date
const releaseTimeFormat = "2006-01-02T15:04:05Z"
