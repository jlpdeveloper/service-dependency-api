package teams

import (
	"context"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customerrors"
	"time"
)

func (c CallsHandler) CreateTeamAssociation(rw http.ResponseWriter, r *http.Request) {
	teamId, ok := internal.GetGuidFromRequestPath("teamId", r)
	if !ok {
		http.Error(rw, "Invalid team ID", http.StatusBadRequest)
		return
	}
	serviceId, ok := internal.GetGuidFromRequestPath("serviceId", r)
	if !ok {
		http.Error(rw, "Invalid service ID", http.StatusBadRequest)
		return
	}
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	err := c.Repository.CreateTeamAssociation(ctxWithTimeout, teamId, serviceId)
	if err != nil {
		customerrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func (c CallsHandler) DeleteTeamAssociation(rw http.ResponseWriter, r *http.Request) {
	teamId, ok := internal.GetGuidFromRequestPath("teamId", r)
	if !ok {
		http.Error(rw, "Invalid team ID", http.StatusBadRequest)
		return
	}
	serviceId, ok := internal.GetGuidFromRequestPath("serviceId", r)
	if !ok {
		http.Error(rw, "Invalid service ID", http.StatusBadRequest)
		return
	}
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	err := c.Repository.DeleteTeamAssociation(ctxWithTimeout, teamId, serviceId)
	if err != nil {
		customerrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}
