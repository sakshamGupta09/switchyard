package health

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r *chi.Mux) {
	handler := NewHandler()
	r.Get("/health", handler.Health)
}
