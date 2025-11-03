package reportRepository

import (
	"context"
	"net/http"
	"service-dependency-api/internal/customerrors"
	nRepo "service-dependency-api/neo4jRepositories"
	"service-dependency-api/repositories"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (r Neo4jReportRepository) GetServicesByTeam(ctx context.Context, teamId string) ([]repositories.Service, error) {
	getServicesByTeamTransaction := func(tx neo4j.ManagedTransaction) (any, error) {
		cypher := `
		MATCH (t:Team {id: $teamId}) -[r:OWNS]-> (s:Service)
		RETURN s
		`
		result, err := tx.Run(ctx, cypher, map[string]any{
			"teamId": teamId,
		})
		if err != nil {
			return nil, customerrors.HTTPError{
				Status: http.StatusInternalServerError,
				Msg:    err.Error(),
			}
		}
		services := []repositories.Service{}
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

			svc := nRepo.MapNodeToService(n)
			services = append(services, svc)
		}
		return services, nil
	}

	servicesResult, err := r.manager.ExecuteRead(ctx, getServicesByTeamTransaction)
	if err != nil {
		return nil, err
	}
	services, ok := servicesResult.([]repositories.Service)
	if !ok {
		return nil, customerrors.HTTPError{
			Status: http.StatusInternalServerError,
			Msg:    "unexpected return type from transaction",
		}
	}
	return services, nil
}
