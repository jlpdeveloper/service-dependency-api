package routes

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"service-dependency-api/api/dependencies"
	"service-dependency-api/api/hello_world"
	"service-dependency-api/api/releases"
	"service-dependency-api/api/services"
	"service-dependency-api/api/system"
)

func SetupRoutes(mux *http.ServeMux, driver *neo4j.DriverWithContext) {
	log.Println("Setting up Routes")
	mux.HandleFunc("GET /helloworld", hello_world.HelloWorld)
	mux.HandleFunc("GET /time", system.GetTime)
	mux.HandleFunc("GET /database", system.GetDbAddress)

	services.Register(mux, driver)
	dependencies.Register(mux, driver)
	releases.Register(mux, driver)
}
