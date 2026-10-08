package health

import (
	"switchyard/internal/app"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(deps app.RouteDeps) *chi.Mux {
	r := chi.NewRouter()

	service := NewService(deps.MongoClient)
	handler := NewHandler(service)

	r.Get("/health", handler.HealthCheck)
	r.Get("/ready", handler.ReadinessCheck)

	return r
}
