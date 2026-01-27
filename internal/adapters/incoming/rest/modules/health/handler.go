package health

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/response"
	healthApplication "perezvonish/plata-test-assignment/internal/application/health/usecases"
)

type Handler struct {
	pingUsecase healthApplication.HealthPingUsecase
}

func NewHandler() Handler {
	return Handler{pingUsecase: healthApplication.NewHealthPingUsecase()}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	result, err := h.pingUsecase.Execute(r.Context())
	if err != nil {
		response.SendResponse(w, response.SendResponseParams[healthApplication.HealthPingResult]{
			Status: http.StatusInternalServerError,
			Error:  err,
		})
		return
	}

	response.SendResponse(w, response.SendResponseParams[healthApplication.HealthPingResult]{
		Status: http.StatusOK,
		Data:   result,
	})
}
