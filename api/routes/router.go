package routes

import (
	"log"
	"net/http"
	"service-dependency-api/api/hello_world"
)

func SetupRoutes(mux *http.ServeMux) {
	log.Println("Setting up Routes")
	mux.HandleFunc("GET /helloworld", hello_world.HelloWorld)
}
