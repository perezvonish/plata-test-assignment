package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestTraceId struct{}

var requestId = requestTraceId{}

func useTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New()
		ctx := context.WithValue(r.Context(), requestId, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
