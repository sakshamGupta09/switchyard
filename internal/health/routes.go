package health

import (
	"switchyard/internal/app"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(deps app.Dependencies) *chi.Mux {
	r := chi.NewRouter()

	service := NewService(deps.MongoClient)
	handler := NewHandler(service)

	r.Get("/health", handler.HealthCheck)
	r.Get("/readiness", handler.ReadinessCheck)

	return r
}
