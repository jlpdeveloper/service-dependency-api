package reportrepository

import (
	"context"
	"errors"
	"testing"

	"service-dependency-api/internal/customerrors"
	nRepo "service-dependency-api/neo4jrepositories"
	"service-dependency-api/neo4jrepositories/debtrepository"
	"service-dependency-api/neo4jrepositories/servicerepository"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestNeo4jReportRepository_GetServiceRiskReport_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx := context.Background()
	// spin up test container
	tc, err := nRepo.NewTestContainerHelper(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tc.Container.Terminate(ctx) })

	driver, err := neo4j.NewDriverWithContext(
		tc.Endpoint,
		neo4j.BasicAuth("neo4j", "letmein!", ""))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = driver.Close(ctx) }()

	reportRepo := New(driver)
	svcRepo := servicerepository.New(driver)
	debtRepo := debtrepository.New(driver)

	// Arrange: create target service
	targetID, err := svcRepo.CreateService(ctx, repositories.Service{
		Name:        "target",
		Description: "target service",
		ServiceType: "api",
		Url:         "https://target",
	})
	if err != nil {
		t.Fatalf("CreateService target error: %v", err)
	}

	// Arrange: create two dependents that depend on target
	dep1ID, err := svcRepo.CreateService(ctx, repositories.Service{Name: "dep-1", ServiceType: "worker", Url: "https://dep1"})
	if err != nil {
		t.Fatalf("CreateService dep1 error: %v", err)
	}
	dep2ID, err := svcRepo.CreateService(ctx, repositories.Service{Name: "dep-2", ServiceType: "worker", Url: "https://dep2"})
	if err != nil {
		t.Fatalf("CreateService dep2 error: %v", err)
	}

	write := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() { _ = write.Close(ctx) }()
	if _, err = write.Run(ctx, "MATCH (a:Service {id: $a}),(b:Service {id: $b}) MERGE (a)-[:DEPENDS_ON]->(b)", map[string]any{"a": dep1ID, "b": targetID}); err != nil {
		t.Fatalf("create dep1->target relationship: %v", err)
	}
	if _, err = write.Run(ctx, "MATCH (a:Service {id: $a}),(b:Service {id: $b}) MERGE (a)-[:DEPENDS_ON {version: '1.0.0'}]->(b)", map[string]any{"a": dep2ID, "b": targetID}); err != nil {
		t.Fatalf("create dep2->target relationship: %v", err)
	}

	// Arrange: create three debts for target: two security and one operational
	if err := debtRepo.CreateDebtItem(ctx, repositories.Debt{Type: "security", Title: "sec-1", Description: "", ServiceId: targetID}); err != nil {
		t.Fatalf("CreateDebtItem sec-1 error: %v", err)
	}
	if err := debtRepo.CreateDebtItem(ctx, repositories.Debt{Type: "security", Title: "sec-2", Description: "", ServiceId: targetID}); err != nil {
		t.Fatalf("CreateDebtItem sec-2 error: %v", err)
	}
	if err := debtRepo.CreateDebtItem(ctx, repositories.Debt{Type: "operational", Title: "op-1", Description: "", ServiceId: targetID}); err != nil {
		t.Fatalf("CreateDebtItem op-1 error: %v", err)
	}

	// Act
	report, err := reportRepo.GetServiceRiskReport(ctx, targetID)
	if err != nil {
		t.Fatalf("GetServiceRiskReport error: %v", err)
	}

	// Assert dependent count
	if report.DependentCount != 2 {
		t.Fatalf("expected DependentCount=2, got %d", report.DependentCount)
	}
	// Assert debt counts
	if report.DebtCount == nil {
		t.Fatalf("expected DebtCount map, got nil")
	}
	if got := report.DebtCount["security"]; got != 2 {
		t.Fatalf("expected security count=2, got %d", got)
	}
	if got := report.DebtCount["operational"]; got != 1 {
		t.Fatalf("expected operational count=1, got %d", got)
	}
}

func TestNeo4jReportRepository_GetServiceRiskReport_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx := context.Background()
	tc, err := nRepo.NewTestContainerHelper(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = tc.Container.Terminate(ctx) })

	driver, err := neo4j.NewDriverWithContext(
		tc.Endpoint,
		neo4j.BasicAuth("neo4j", "letmein!", ""))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = driver.Close(ctx) }()

	reportRepo := New(driver)

	_, err = reportRepo.GetServiceRiskReport(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Fatalf("expected error for non-existent service")
	}
	var httpErr *customerrors.HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *customerrors.HTTPError, got %T: %v", err, err)
	}
	if httpErr.Status != 404 {
		t.Fatalf("expected HTTP 404, got %d", httpErr.Status)
	}
}
