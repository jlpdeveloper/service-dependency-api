package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"service-dependency-api/api/routes"
	"service-dependency-api/internal/config"
	"syscall"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
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

	mux := routes.SetupRouter(&driver)

	server := &http.Server{
		Handler: mux,
		Addr:    config.GetConfigValue("address"),
	}

	log.Println("Starting Web Server")
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
