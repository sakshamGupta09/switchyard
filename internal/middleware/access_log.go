package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(w, r)

		slog.Info("HTTP Request summary",
			"request_id", GetRequestID(r.Context()),
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", recorder.Status(),
			"duration", time.Since(start),
		)
	})
}
