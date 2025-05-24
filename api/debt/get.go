package debt

import (
	"encoding/json"
	"net/http"
	"service-dependency-api/internal"
	"service-dependency-api/internal/customErrors"
)

func (c CallsHandler) getDebtByServiceId(rw http.ResponseWriter, r *http.Request) {
	id, ok := internal.GetGuidFromRequestPath("id", r)
	if !ok {
		http.Error(rw, "service id not valid", http.StatusBadRequest)
		return
	}
	debt, err := c.Repository.GetDebtByServiceId(r.Context(), id)
	if err != nil {
		customErrors.HandleError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(debt)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
