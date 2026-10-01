package server

import (
	"switchyard/internal/middleware"

	"github.com/go-chi/chi/v5"
)

const defaultMaxRequestBodySize = 1 << 20

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Limit the size of the request body to 1 MB
	r.Use(middleware.MaxRequestBodySize(defaultMaxRequestBodySize)) // 1 MB

	return r
}
