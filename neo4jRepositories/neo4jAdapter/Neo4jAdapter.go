package neo4jAdapter

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/databaseAdapter"
)

func NewDriverManager(driver neo4j.DriverWithContext) databaseAdapter.DriverManager {
	return &Neo4jDriverAdapter{driver}
}

// Neo4jSessionAdapter adapts Neo4j session to SessionManager
type Neo4jSessionAdapter struct {
	session neo4j.SessionWithContext
}

func (n Neo4jSessionAdapter) ExecuteWrite(ctx context.Context, work func(tx databaseAdapter.TransactionManager) (any, error)) (any, error) {
	return n.session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		adapter := Neo4jTransactionAdapter{tx}
		return work(adapter)
	})
}

func (n Neo4jSessionAdapter) Close(ctx context.Context) error {
	return n.session.Close(ctx)
}

// Neo4jDriverAdapter adapts Neo4j driver to DriverManager
type Neo4jDriverAdapter struct {
	driver neo4j.DriverWithContext
}

func (n Neo4jDriverAdapter) NewSession(ctx context.Context, config databaseAdapter.SessionConfig) databaseAdapter.SessionManager {
	am := neo4j.AccessModeRead
	if config.AccessMode == "write" {
		am = neo4j.AccessModeWrite
	}
	return &Neo4jSessionAdapter{
		n.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: am}),
	}
}

type Neo4jTransactionAdapter struct {
	tx neo4j.ManagedTransaction
}

func (n Neo4jTransactionAdapter) Run(ctx context.Context, cypher string, params map[string]any) (databaseAdapter.ResultManager, error) {
	results, err := n.tx.Run(ctx, cypher, params)
	return &Neo4jResultAdapter{result: results}, err
}

func (n Neo4jTransactionAdapter) ReturnBase() any {
	return n.tx
}

type Neo4jResultAdapter struct {
	result neo4j.ResultWithContext
}

func (n Neo4jResultAdapter) Collect(ctx context.Context) (any, error) {
	records, err := n.result.Collect(ctx)
	return records, err
}

func (n Neo4jResultAdapter) Single(ctx context.Context) (any, error) {
	record, err := n.result.Single(ctx)
	return record, err
}

func (n Neo4jResultAdapter) HasRecords(ctx context.Context) (bool, error) {
	records, err := n.result.Collect(ctx)
	if err != nil {
		return false, err
	}
	return len(records) > 0, nil
}

func (n Neo4jResultAdapter) Next(ctx context.Context) bool {
	return n.result.Next(ctx)
}

func (n Neo4jResultAdapter) GetProperty(ctx context.Context, key string) (any, error) {
	record, err := n.result.Single(ctx)
	if err != nil {
		return nil, err
	}
	val, ok := record.Get(key)
	if !ok {
		return nil, nil
	}
	return val, nil
}
