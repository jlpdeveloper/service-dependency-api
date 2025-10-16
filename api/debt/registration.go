package debt

import (
	"net/http"
	"service-dependency-api/neo4jRepositories/debtRepository"
	"service-dependency-api/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CallsHandler struct {
	Repository repositories.DebtRepository
}

func (c CallsHandler) Register(mux *chi.Mux) {
	calls := map[string]func(http.ResponseWriter, *http.Request){
		"POST /services/{id}/debt": c.createDebt,
		"GET /services/{id}/debt":  c.getDebtByServiceId,
		"PATCH /debt/{id}":         c.updateDebtStatus,
	}
	for path, f := range calls {
		mux.HandleFunc(path, f)
	}
}

func Register(mux *chi.Mux, driver *neo4j.DriverWithContext) {
	handler := CallsHandler{
		Repository: debtRepository.New(*driver),
	}
	handler.Register(mux)
}
