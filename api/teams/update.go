package teams

import (
	"encoding/json"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
	"service-dependency-api/repositories"
)

func (c CallsHandler) UpdateTeam(rw http.ResponseWriter, r *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", r)
	if !ok {
		http.Error(rw, "Invalid team ID", http.StatusBadRequest)
		return
	}
	team := &repositories.Team{}
	err := json.NewDecoder(r.Body).Decode(team)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if err := team.Validate(); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	if id != team.Id {
		http.Error(rw, "Invalid team ID", http.StatusBadRequest)
		return
	}
	err = c.Repository.UpdateTeam(r.Context(), *team)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}
