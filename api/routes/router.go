package routes

import (
	"log/slog"
	"net/http"
	"service-dependency-api/api/debt"
	"service-dependency-api/api/dependencies"
	"service-dependency-api/api/hello_world"
	"service-dependency-api/api/releases"
	"service-dependency-api/api/reports"
	"service-dependency-api/api/services"
	"service-dependency-api/api/system"
	"service-dependency-api/internal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func SetupRouter(driver *neo4j.DriverWithContext) http.Handler {
	slog.Debug("Setting up router")
	router := chi.NewRouter()

	router.Use(internal.StructuredLogger(slog.Default()))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))
	setupSystemCalls(router)
	services.Register(router, driver)
	dependencies.Register(router, driver)
	releases.Register(router, driver)
	debt.Register(router, driver)
	reports.Register(router, driver)
	return router
}

func setupSystemCalls(r *chi.Mux) {
	slog.Debug("Setting up system calls")
	r.Get("/time", system.GetTime)
	r.Get("/database", system.GetDbAddress)
	r.Get("/helloworld", hello_world.HelloWorld)
}
