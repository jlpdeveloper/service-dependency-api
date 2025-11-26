package teamrepository

import (
	"service-atlas/databaseadapter"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jTeamRepository struct {
	manager databaseadapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jTeamRepository {
	return &Neo4jTeamRepository{manager: databaseadapter.NewDriverManager(driver)}
}
