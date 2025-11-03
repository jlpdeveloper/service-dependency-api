package teams

import (
	"context"
	"encoding/json"
	"net/http"
	"service-dependency-api/internal/customerrors"
	"service-dependency-api/repositories"
	"time"
)

func (c CallsHandler) CreateTeam(rw http.ResponseWriter, r *http.Request) {
	team := &repositories.Team{}
	const maxBodySize = 1 << 20 // 1 MB
	r.Body = http.MaxBytesReader(rw, r.Body, maxBodySize)
	err := json.NewDecoder(r.Body).Decode(team)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err = team.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	id, err := c.Repository.CreateTeam(ctxWithTimeout, *team)
	if err != nil {
		customerrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
	_, _ = rw.Write([]byte(id))
}
