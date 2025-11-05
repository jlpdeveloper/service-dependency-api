package neo4jrepositories

import (
	"service-dependency-api/repositories"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// MapNodeToService converts a Neo4j node to a Service object
func MapNodeToService(n neo4j.Node) repositories.Service {
	svc := repositories.Service{}

	// Safely extract name with validation
	if name, ok := n.Props["name"]; ok {
		if nameStr, ok := name.(string); ok {
			svc.Name = nameStr
		}
	}

	// Safely extract description with validation
	if desc, ok := n.Props["description"]; ok {
		if descStr, ok := desc.(string); ok {
			svc.Description = descStr
		}
	}

	// Safely extract service type with validation
	if svcType, ok := n.Props["type"]; ok {
		if typeStr, ok := svcType.(string); ok {
			svc.ServiceType = typeStr
		}
	}

	// Safely extract ID with validation
	if id, ok := n.Props["id"]; ok {
		if idStr, ok := id.(string); ok {
			svc.Id = idStr
		}
	}

	if url, ok := n.Props["url"]; ok {
		if urlStr, ok := url.(string); ok {
			svc.Url = urlStr
		}
	}

	// Safely extract created date with validation
	if date, ok := n.Props["created"]; ok {
		if dateStr, ok := date.(time.Time); ok {
			svc.Created = dateStr
		}
	}

	if date, ok := n.Props["updated"]; ok {
		if dateStr, ok := date.(time.Time); ok {
			svc.Updated = dateStr
		}
	}
	return svc
}

// MapNodeToTeam converts a Neo4j node to a Team object
func MapNodeToTeam(n neo4j.Node) repositories.Team {
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
