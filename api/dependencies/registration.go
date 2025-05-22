package dependencies

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"service-dependency-api/neo4jRepositories/dependencyRepository"
	"service-dependency-api/repositories"
)

type ServiceCallsHandler struct {
	Repository repositories.DependencyRepository
}

func (s *ServiceCallsHandler) Register(mux *http.ServeMux) {
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

func Register(mux *http.ServeMux, driver *neo4j.DriverWithContext) {
	handler := ServiceCallsHandler{
		Repository: dependencyRepository.New(*driver),
	}
	handler.Register(mux)

}
