package teams

import (
	"encoding/json"
	"log"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
)

func (c CallsHandler) GetTeam(rw http.ResponseWriter, r *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", r)
	if !ok {
		http.Error(rw, "Invalid team ID", http.StatusBadRequest)
		return
	}
	team, err := c.Repository.GetTeam(r.Context(), id)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(team)
	if err != nil {
		log.Println(err)
	}

}
