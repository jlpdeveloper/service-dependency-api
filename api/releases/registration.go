package releases

import (
	"net/http"
	"service-dependency-api/neo4jRepositories/releaseRepository"
	"service-dependency-api/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ServiceCallsHandler struct {
	Repository repositories.ReleaseRepository
}

func (s *ServiceCallsHandler) Register(mux *chi.Mux) {
	paths := map[string]func(http.ResponseWriter, *http.Request){
		"POST /services/{id}/release":         s.createRelease,
		"GET /services/{id}/releases":         s.getReleasesByServiceId,
		"GET /releases/{startDate}/{endDate}": s.getReleasesInDateRange,
	}
	for k, v := range paths {
		mux.HandleFunc(k, v)
	}
}

func Register(mux *chi.Mux, driver *neo4j.DriverWithContext) {
	handler := ServiceCallsHandler{
		Repository: releaseRepository.New(*driver),
	}
	handler.Register(mux)

}
