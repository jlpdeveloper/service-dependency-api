package neo4jrepositories

import (
	"context"

	"service-atlas/databaseadapter"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ServiceFulltextIndexName is the name of the fulltext index used for service fuzzy search.
const ServiceFulltextIndexName = "service_fulltext_index"

// Startup ensures required database constructs exist (e.g., full-text indexes).
// It is safe to call multiple times; the Cypher uses IF NOT EXISTS for idempotency.
func Startup(ctx context.Context, driver neo4j.DriverWithContext) error {
	manager := databaseadapter.NewDriverManager(driver)

	// Create a full-text index on key Service fields commonly used for search.
	// Using Neo4j 5 syntax with IF NOT EXISTS for idempotency.
	_, err := manager.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, runErr := tx.Run(ctx, `
            CREATE FULLTEXT INDEX `+ServiceFulltextIndexName+` IF NOT EXISTS
            FOR (s:Service) ON EACH [s.name, s.description, s.type, s.url]
        `, nil)
		if runErr != nil {
			return nil, runErr
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}
