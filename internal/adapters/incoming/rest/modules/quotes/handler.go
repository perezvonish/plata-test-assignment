package quotes

import (
	"net/http"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
)

type Handler struct {
	updateUsecase        usecases.QuoteUpdateUsecase
	getByUpdateIdUsecase usecases.QuoteGetByUpdateIdUsecase
	getLatestUsecase     usecases.QuoteGetLatestUsecase
}

func newHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetByUpdateId(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {}
