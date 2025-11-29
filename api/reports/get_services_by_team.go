package reports

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"service-atlas/internal"
	"service-atlas/internal/customerrors"
)

func (c *CallsHandler) GetServicesByTeam(rw http.ResponseWriter, r *http.Request) {
	teamId, ok := internal.GetGuidFromRequestPath("teamId", r)
	logger := internal.LoggerFromContext(r.Context())
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
		logger.Debug("Error encoding services json",
			slog.String("error", err.Error()),
		)
	}

}
