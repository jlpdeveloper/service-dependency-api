package reportRepository

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/repositories"
	"sync"
)

func (n Neo4jReportRepository) GetServiceRiskReport(ctx context.Context, serviceId string) (*repositories.ServiceRiskReport, error) {
	result, err := n.manager.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		var wg sync.WaitGroup
		errCh := make(chan error, 2)
		report := repositories.ServiceRiskReport{}
		wg.Add(1)
		//get service dependents
		go func() {
			defer wg.Done()
			cypher := `
			MATCH (s:Service)-[r:DEPENDS_ON]->(d:Service {id: $serviceId})
			RETURN count(d) as count
			`
			result, err := tx.Run(ctx, cypher, map[string]any{
				"serviceId": serviceId,
			})
			if err != nil {
				errCh <- err
				return
			}
			if result.Next(ctx) {
				ctr, _ := result.Record().Get("count")
				if ct, ok := ctr.(int); ok {
					report.DependentCount = ct
				}
			}
			if result.Err() != nil {
				errCh <- result.Err()
				return
			}

		}()
		wg.Add(1)
		//get map of debt types
		go func() {
			defer wg.Done()
			cypher := `
			MATCH (s:Service)-[r:OWNS]->(d:Debt)
			WHERE s.id = $serviceId
			RETURN d.type as type, count(d) as count
			`
			result, err := tx.Run(ctx, cypher, map[string]any{
				"serviceId": serviceId,
			})
			if err != nil {
				errCh <- err
				return
			}
			for result.Next(ctx) {
				record := result.Record().AsMap()
				report.DebtCount[record["type"].(string)] = record["count"].(int)
			}
			if result.Err() != nil {
				errCh <- result.Err()
				return
			}

		}()

		wg.Wait()
		close(errCh)
		for err := range errCh {
			if err != nil {
				return nil, err
			}
		}
		return report, nil

	})
	if err != nil {
		return nil, err
	}
	convertedReport, ok := result.(repositories.ServiceRiskReport)
	if !ok {
		return nil, errors.New("failed to convert result to ServiceRiskReport type")
	}
	return &convertedReport, nil
}
