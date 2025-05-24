package debtRepository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
)

func (n Neo4jDebtRepository) GetDebtByServiceId(ctx context.Context, id string, page, pageSize int) ([]repositories.Debt, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, &customErrors.HTTPError{
			Status: http.StatusBadRequest,
			Msg:    "page and page size must be positive integers",
		}
	}
	getByServiceIdTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (s:Service {id: $serviceId})-[:OWNS]->(d:Debt)
			RETURN d.title as title, d.description as description, d.type as type, d.status as status
			ORDER BY d.created DESC
			SKIP $skip
			LIMIT $limit
		`, map[string]any{
			"serviceId": id,
			"skip":      (page - 1) * pageSize,
			"limit":     pageSize,
		})
		if err != nil {
			return nil, err
		}
		debtList := make([]repositories.Debt, 0)
		for result.Next(ctx) {
			record := result.Record()
			debt := repositories.Debt{
				ServiceId: id,
			}
			if title, ok := record.Get("title"); ok && title != nil {
				debt.Title = title.(string)
			}
			if debtType, ok := record.Get("type"); ok && debtType != nil {
				debt.Type = debtType.(string)
			}
			if description, ok := record.Get("description"); ok && description != nil {
				debt.Description = description.(string)
			}
			if status, ok := record.Get("status"); ok && status != nil {
				debt.Status = status.(string)
			}

			debtList = append(debtList, debt)

		}
		return debtList, nil
	}
	debtList, err := n.manager.ExecuteRead(ctx, getByServiceIdTransaction)
	return debtList.([]repositories.Debt), err
}
