package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"service-dependency-api/api/routes"
	"service-dependency-api/internal/config"
	"strings"
	"syscall"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	ctx := context.Background()
	logger := getLogger()
	slog.SetDefault(logger)
	driver, err := neo4j.NewDriverWithContext(
		config.GetConfigValue("DB_URL"),
		neo4j.BasicAuth(config.GetConfigValue("DB_USERNAME"), config.GetConfigValue("DB_PASSWORD"), ""))
	defer func() {
		closeErr := driver.Close(ctx)
		if closeErr != nil {
			slog.Error("error closing driver: ", slog.Any("error", closeErr))
		}
	}()
	if err != nil {
		slog.Error("Error creating driver: ", slog.Any("error", err))
		os.Exit(1)
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	mux := routes.SetupRouter(driver)

	server := &http.Server{
		Handler: mux,
		Addr:    config.GetConfigValue("address"),
	}

	slog.Info("Starting Web Server")
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen error", slog.Any("error", err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", slog.Any("error", err))
	}
}

func getLogger() *slog.Logger {
	lvlEnv, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		lvlEnv = "info"
	}
	switch strings.ToLower(lvlEnv) {
	case "debug":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "error":
		return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))
	case "warning":
		return slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}))
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))

}
