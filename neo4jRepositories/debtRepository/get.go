package debtRepository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"service-dependency-api/repositories"
)

func (n Neo4jDebtRepository) GetDebtByServiceId(ctx context.Context, id string) ([]repositories.Debt, error) {
	getByServiceIdTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (s:Service {id: $serviceId})-[:OWNS]->(d:Debt)
			RETURN d
		`, map[string]any{
			"serviceId": id,
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
			node, ok := record.Get("d")
			if ok {
				debtNode := node.(neo4j.Node)
				if title, found := debtNode.Props["title"]; found {
					debt.Title = title.(string)
				}

				if debtType, found := debtNode.Props["type"]; found {
					debt.Type = debtType.(string)
				}
				if description, found := debtNode.Props["description"]; found {
					debt.Description = description.(string)
				}

				if status, found := debtNode.Props["status"]; found {
					debt.Status = status.(string)
				}

				debtList = append(debtList, debt)
			}
		}
		return debtList, nil
	}
	debtList, err := n.manager.ExecuteRead(ctx, getByServiceIdTransaction)
	return debtList.([]repositories.Debt), err
}
