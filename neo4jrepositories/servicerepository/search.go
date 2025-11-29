package servicerepository

import (
	"context"
	nRepo "service-atlas/neo4jrepositories"
	"service-atlas/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Search performs a fuzzy search against the Service full-text index
// and returns matching services ordered by relevance.
func (d *Neo4jServiceRepository) Search(ctx context.Context, query string) ([]repositories.Service, error) {
	services := make([]repositories.Service, 0)
	if query == "" {
		return services, nil
	}
	if string(query[len(query)-1]) != "~" {
		query += "~"
	}
	work := func(tx neo4j.ManagedTransaction) (any, error) {
		localServices := make([]repositories.Service, 0)
		result, err := tx.Run(ctx, `
            CALL db.index.fulltext.queryNodes($indexName, $q)
            YIELD node, score
            RETURN node AS s, score
            ORDER BY score DESC
            LIMIT 50
        `, map[string]any{
			"indexName": nRepo.ServiceFulltextIndexName,
			"q":         query,
		})
		if err != nil {
			return nil, err
		}

		for result.Next(ctx) {
			record := result.Record()
			node, ok := record.Get("s")
			if !ok {
				continue
			}
			n, ok := node.(neo4j.Node)
			if !ok {
				continue
			}
			localServices = append(localServices, nRepo.MapNodeToService(n))
		}
		services = localServices
		return nil, nil
	}

	if _, err := d.manager.ExecuteRead(ctx, work); err != nil {
		return nil, err
	}
	return services, nil
}
