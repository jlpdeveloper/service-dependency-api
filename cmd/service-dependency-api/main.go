package main

import (
	"log"
	"net/http"
	"service-dependency-api/api/routes"
	"service-dependency-api/internal/config"
)

func main() {
	mux := http.NewServeMux()
	routes.SetupRoutes(mux)

	log.Println("Starting Web Server")
	log.Fatal(http.ListenAndServe(config.GetConfigValue("address"), mux))
}
