package releases

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/api/releases/internal/releaseRepository"
)

type ServiceCallsHandler struct {
	Repository releaseRepository.ReleaseRepository
}

func (s *ServiceCallsHandler) Register(mux *http.ServeMux) {
	paths := map[string]func(http.ResponseWriter, *http.Request){
		"POST /services/{id}/release":         s.createRelease,
		"GET /services/{id}/releases":         s.getReleasesByServiceId,
		"GET /releases/{startDate}/{endDate}": s.getReleasesInDateRange,
	}
	for k, v := range paths {
		mux.HandleFunc(k, v)
	}
}

func Register(mux *http.ServeMux, driver *neo4j.DriverWithContext) {
	repo := releaseRepository.New(*driver)

	handler := ServiceCallsHandler{
		Repository: repo,
	}
	handler.Register(mux)

}
