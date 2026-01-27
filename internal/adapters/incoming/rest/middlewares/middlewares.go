package middlewares

import "net/http"

func Use(sm *http.ServeMux) http.Handler {
	return useTrace(useLogger(sm))
}
