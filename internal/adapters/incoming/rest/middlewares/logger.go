package middlewares

import (
	"log"
	"net/http"
	"perezvonish/plata-test-assignment/internal/shared/utils"
	"time"
)

func useLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := utils.GenerateUUID()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		log.Printf("%s | %s %s | %d | %v\n",
			requestID,
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
		)
	})
}
