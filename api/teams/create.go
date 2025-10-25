package teams

import (
	"context"
	"encoding/json"
	"net/http"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
	"time"
)

func (c CallsHandler) CreateTeam(r *http.Request, rw http.ResponseWriter) {
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
	ctx_witn_timeout, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	err = c.Repository.CreateTeam(ctx_witn_timeout, *team)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}
