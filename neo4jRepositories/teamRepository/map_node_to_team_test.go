package teamRepository

import (
	"testing"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Test_mapNodeToTeam(t *testing.T) {
	var r Neo4jTeamRepository // zero value is fine; method doesn't use manager

	now := time.Now().UTC().Truncate(time.Microsecond)
	later := now.Add(1 * time.Hour).UTC().Truncate(time.Microsecond)

	tests := []struct {
		name        string
		node        neo4j.Node
		wantName    string
		wantId      string
		wantCreated time.Time
		wantUpdated time.Time
	}{
		{
			name: "all properties present with correct types",
			node: neo4j.Node{Props: map[string]any{
				"name":    "team-a",
				"id":      "abc-123",
				"created": now,
				"updated": later,
			}},
			wantName:    "team-a",
			wantId:      "abc-123",
			wantCreated: now,
			wantUpdated: later,
		},
		{
			name: "missing optional properties are zero-valued",
			node: neo4j.Node{Props: map[string]any{
				"name": "only-name",
			}},
			wantName:    "only-name",
			wantId:      "",
			wantCreated: time.Time{},
			wantUpdated: time.Time{},
		},
		{
			name: "incorrect types are ignored (leave zero values)",
			node: neo4j.Node{Props: map[string]any{
				"name":    123,          // not a string
				"id":      456,          // not a string
				"created": "2021-01-01", // not time.Time
				"updated": struct{}{},   // not time.Time
			}},
			wantName:    "",
			wantId:      "",
			wantCreated: time.Time{},
			wantUpdated: time.Time{},
		},
		{
			name: "extra properties are ignored",
			node: neo4j.Node{Props: map[string]any{
				"name":    "extra",
				"id":      "id-1",
				"created": now,
				"updated": later,
				"foo":     "bar",
			}},
			wantName:    "extra",
			wantId:      "id-1",
			wantCreated: now,
			wantUpdated: later,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.mapNodeToTeam(tt.node)

			if got.Name != tt.wantName {
				t.Errorf("Name: expected %q, got %q", tt.wantName, got.Name)
			}
			if got.Id != tt.wantId {
				t.Errorf("Id: expected %q, got %q", tt.wantId, got.Id)
			}
			// Created
			if tt.wantCreated.IsZero() {
				if !got.Created.IsZero() {
					t.Errorf("Created: expected zero value, got %v", got.Created)
				}
			} else if !got.Created.Equal(tt.wantCreated) {
				t.Errorf("Created: expected %v, got %v", tt.wantCreated, got.Created)
			}
			// Updated
			if tt.wantUpdated.IsZero() {
				if !got.Updated.IsZero() {
					t.Errorf("Updated: expected zero value, got %v", got.Updated)
				}
			} else if !got.Updated.Equal(tt.wantUpdated) {
				t.Errorf("Updated: expected %v, got %v", tt.wantUpdated, got.Updated)
			}
		})
	}
}
