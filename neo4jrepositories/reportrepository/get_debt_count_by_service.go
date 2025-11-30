package reportrepository

import (
	"context"
	"service-atlas/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jReportRepository) GetDebtCountByService(ctx context.Context) ([]repositories.ServiceDebtReport, error) {
	debtReport := make([]repositories.ServiceDebtReport, 0)
	work := func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := `
        MATCH (s:Service)-[:OWNS]->(d:Debt)
        WHERE d.status IN $statuses
        RETURN s.name AS name, s.id AS id, count(d) AS count
        ORDER BY count DESC, name ASC
        `
		params := map[string]any{
			"statuses": []string{"in_progress", "pending"},
		}
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		for result.Next(ctx) {
			rec := result.Record()
			nameVal, _ := rec.Get("name")
			idVal, _ := rec.Get("id")
			countVal, _ := rec.Get("count")

			name, _ := nameVal.(string)
			id, _ := idVal.(string)
			var countInt int
			if c64, ok := countVal.(int64); ok {
				countInt = int(c64)
			}
			debtReport = append(debtReport, repositories.ServiceDebtReport{
				Name:  name,
				Id:    id,
				Count: countInt,
			})
		}
		if err := result.Err(); err != nil {
			return nil, err
		}
		return nil, nil
	}
	_, err := r.manager.ExecuteRead(ctx, work)
	if err != nil {
		return nil, err
	}
	return debtReport, nil
}
