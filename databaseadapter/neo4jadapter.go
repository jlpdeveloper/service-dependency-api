package databaseadapter

import (
	"context"
	"log/slog"
	"service-dependency-api/internal"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NewDriverManager(driver neo4j.DriverWithContext) DriverManager {
	return &Neo4jDriverAdapter{driver}
}

// Neo4jDriverAdapter adapts Neo4j driver to DriverManager
type Neo4jDriverAdapter struct {
	driver neo4j.DriverWithContext
}

func (n Neo4jDriverAdapter) executeInSession(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error), mode neo4j.AccessMode) (any, error) {
	session := n.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: mode})
	logger := internal.LoggerFromContext(ctx)
	defer func() {
		err := session.Close(ctx)
		if err != nil {
			logger.Error("Error closing session: ",
				slog.String("error", err.Error()),
			)
		}
	}()
	requestId := internal.GetRequestIdFromContext(ctx)
	if mode == neo4j.AccessModeWrite {
		return session.ExecuteWrite(ctx, work,
			neo4j.WithTxMetadata(map[string]any{
				"requestId": requestId,
			}))
	}
	return session.ExecuteRead(ctx, work, neo4j.WithTxMetadata(map[string]any{
		"requestId": requestId,
	}))
}

func (n Neo4jDriverAdapter) ExecuteWrite(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error)) (any, error) {
	return n.executeInSession(ctx, work, neo4j.AccessModeWrite)
}

func (n Neo4jDriverAdapter) ExecuteRead(ctx context.Context, work func(tx neo4j.ManagedTransaction) (any, error)) (any, error) {
	return n.executeInSession(ctx, work, neo4j.AccessModeRead)
}
