package teamRepository

import (
	"context"
	"service-dependency-api/neo4jRepositories"
	"service-dependency-api/repositories"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Using create_test.go as an example, this test verifies the happy-path for GetTeam.
func TestNeo4jTeamRepository_GetTeam(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	ctx := context.Background()

	// Start Neo4j test container
	tc, err := neo4jRepositories.NewTestContainerHelper(ctx)
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

	// Create a team we can fetch
	team := repositories.Team{Name: "get-team-test"}
	now := time.Now()
	if _, err := repo.CreateTeam(ctx, team); err != nil {
		t.Fatal(err)
	}

	// Read session to fetch the created team's id by name (deterministic lookup)
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer func() { _ = session.Close(ctx) }()

	res, err := session.Run(ctx,
		"MATCH (n:Team {name: $name}) RETURN n.id AS id, n.name AS name, n.created AS created",
		map[string]any{"name": team.Name},
	)
	if err != nil {
		t.Fatal(err)
	}
	rec, err := res.Single(ctx)
	if err != nil || rec == nil {
		t.Fatalf("expected single record, got error: %v", err)
	}
	idVal, ok := rec.Get("id")
	if !ok {
		t.Fatalf("missing 'id' in created team record")
	}
	idStr, ok := idVal.(string)
	if !ok || idStr == "" {
		t.Fatalf("invalid 'id' type/value: %v (%T)", idVal, idVal)
	}

	// Exercise the method under test
	got, err := repo.GetTeam(ctx, idStr)
	if err != nil {
		t.Fatalf("GetTeam returned error: %v", err)
	}
	if got == nil {
		t.Fatalf("GetTeam returned nil team")
	}

	// Validate core fields
	if got.Id != idStr {
		t.Errorf("expected Id %q, got %q", idStr, got.Id)
	}
	if got.Name != team.Name {
		t.Errorf("expected Name %q, got %q", team.Name, got.Name)
	}

	// created should be within a reasonable window of now
	if got.Created.Before(now) || got.Created.After(now.Add(10*time.Second)) {
		t.Errorf("expected Created between %s and %s, got %s", now, now.Add(10*time.Second), got.Created)
	}
}
