package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/response"
)

type contextKey string

const IdempotencyKey contextKey = "idempotency_key"

func Idempotency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Idempotency-Key")

		if key == "" {
			response.SendResponse(w, response.SendResponseParams[any]{
				Status: http.StatusBadRequest,
				Error:  fmt.Errorf("header 'X-Idempotency-Key' is required for this operation"),
			})
			return
		}

		ctx := context.WithValue(r.Context(), IdempotencyKey, key)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
