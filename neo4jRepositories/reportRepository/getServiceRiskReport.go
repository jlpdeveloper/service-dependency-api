package reportRepository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
	"sync"
)

func (n Neo4jReportRepository) GetServiceRiskReport(ctx context.Context, serviceId string) (*repositories.ServiceRiskReport, error) {
	_, err := n.manager.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Check if the service exists
		checkQuery := `
			MATCH (s:Service {id: $serviceId})
			RETURN s
		`
		result, err := tx.Run(ctx, checkQuery, map[string]any{
			"serviceId": serviceId,
		})
		if err != nil {
			return nil, err
		}

		// If no records are returned, the service doesn't exist
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			return nil, &customErrors.HTTPError{
				Status: 404,
				Msg:    fmt.Sprintf("Service not found: %s", serviceId),
			}
		}
		return nil, nil

	})
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	report := repositories.ServiceRiskReport{
		DebtCount: make(map[string]int64),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, dependentErr := n.manager.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			//get service dependent
			cypher := `
			MATCH (s:Service)-[r:DEPENDS_ON]->(d:Service {id: $serviceId})
			RETURN count(s) as count
			`
			result, err := tx.Run(ctx, cypher, map[string]any{
				"serviceId": serviceId,
			})
			if err != nil {
				return nil, err
			}
			if result.Next(ctx) {
				ctr, _ := result.Record().Get("count")
				if ct, ok := ctr.(int64); ok {
					report.DependentCount = ct
				}
			}
			if result.Err() != nil {

				return nil, result.Err()
			}
			return nil, nil

		})
		if dependentErr != nil {
			errCh <- dependentErr
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, debtErr := n.manager.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			//get map of debt types
			cypher := `
			MATCH (s:Service)-[r:OWNS]->(d:Debt)
			WHERE s.id = $serviceId
			RETURN d.type as type, count(d) as count
			`
			result, err := tx.Run(ctx, cypher, map[string]any{
				"serviceId": serviceId,
			})
			if err != nil {
				return nil, err
			}
			for result.Next(ctx) {
				record := result.Record().AsMap()
				report.DebtCount[record["type"].(string)] = record["count"].(int64)
			}
			if result.Err() != nil {
				return nil, result.Err()
			}
			return nil, nil
		})
		if debtErr != nil {
			errCh <- debtErr
		}
	}()

	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}
	return &report, nil
}
