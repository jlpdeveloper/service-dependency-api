package services

import (
	"context"
	"log/slog"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customerrors"
	"time"
)

func (u *ServiceCallsHandler) DeleteServiceById(rw http.ResponseWriter, r *http.Request) {
	logger := internal.LoggerFromContext(r.Context())
	id, ok := internal.GetGuidFromRequestPath("id", r)
	logger.Debug("Request received - DeleteServiceById - " + id)
	if !ok {
		http.Error(rw, "Invalid Request", http.StatusBadRequest)
		logger.Debug("Invalid Request - " + id)
		return
	}
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	err := u.Repository.DeleteService(ctxWithTimeout, id)
	if err != nil {
		logger.Debug("Error deleting service:",
			slog.String("error", err.Error()))
		customerrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusNoContent)

}
