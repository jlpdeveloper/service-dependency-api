package routes

import (
	"log"
	"net/http"
	"service-dependency-api/api/hello_world"
	"service-dependency-api/api/system"
)

func SetupRoutes(mux *http.ServeMux) {
	log.Println("Setting up Routes")
	mux.HandleFunc("GET /helloworld", hello_world.HelloWorld)
	mux.HandleFunc("GET /time", system.GetTime)
	mux.HandleFunc("GET /database", system.GetDbAddress)
}
