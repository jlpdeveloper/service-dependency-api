package dependencies

import (
	"net/http"
	"service-dependency-api/neo4jRepositories/dependencyRepository"
	"service-dependency-api/repositories"

	"github.com/go-chi/chi/v5"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ServiceCallsHandler struct {
	Repository repositories.DependencyRepository
}

func (s *ServiceCallsHandler) Register(mux *chi.Mux) {
	paths := map[string]func(http.ResponseWriter, *http.Request){
		"POST /services/{id}/dependency":         s.createDependency,
		"GET /services/{id}/dependencies":        s.getDependencies,
		"GET /services/{id}/dependents":          s.getDependents,
		"DELETE /services/{id}/dependency/{id2}": s.deleteDependency,
	}
	for k, v := range paths {
		mux.HandleFunc(k, v)
	}
}

func Register(mux *chi.Mux, driver *neo4j.DriverWithContext) {
	handler := ServiceCallsHandler{
		Repository: dependencyRepository.New(*driver),
	}
	handler.Register(mux)

}
