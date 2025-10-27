package teamRepository

import (
	"context"
	"errors"
	"net/http"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/neo4jRepositories"
	"service-dependency-api/repositories"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestNeo4jTeamRepository_DeleteTeam(t *testing.T) {
	ctx := context.Background()
	tc, err := neo4jRepositories.NewTestContainerHelper(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = tc.Container.Terminate(ctx)
	})

	driver, err := neo4j.NewDriverWithContext(
		tc.Endpoint,
		neo4j.BasicAuth("neo4j", "letmein!", ""))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = driver.Close(ctx)
	}()
	repo := New(driver)
	team := repositories.Team{
		Name: "test",
	}
	id, err := repo.CreateTeam(ctx, team)

	if err != nil {
		t.Fatal(err)
	}

	err = repo.DeleteTeam(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	session := driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})
	defer func() {
		_ = session.Close(ctx)
	}()
	result, err := session.Run(ctx, "MATCH (t:Team) RETURN count(t) as ctr", map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	ctrResult, err := result.Single(ctx)
	if err != nil {
		t.Fatalf("expected single team record, got error: %v", err)
	}
	ctr, ok := ctrResult.Get("ctr")
	if !ok {
		t.Fatalf("missing 'ctr' field in record")
	}
	if ctr.(int64) != 0 {
		t.Fatalf("expected 0 teams, got %d", ctr.(int64))
	}

}

func TestNeo4jTeamRepository_DeleteTeam_Returns404(t *testing.T) {
	ctx := context.Background()
	tc, err := neo4jRepositories.NewTestContainerHelper(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = tc.Container.Terminate(ctx)
	})

	driver, err := neo4j.NewDriverWithContext(
		tc.Endpoint,
		neo4j.BasicAuth("neo4j", "letmein!", ""))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = driver.Close(ctx)
	}()
	repo := New(driver)
	team := repositories.Team{
		Name: "test",
	}
	_, err = repo.CreateTeam(ctx, team)

	if err != nil {
		t.Fatal(err)
	}
	//using team name here was purposeful, to ensure that the team is not found by id
	err = repo.DeleteTeam(ctx, team.Name)
	if err == nil {
		t.Fatal("expected error")
	}
	var cErr *customErrors.HTTPError
	errors.As(err, &cErr)
	if cErr.Error() != "Team not found" {
		t.Errorf("expected error message 'team not found', got '%s'", err.Error())
	}
	if cErr.Status != http.StatusNotFound {
		t.Errorf("expected status code 404, got %d", cErr.Status)
	}

}
