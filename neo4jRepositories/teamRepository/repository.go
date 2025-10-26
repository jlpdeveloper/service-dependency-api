package teamRepository

import (
	"service-dependency-api/databaseAdapter"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jTeamRepository struct {
	manager databaseAdapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jTeamRepository {
	return &Neo4jTeamRepository{manager: databaseAdapter.NewDriverManager(driver)}
}
