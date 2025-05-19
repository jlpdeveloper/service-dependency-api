package debt

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/api/debt/internal/debtRepository"
)

type CallsHandler struct {
	Repository debtRepository.Repository
}

func (c CallsHandler) Register(mux *http.ServeMux) {
	calls := map[string]func(http.ResponseWriter, *http.Request){
		"POST /services/{id}/debt": c.createDebt,
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
