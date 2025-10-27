package teams

import (
	"context"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
	"time"
)

func (c CallsHandler) DeleteTeam(rw http.ResponseWriter, r *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", r)
	if !ok {
		http.Error(rw, "Invalid team ID", http.StatusBadRequest)
		return
	}
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	err := c.Repository.DeleteTeam(ctxWithTimeout, id)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}
