package teamRepository

import (
	"service-dependency-api/databaseAdapter"
	"service-dependency-api/repositories"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jTeamRepository struct {
	manager databaseAdapter.DriverManager
}

func New(driver neo4j.DriverWithContext) *Neo4jTeamRepository {
	return &Neo4jTeamRepository{manager: databaseAdapter.NewDriverManager(driver)}
}

func (r Neo4jTeamRepository) mapNodeToTeam(n neo4j.Node) repositories.Team {
	team := repositories.Team{}

	// Safely extract name with validation
	if name, ok := n.Props["name"]; ok {
		if nameStr, ok := name.(string); ok {
			team.Name = nameStr
		}
	}

	// Safely extract ID with validation
	if id, ok := n.Props["id"]; ok {
		if idStr, ok := id.(string); ok {
			team.Id = idStr
		}
	}

	// Safely extract created date with validation
	if date, ok := n.Props["created"]; ok {
		if dateStr, ok := date.(time.Time); ok {
			team.Created = dateStr
		}
	}

	if date, ok := n.Props["updated"]; ok {
		if dateStr, ok := date.(time.Time); ok {
			team.Updated = dateStr
		}
	}
	return team
}
