package quotes

import (
	"errors"
	"net/http"
	"perezvonish/plata-test-assignment/internal/adapters/incoming/rest/response"
	applicationQuote "perezvonish/plata-test-assignment/internal/application/quote"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
)

type Handler struct {
	paramsValidator ParamsValidator

	updateUsecase        usecases.QuoteUpdateUsecase
	getByUpdateIdUsecase usecases.QuoteGetByUpdateIdUsecase
	getLatestUsecase     usecases.QuoteGetLatestUsecase
}

type HandlerInitParams struct {
	UpdateUsecase        usecases.QuoteUpdateUsecase
	GetByUpdateIdUsecase usecases.QuoteGetByUpdateIdUsecase
	GetLatestUsecase     usecases.QuoteGetLatestUsecase
}

func newHandler(params HandlerInitParams) *Handler {
	paramsValidator := NewParamsValidator()

	return &Handler{
		paramsValidator:      paramsValidator,
		updateUsecase:        params.UpdateUsecase,
		getByUpdateIdUsecase: params.GetByUpdateIdUsecase,
		getLatestUsecase:     params.GetLatestUsecase,
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

// GetByUpdateId
// @Summary      Получение котировки по идентификатору задачи на обновление
// @Description  Принимает id задачи на обновлние
// @Tags         quotes
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "UUID исходной валюты" format(uuid)
// @Success      200      {object}  response.SendResponseParams[GetByUpdateIdOutput] ""
// @Failure      400      {object}  response.SendResponseParams[any] "Ошибка валидации (не корректный идентификатор задачи, не передан id и т.д.)"
// @Failure      404      {object}  response.SendResponseParams[any] "Задачи с таким идентификатором не существует"
// @Failure      500      {object}  response.SendResponseParams[any] "Внутренняя ошибка сервера"
// @Router       /v1/quotes/update-task/{id} [get]
func (h *Handler) GetByUpdateId(w http.ResponseWriter, r *http.Request) {
	params, err := h.paramsValidator.ValidateGetByUpdateId(r.Context(), r)
	if err != nil {
		if errors.Is(err, ErrorNotPassedJobId) || errors.Is(err, ErrorNotValidJobId) {
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

	result, err := h.getByUpdateIdUsecase.Execute(r.Context(), usecases.QuoteGetByUpdateIdUsecaseInput{
		UpdateId: params.Id,
	})
	if err != nil {
		if errors.Is(err, applicationQuote.ErrorWhileFindingJob) || errors.Is(err, applicationQuote.ErrorNotFoundJob) {
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

	response.SendResponse(w, response.SendResponseParams[GetByUpdateIdOutput]{
		Status: http.StatusOK,
		Data: GetByUpdateIdOutput{
			Job: result,
		},
	})
}

// GetLatest
// @Summary      Получение последнего обновления котировки
// @Tags         quotes
// @Accept       json
// @Produce      json
// @Param        from    query     string  true  "UUID исходной валюты" format(string)
// @Param        to      query     string  true  "UUID целевой валюты"  format(string)
// @Success      200      {object}  response.SendResponseParams[UpdateOutput] "Получено последнее обновление котировки"
// @Failure      400      {object}  response.SendResponseParams[any] "Ошибка валидации (одинаковые валюты, неверный формат и т.д.)"
// @Failure      404      {object}  response.SendResponseParams[any] "Котировка еще ни разу не обновлялась"
// @Failure      500      {object}  response.SendResponseParams[any] "Внутренняя ошибка сервера"
// @Router       /v1/quotes/latest [get]
func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {
	params, err := h.paramsValidator.ValidateGetLatest(r.Context(), r)
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

	result, err := h.getLatestUsecase.Execute(r.Context(), usecases.QuoteGetLatestUsecaseParams{
		From: params.From,
		To:   params.To,
	})
	if err != nil {
		if errors.Is(err, applicationQuote.ErrorWhileFindingJob) || errors.Is(err, applicationQuote.ErrorNotFoundJob) {
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

	response.SendResponse(w, response.SendResponseParams[GetLatestOutput]{
		Status: http.StatusOK,
		Data: GetLatestOutput{
			Job: result,
		},
	})
}
