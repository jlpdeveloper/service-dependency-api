package debtrepository

import (
	"context"
	"net/http"
	"service-atlas/internal/customerrors"
	"service-atlas/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (n Neo4jDebtRepository) GetDebtByServiceId(ctx context.Context, id string, page, pageSize int, onlyResolved bool) ([]repositories.Debt, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, &customerrors.HTTPError{
			Status: http.StatusBadRequest,
			Msg:    "page and page size must be positive integers",
		}
	}
	getByServiceIdTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := `
			MATCH (s:Service {id: $serviceId})-[:OWNS]->(d:Debt)
		`
		params := map[string]any{
			"serviceId": id,
			"skip":      (page - 1) * pageSize,
			"limit":     pageSize,
		}
		if onlyResolved {
			cypher += `
				WHERE d.status = $status
			`
			params["status"] = "remediated"
		}
		cypher += `
			RETURN d.title as title, d.description as description, d.type as type, d.status as status, d.id as id
			ORDER BY d.created DESC
			SKIP $skip
			LIMIT $limit`
		result, err := tx.Run(ctx, cypher, params)
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

			if id, ok := record.Get("id"); ok && id != nil {
				debt.Id = id.(string)
			} else {
				return nil, &customerrors.HTTPError{
					Status: 500,
					Msg:    "debt id is nil",
				}
			}

			debtList = append(debtList, debt)

		}

		if err := result.Err(); err != nil {
			return nil, err
		}
		return debtList, nil
	}

	debtList, err := n.manager.ExecuteRead(ctx, getByServiceIdTransaction)
	if err != nil {
		return nil, err
	}
	if typedDebtList, ok := debtList.([]repositories.Debt); ok {
		return typedDebtList, nil
	}
	return nil, &customerrors.HTTPError{
		Status: http.StatusInternalServerError,
		Msg:    "unexpected return type from transaction",
	}
}
