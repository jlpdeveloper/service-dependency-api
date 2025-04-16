package main

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"service-dependency-api/api/routes"
	"service-dependency-api/internal/config"
)

func main() {
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		config.GetConfigValue("NEO4J_URL"),
		neo4j.BasicAuth(config.GetConfigValue("NEO4J_USERNAME"), config.GetConfigValue("NEO4J_PASSWORD"), ""))
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	routes.SetupRoutes(mux, &ctx, &driver)

	log.Println("Starting Web Server")
	log.Fatal(http.ListenAndServe(config.GetConfigValue("address"), mux))
}
