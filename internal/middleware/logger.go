package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type resposeWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *resposeWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Loggingmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		ww := &resposeWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		slog.Info("incoming requst",
			"method", r.Method,
			"path", r.URL.Path,
			"remoteaddr", r.RemoteAddr,
			"agent", r.UserAgent(),
		)

		next.ServeHTTP(ww, r)

		slog.Info("requst completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", ww.statusCode,
			"duration", time.Since(start).String(),
		)
	})
}
