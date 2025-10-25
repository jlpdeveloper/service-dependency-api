package teamRepository

import (
	"context"
	"service-dependency-api/repositories"
	"strings"
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	neo4j_tc "github.com/testcontainers/testcontainers-go/modules/neo4j"
)

func TestNeo4jTeamRepository_CreateTeam(t *testing.T) {
	neo4jContainer, err := neo4j_tc.Run(context.Background(),
		"neo4j:latest",
		neo4j_tc.WithAdminPassword("letmein!"),
	)
	if err != nil {
		t.Fatal(err)
	}
	neo4jContainer.Start(context.Background())
	defer neo4jContainer.Terminate(context.Background())
	db_port, err := neo4jContainer.MappedPort(context.Background(), "7687")
	if err != nil {
		t.Fatal(err)
	}
	err = neo4jContainer.Start(context.Background())
	if err != nil {
	}
	endpoint, err := neo4jContainer.Endpoint(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	endpoint = "neo4j://" + strings.Split(endpoint, ":")[0] + ":" + db_port.Port()
	driver, err := neo4j.NewDriverWithContext(
		endpoint,
		neo4j.BasicAuth("neo4j", "letmein!", ""))
	if err != nil {
		t.Fatal(err)
	}
	defer driver.Close(context.Background())
	repo := New(driver)
	team := repositories.Team{
		Name: "test",
	}
	now := time.Now()
	err = repo.CreateTeam(context.Background(), team)

	session := driver.NewSession(context.Background(), neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	defer session.Close(context.Background())

	result, err := session.Run(context.Background(), "MATCH (n:Team) RETURN n.name as name, n.id as id, n.created as created", nil)
	if err != nil {
		t.Fatal(err)
	}

	returned_team, err := result.Single(context.Background())
	if err != nil || returned_team == nil {
		t.Fatal(err)
	}

	if n, _ := returned_team.Get("name"); n != team.Name {
		t.Errorf("Expected %s, got %s", team.Name, n)
	}
	if n, _ := returned_team.Get("id"); n == "" {
		t.Errorf("Expected ID, got %s", n)
	}
	if n, _ := returned_team.Get("created"); n.(time.Time).Before(now) || n.(time.Time).After(now.Add(time.Second*10)) {
		t.Errorf("Expected created time between %s and %s, got %s", now, now.Add(time.Second*10), n)
	}
}
