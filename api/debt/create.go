package debt

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customerrors"
	"service-dependency-api/repositories"
	"time"
)

func (c CallsHandler) CreateDebt(rw http.ResponseWriter, r *http.Request) {
	logger := internal.LoggerFromContext(r.Context())
	id, ok := internal.GetGuidFromRequestPath("id", r)
	if !ok {
		http.Error(rw, "service id not valid", http.StatusBadRequest)
		return
	}
	debt := &repositories.Debt{}
	const maxBodySize = 1 << 20 // 1 MB
	r.Body = http.MaxBytesReader(rw, r.Body, maxBodySize)
	err := json.NewDecoder(r.Body).Decode(debt)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	debt.ServiceId = id

	if err = debt.Validate(); err != nil {
		logger.Debug("Invalid debt item:", slog.String("error", err.Error()))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctxWithTimeout, cancel := context.WithTimeoutCause(r.Context(), 15*time.Second, errors.New("database timeout"))
	defer cancel()
	if err = c.Repository.CreateDebtItem(ctxWithTimeout, *debt); err != nil {
		logger.Debug("Error creating debt item:",
			slog.String("error", err.Error()))
		customerrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
