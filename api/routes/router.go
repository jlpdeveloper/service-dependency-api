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
	"service-dependency-api/api/teams"
	"service-dependency-api/internal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func SetupRouter(driver neo4j.DriverWithContext) http.Handler {
	slog.Debug("Setting up router")
	router := chi.NewRouter()

	router.Use(internal.StructuredLogger(slog.Default()))
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))
	setupSystemCalls(router)

	serviceHandler := services.New(driver)
	debtHandler := debt.New(driver)
	dependencyHandler := dependencies.New(driver)
	releaseHandler := releases.New(driver)
	reportHandler := reports.New(driver)
	teamHandler := teams.New(driver)

	router.Get("/releases/{startDate}/{endDate}", releaseHandler.GetReleasesInDateRange)
	router.Get("/reports/services/{id}/risk", reportHandler.GetServiceRiskReport)
	router.Patch("/debt/{id}", debtHandler.UpdateDebtStatus)

	router.Route("/services", func(r chi.Router) {
		r.Get("/", serviceHandler.GetAllServices)
		r.Post("/", serviceHandler.CreateService)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", serviceHandler.GetById)
			r.Put("/", serviceHandler.UpdateService)
			r.Delete("/", serviceHandler.DeleteServiceById)

			r.Get("/dependencies", dependencyHandler.GetDependencies)
			r.Get("/dependents", dependencyHandler.GetDependents)
			r.Post("/dependency", dependencyHandler.CreateDependency)
			r.Delete("/dependency/{id2}", dependencyHandler.DeleteDependency)

			r.Route("/debt", func(r chi.Router) {
				r.Post("/", debtHandler.CreateDebt)
				r.Get("/", debtHandler.GetDebtByServiceId)
			})

			r.Route("/release", func(r chi.Router) {
				r.Post("/", releaseHandler.CreateRelease)
				r.Get("/", releaseHandler.GetReleasesByServiceId)
			})

		})
	})

	router.Route("/teams", func(r chi.Router) {
		r.Post("/", teamHandler.CreateTeam)
		r.Delete("/{id}", teamHandler.DeleteTeam)
		r.Get("/{id}", teamHandler.GetTeam)
		r.Get("/", teamHandler.GetTeams)
		r.Put("/{id}", teamHandler.UpdateTeam)
	})
	return router
}

func setupSystemCalls(r chi.Router) {
	slog.Debug("Setting up system calls")
	r.Get("/time", system.GetTime)
	r.Get("/database", system.GetDbAddress)
	r.Get("/helloworld", hello_world.HelloWorld)
}
