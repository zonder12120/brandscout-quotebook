package middleware

import (
	"net/http"
	"time"

	"github.com/zonder12120/brandscout-quotebook/pkg/logger"
)

func Logging(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lw := &loggingResponseWriter{w, http.StatusOK}
			next.ServeHTTP(lw, r)

			log.Info().Msgf(
				"method=%s path=%s status=%d duration=%s",
				r.Method,
				r.URL.Path,
				lw.status,
				time.Since(start),
			)
		})
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}
