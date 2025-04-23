package serviceRepository

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

// ServiceNeo4jRepository type implements the interface for service repository above
type ServiceNeo4jRepository struct {
	Driver neo4j.DriverWithContext
}

// mapNodeToService converts a Neo4j node to a Service object
func (d *ServiceNeo4jRepository) mapNodeToService(n neo4j.Node) Service {
	svc := Service{}

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
