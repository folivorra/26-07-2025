package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler { // TODO: обернуть хэндлеры в мидлварь
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger.Info("incoming request",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("remote", r.RemoteAddr),
			)

			next.ServeHTTP(w, r)

			logger.Info("request completed",
				slog.String("method", r.Method),
				slog.String("url", r.URL.Path),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}
