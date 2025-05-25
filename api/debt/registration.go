package debt

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/neo4jRepositories/debtRepository"
	"service-dependency-api/repositories"
)

type CallsHandler struct {
	Repository repositories.DebtRepository
}

func (c CallsHandler) Register(mux *http.ServeMux) {
	calls := map[string]func(http.ResponseWriter, *http.Request){
		"POST /services/{id}/debt": c.createDebt,
		"GET /services/{id}/debt":  c.getDebtByServiceId,
		"PATCH /debt/{id}":         c.updateDebtStatus,
	}
	for path, f := range calls {
		mux.HandleFunc(path, f)
	}
}

func Register(mux *http.ServeMux, driver *neo4j.DriverWithContext) {
	handler := CallsHandler{
		Repository: debtRepository.New(*driver),
	}
	handler.Register(mux)
}
