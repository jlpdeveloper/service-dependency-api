package tests

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/api/services"
	"service-dependency-api/internal/database"
	"testing"
)

// --- Mock Dependencies ---

type mockSession struct {
	data *map[string]any
}

func (m *mockSession) ExecuteWrite(ctx context.Context, fn database.ManagedTransactionWork) (any, error) {
	return fn(&mockTx{
		data: m.data,
	})
}

func (m *mockSession) Close(ctx context.Context) error {
	return nil
}

type mockTx struct {
	data *map[string]any
}

func (t *mockTx) Run(ctx context.Context, query string, params map[string]any) (database.Neo4jResult, error) {
	return &mockResult{
		data: t.data,
	}, nil
}

type mockResult struct {
	data *map[string]any
}

func (r *mockResult) Single(ctx context.Context) (database.Neo4jRecord, error) {
	return &mockRecord{
		data: r.data,
	}, nil
}

type mockRecord struct {
	data *map[string]any
}

func (r *mockRecord) AsMap() map[string]any {
	return *r.data
}

type mockDriver struct {
	data map[string]any
}

func (d *mockDriver) NewSession(ctx context.Context, config neo4j.SessionConfig) database.Neo4jSession {
	return &mockSession{
		data: &d.data,
	}
}

func TestCreateService(t *testing.T) {
	svc := &services.ServiceNeo4jService{
		Driver: &mockDriver{
			data: map[string]any{
				"id":   "mock-id-123",
				"name": "Test Service",
			},
		},
		Ctx: context.Background(),
	}

	service := services.Service{
		Name:        "MockService",
		ServiceType: "Internal",
		Description: "Unit test service",
	}

	id, err := svc.CreateService(service)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if id == "" {
		t.Fatal("expected a valid UUID, got empty string")
	}
}
