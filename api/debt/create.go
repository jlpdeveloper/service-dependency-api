package debt

import (
	"context"
	"encoding/json"
	"net/http"
	"service-dependency-api/api/debt/internal/debtRepository"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
	"time"
)

func (c CallsHandler) createDebt(rw http.ResponseWriter, r *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", r)
	if !ok {
		http.Error(rw, "service id not valid", http.StatusBadRequest)
		return
	}
	debt := &debtRepository.Debt{}
	const maxBodySize = 1 << 20 // 1 MB
	r.Body = http.MaxBytesReader(rw, r.Body, maxBodySize)
	err := json.NewDecoder(r.Body).Decode(debt)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	debt.ServiceId = id

	if err = debt.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	if err = c.Repository.CreateDebtItem(ctxWithTimeout, *debt); err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
