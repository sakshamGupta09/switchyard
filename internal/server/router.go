package server

import (
	"switchyard/internal/app"
	"switchyard/internal/health"
	"switchyard/internal/middleware"

	"github.com/go-chi/chi/v5"
)

const defaultMaxRequestBodySize = 1 << 20

func NewRouter(deps app.RouteDeps) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.AccessLog)
	r.Use(middleware.Recoverer)
	r.Use(middleware.MaxRequestBodySize(defaultMaxRequestBodySize)) // 1 MB

	// Register health check routes
	r.Mount("/", health.RegisterRoutes(deps))

	return r
}
