package database

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jWrapper struct {
	Driver neo4j.DriverWithContext
}

func (r *Neo4jWrapper) NewSession(ctx context.Context, config neo4j.SessionConfig) Neo4jSession {
	return &SessionWrapper{r.Driver.NewSession(ctx, config)}
}

type SessionWrapper struct {
	session neo4j.SessionWithContext
}

func (s *SessionWrapper) ExecuteWrite(ctx context.Context, fn ManagedTransactionWork) (any, error) {
	return s.session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		return fn(&TransactionWrapper{tx})
	})
}

func (s *SessionWrapper) Close(ctx context.Context) error {
	return s.session.Close(ctx)
}

type TransactionWrapper struct {
	tx neo4j.ManagedTransaction
}

func (t *TransactionWrapper) Run(ctx context.Context, query string, params map[string]any) (Neo4jResult, error) {
	result, err := t.tx.Run(ctx, query, params)
	if err != nil {
		return nil, err
	}
	return &ResultWrapper{result}, nil
}

type ResultWrapper struct {
	result neo4j.ResultWithContext
}

func (r *ResultWrapper) Single(ctx context.Context) (Neo4jRecord, error) {
	rec, err := r.result.Single(ctx)
	if err != nil {
		return nil, err
	}
	return &RecordWrapper{*rec}, nil
}

type RecordWrapper struct {
	rec neo4j.Record
}

func (r *RecordWrapper) AsMap() map[string]any {
	return r.rec.AsMap()
}
