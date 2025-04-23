package serviceRepository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"service-dependency-api/internal"
)

func (d *ServiceNeo4jRepository) DeleteService(ctx context.Context, id string) (err error) {
	session := d.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		closeErr := session.Close(ctx)
		if err == nil {
			err = closeErr
		}
	}()
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
				return nil, &internal.HTTPError{
					Status: 404,
					Msg:    "Service not found",
				}
			}
		}
		result, err = tx.Run(ctx, `
		MATCH(s:Service { id: $id})
		DELETE s;`, map[string]interface{}{"id": id})
		if err != nil {
			log.Println("Error deleting service: " + id)
			return nil, err
		}

		summary, err := result.Consume(ctx)
		if err != nil {
			return nil, &internal.HTTPError{Status: 500, Msg: "Error deleting service: " + id}
		}

		if summary.Counters().NodesDeleted() == 0 {
			log.Println("Error deleting service: " + id + ". Database transaction not successful")
			return nil, &internal.HTTPError{Status: 500, Msg: "Error deleting service: " + id}
		}
		return nil, nil
	}

	_, err = session.ExecuteWrite(ctx, deleteServiceTransaction)
	return err
}
