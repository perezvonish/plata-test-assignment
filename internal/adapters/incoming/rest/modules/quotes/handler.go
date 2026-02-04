package quotes

import (
	"errors"
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/response"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
	"perezvonish/plata-test-assignment/internal/infrastructure/database/postgres/job"
	"perezvonish/plata-test-assignment/internal/infrastructure/database/postgres/quote"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	paramsValidator ParamsValidator

	updateUsecase        usecases.QuoteUpdateUsecase
	getByUpdateIdUsecase usecases.QuoteGetByUpdateIdUsecase
	getLatestUsecase     usecases.QuoteGetLatestUsecase
}

type HandlerInitParams struct {
	Pool *pgxpool.Pool

	JobChannel chan<- uuid.UUID
}

func newHandler(params HandlerInitParams) *Handler {
	paramsValidator := newParamsValidator()

	jobRepository := job.NewRepository(params.Pool)
	quoteRepository := quote.NewRepository(params.Pool)

	updateUsecase := usecases.NewQuoteUpdateUsecase(usecases.QuoteUpdateUsecaseInitParams{
		JobRepository:   jobRepository,
		JobChannel:      params.JobChannel,
		QuoteRepository: quoteRepository,
	})

	return &Handler{
		paramsValidator: paramsValidator,
		updateUsecase:   updateUsecase,
	}
}

// Update
// @Summary      Создать задачу на обновление котировки
// @Description  Принимает валютную пару, валидирует её и создает задачу (Job) на обновление цены.
// @Tags         quotes
// @Accept       json
// @Produce      json
// @Param        from    query     string  true  "UUID исходной валюты" format(string)
// @Param        to      query     string  true  "UUID целевой валюты"  format(string)
// @Param        X-Idempotency-Key  header    string  true  "Ключ идемпотентности"
// @Success      201      {object}  response.SendResponseParams[UpdateOutput] "Задача успешно создана"
// @Failure      400      {object}  response.SendResponseParams[any] "Ошибка валидации (одинаковые валюты, неверный формат и т.д.)"
// @Failure      500      {object}  response.SendResponseParams[any] "Внутренняя ошибка сервера"
// @Router       /v1/quotes [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	params, err := h.paramsValidator.ValidateUpdate(r.Context(), r)
	if err != nil {
		if errors.Is(err, ErrorFromAndToAreRequired) || errors.Is(err, ErrorInvalidCurrency) || errors.Is(err, ErrorIdenticalCurrency) {
			response.SendResponse(w, response.SendResponseParams[any]{
				Status: http.StatusBadRequest,
				Error:  err,
			})
			return
		}

		response.SendResponse(w, response.SendResponseParams[any]{
			Status: http.StatusInternalServerError,
			Error:  err,
		})
		return
	}

	updateId, err := h.updateUsecase.Execute(r.Context(), usecases.QuoteUpdateUsecaseInput{
		FromCurrency: params.From,
		ToCurrency:   params.To,
	})
	if err != nil {
		response.SendResponse(w, response.SendResponseParams[any]{
			Status: http.StatusInternalServerError,
			Error:  err,
		})
		return
	}

	response.SendResponse(w, response.SendResponseParams[UpdateOutput]{
		Status: http.StatusCreated,
		Data: UpdateOutput{
			UpdateId: updateId,
		},
	})
}

func (h *Handler) GetByUpdateId(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {}
