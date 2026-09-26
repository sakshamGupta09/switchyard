package server

import (
	"switchyard/internal/health"

	"github.com/go-chi/chi/v5"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	health.RegisterRoutes(r)

	return r
}
