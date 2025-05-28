package main

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"service-dependency-api/api/routes"
	"service-dependency-api/internal/config"
	"service-dependency-api/middleware"
)

func main() {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		config.GetConfigValue("NEO4J_URL"),
		neo4j.BasicAuth(config.GetConfigValue("NEO4J_USERNAME"), config.GetConfigValue("NEO4J_PASSWORD"), ""))
	defer func() {
		closeErr := driver.Close(ctx)
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	routes.SetupRoutes(mux, &driver)

	middlewareHandler := middleware.CreateStack(
		middleware.Logging,
	)

	server := &http.Server{
		Handler: middlewareHandler(mux),
		Addr:    config.GetConfigValue("address"),
	}

	log.Println("Starting Web Server")
	log.Fatal(server.ListenAndServe())
}
