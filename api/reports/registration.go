package reports

import (
	"net/http"
	"service-dependency-api/neo4jRepositories/reportRepository"
	"service-dependency-api/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CallsHandler struct {
	repository repositories.ReportRepository
}

func (c *CallsHandler) Register(mux *chi.Mux) {
	paths := map[string]func(http.ResponseWriter, *http.Request){
		"GET /reports/services/{id}/risk": c.getServiceRiskReport,
	}
	for k, v := range paths {
		mux.HandleFunc(k, v)
	}
}

func Register(mux *chi.Mux, driver *neo4j.DriverWithContext) {
	handler := CallsHandler{
		repository: reportRepository.New(*driver),
	}
	handler.Register(mux)
}
