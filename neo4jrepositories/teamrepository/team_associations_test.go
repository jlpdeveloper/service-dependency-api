package teamrepository

import (
	"context"
	"testing"

	"service-atlas/neo4jrepositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestNeo4jTeamRepository_CreateTeamAssociation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx := context.Background()

	// Start Neo4j test container
	tc, err := neo4jrepositories.NewTestContainerHelper(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = tc.Container.Terminate(ctx)
	})

	// Connect driver
	driver, err := neo4j.NewDriverWithContext(
		tc.Endpoint,
		neo4j.BasicAuth("neo4j", "letmein!", ""),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = driver.Close(ctx) }()

	repo := New(driver)

	// Arrange: create Team and Service nodes with GUID-like ids
	teamID := "11111111-1111-1111-1111-111111111111"
	serviceID := "22222222-2222-2222-2222-222222222222"

	write := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() { _ = write.Close(ctx) }()

	_, err = write.Run(ctx,
		"CREATE (t:Team {id: $tid, name: $tname}) RETURN t",
		map[string]any{"tid": teamID, "tname": "assoc-team"},
	)
	if err != nil {
		t.Fatalf("failed to create team node: %v", err)
	}
	_, err = write.Run(ctx,
		"CREATE (s:Service {id: $sid, name: $sname}) RETURN s",
		map[string]any{"sid": serviceID, "sname": "assoc-service"},
	)
	if err != nil {
		t.Fatalf("failed to create service node: %v", err)
	}

	// Act
	if err := repo.CreateTeamAssociation(ctx, teamID, serviceID); err != nil {
		t.Fatalf("CreateTeamAssociation returned error: %v", err)
	}

	// Assert: relationship exists
	read := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() { _ = read.Close(ctx) }()

	res, err := read.Run(ctx,
		"MATCH (t:Team {id: $tid})-[:OWNS]->(s:Service {id: $sid}) RETURN count(*) AS c",
		map[string]any{"tid": teamID, "sid": serviceID},
	)
	if err != nil {
		t.Fatalf("failed to verify relationship: %v", err)
	}
	rec, err := res.Single(ctx)
	if err != nil || rec == nil {
		t.Fatalf("expected single record verifying relationship, got err=%v", err)
	}
	cVal, _ := rec.Get("c")
	switch v := cVal.(type) {
	case int64:
		if v != 1 {
			t.Fatalf("expected 1 relationship, got %d", v)
		}
	default:
		t.Fatalf("unexpected count type %T: %v", cVal, cVal)
	}
}

func TestNeo4jTeamRepository_DeleteTeamAssociation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx := context.Background()

	// Start Neo4j test container
	tc, err := neo4jrepositories.NewTestContainerHelper(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = tc.Container.Terminate(ctx)
	})

	// Connect driver
	driver, err := neo4j.NewDriverWithContext(
		tc.Endpoint,
		neo4j.BasicAuth("neo4j", "letmein!", ""),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = driver.Close(ctx) }()

	repo := New(driver)

	// Arrange: create Team, Service, and the OWNS relationship
	teamID := "33333333-3333-3333-3333-333333333333"
	serviceID := "44444444-4444-4444-4444-444444444444"

	write := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() { _ = write.Close(ctx) }()

	if _, err = write.Run(ctx,
		"CREATE (t:Team {id: $tid, name: $tname}) RETURN t",
		map[string]any{"tid": teamID, "tname": "assoc-team-delete"},
	); err != nil {
		t.Fatalf("failed to create team node: %v", err)
	}
	if _, err = write.Run(ctx,
		"CREATE (s:Service {id: $sid, name: $sname}) RETURN s",
		map[string]any{"sid": serviceID, "sname": "assoc-service-delete"},
	); err != nil {
		t.Fatalf("failed to create service node: %v", err)
	}
	if _, err = write.Run(ctx,
		"MATCH (t:Team {id: $tid}), (s:Service {id: $sid}) CREATE (t)-[:OWNS]->(s)",
		map[string]any{"tid": teamID, "sid": serviceID},
	); err != nil {
		t.Fatalf("failed to create relationship: %v", err)
	}

	// Act
	if err := repo.DeleteTeamAssociation(ctx, teamID, serviceID); err != nil {
		// Note: If implementation doesn't RETURN a record on delete, this may fail
		// and should be fixed in the repository. The test asserts desired behavior.
		t.Fatalf("DeleteTeamAssociation returned error: %v", err)
	}

	// Assert: relationship no longer exists
	read := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() { _ = read.Close(ctx) }()

	res, err := read.Run(ctx,
		"MATCH (t:Team {id: $tid})-[r:OWNS]->(s:Service {id: $sid}) RETURN count(r) AS c",
		map[string]any{"tid": teamID, "sid": serviceID},
	)
	if err != nil {
		t.Fatalf("failed to verify relationship deletion: %v", err)
	}
	rec, err := res.Single(ctx)
	if err != nil || rec == nil {
		t.Fatalf("expected single record verifying deletion, got err=%v", err)
	}
	cVal, _ := rec.Get("c")
	switch v := cVal.(type) {
	case int64:
		if v != 0 {
			t.Fatalf("expected 0 relationships after delete, got %d", v)
		}
	default:
		t.Fatalf("unexpected count type %T: %v", cVal, cVal)
	}
}
