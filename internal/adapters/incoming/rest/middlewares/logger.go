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

		timestamp := start.Format("02/01/2006 15:04:05.000")

		log.Printf("%s | %s | %s %s | %d | %v\n",
			requestID,
			timestamp,
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
		)
	})
}
