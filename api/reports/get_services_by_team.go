package reports

import (
	"encoding/json"
	"log"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customerrors"
)

func (c *CallsHandler) GetServicesByTeam(rw http.ResponseWriter, r *http.Request) {
	teamId, ok := internal.GetGuidFromRequestPath("teamId", r)
	if !ok {
		http.Error(rw, "Invalid team ID", http.StatusBadRequest)
		return
	}
	services, err := c.repository.GetServicesByTeam(r.Context(), teamId)
	if err != nil {
		customerrors.HandleError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(services)
	if err != nil {
		log.Println(err)
	}

}
