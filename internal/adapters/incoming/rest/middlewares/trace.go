package middlewares

import (
	"context"
	"net/http"
	"perezvonish/plata-test-assignment/internal/shared/utils"
)

type requestTraceId struct{}

var requestId = requestTraceId{}

func useTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := utils.GenerateUUID()
		ctx := context.WithValue(r.Context(), requestId, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
