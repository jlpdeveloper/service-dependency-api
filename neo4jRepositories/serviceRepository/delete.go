package serviceRepository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"service-dependency-api/internal/customErrors"
)

func (d *Neo4jServiceRepository) DeleteService(ctx context.Context, id string) (err error) {
	deleteServiceTransaction := func(tx neo4j.ManagedTransaction) (any, error) {

		result, err := tx.Run(ctx, `
    		MATCH (s:Service { id: $id })
    		RETURN count(s) as count
		`, map[string]interface{}{"id": id})

		if err != nil {
			return nil, err
		}

		if record, err := result.Single(ctx); err == nil {
			count, _ := record.Get("count")
			if count.(int64) == 0 {
				return nil, &customErrors.HTTPError{
					Status: 404,
					Msg:    "Service not found",
				}
			}
		}
		result, err = tx.Run(ctx, `
		MATCH(s:Service { id: $id})
		DETACH DELETE s;`, map[string]interface{}{"id": id})
		if err != nil {
			log.Println("Error deleting service: " + id)
			return nil, err
		}

		summary, err := result.Consume(ctx)
		if err != nil {
			return nil, &customErrors.HTTPError{Status: 500, Msg: "Error deleting service: " + id}
		}

		if summary.Counters().NodesDeleted() == 0 {
			log.Println("Error deleting service: " + id + ". Database transaction not successful")
			return nil, &customErrors.HTTPError{Status: 500, Msg: "Error deleting service: " + id}
		}
		return nil, nil
	}

	_, err = d.manager.ExecuteWrite(ctx, deleteServiceTransaction)
	return err
}
