package usecases_test

import (
	"context"
	"fmt"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"perezvonish/plata-test-assignment/internal/domain/job"
	"perezvonish/plata-test-assignment/internal/domain/quote"

	svcMocks "perezvonish/plata-test-assignment/mocks/application/quote/services"
	jobMocks "perezvonish/plata-test-assignment/mocks/domain/job"
	quoteMocks "perezvonish/plata-test-assignment/mocks/domain/quote"
)

func TestProcessQuoteUpdateJobUsecase_Execute(t *testing.T) {
	ctx := context.Background()
	jobID := uuid.New()
	quoteID := uuid.New()

	// Тестовые данные
	mockJob := &job.Job{
		Id:      jobID,
		QuoteId: quoteID,
		Status:  job.StatusPending,
	}
	mockQuote := &quote.Quote{
		Id:           quoteID,
		FromCurrency: "USD",
		ToCurrency:   "RUB",
	}

	rate := 92.5
	expectedPriceE8 := int64(9250000000) // 92.5 * 10^8

	t.Run("Success", func(t *testing.T) {
		// Инициализация моков
		jobRepo := jobMocks.NewRepository(t)
		quoteRepo := quoteMocks.NewRepository(t)
		priceSvc := svcMocks.NewExchangePrice(t)

		u := usecases.NewProcessQuoteUpdateUsecase(usecases.ProcessQuoteUpdateUsecaseInitParams{
			// Здесь передаем моки. NewProcessQuoteUpdateUsecase нужно будет
			// слегка подправить, чтобы он принимал интерфейсы, а не инициализировал их сам,
			// либо создать конструктор специально для тестов.
		})

		// ВАЖНО: В реальном коде лучше передавать зависимости в структуру напрямую,
		// чтобы тест мог подменить их моками.

		// 1. Get Job
		jobRepo.On("GetById", ctx, jobID).Return(mockJob, nil).Once()

		// 2. Mark as Processing
		jobRepo.On("UpdateStatus", ctx, mock.MatchedBy(func(p job.UpdateStatusParams) bool {
			return p.Id == jobID && p.Status == job.StatusProcessing
		})).Return(nil).Once()

		// 3. Get Quote
		quoteRepo.On("GetById", ctx, quoteID).Return(mockQuote, nil).Once()

		// 4. Get Rate
		priceSvc.On("GetRate", ctx, mock.Anything).Return(rate, nil).Once()

		// 5. Update Job Price
		jobRepo.On("UpdatePrice", ctx, mock.MatchedBy(func(p job.UpdatePriceParams) bool {
			return p.Id == jobID && p.PriceE8Rate == expectedPriceE8
		})).Return(nil).Once()

		// 6. Update Quote Price
		quoteRepo.On("UpdatePrice", ctx, quote.UpdatePriceParams{
			Id:    quoteID,
			Price: expectedPriceE8,
		}).Return(nil).Once()

		// 7. Success Status
		jobRepo.On("UpdateStatus", ctx, mock.MatchedBy(func(p job.UpdateStatusParams) bool {
			return p.Status == job.StatusSuccess
		})).Return(nil).Once()

		// Вызов
		// (Предполагается, что зависимости внедрены в usecaseImpl)
		err := u.Execute(ctx, usecases.ProcessQuoteUpdateJobUsecaseParams{Id: jobID})

		assert.NoError(t, err)
	})

	t.Run("Failure_On_External_Service", func(t *testing.T) {
		jobRepo := jobMocks.NewRepository(t)
		quoteRepo := quoteMocks.NewRepository(t)
		priceSvc := svcMocks.NewExchangePrice(t)

		// Те же шаги до ошибки сервиса
		jobRepo.On("GetById", ctx, jobID).Return(mockJob, nil).Once()
		jobRepo.On("UpdateStatus", ctx, mock.Anything).Return(nil).Once()
		quoteRepo.On("GetById", ctx, quoteID).Return(mockQuote, nil).Once()

		// Сервис возвращает ошибку
		priceSvc.On("GetRate", ctx, mock.Anything).Return(0.0, fmt.Errorf("api down")).Once()

		// Проверяем, что статус стал Failure
		jobRepo.On("UpdateStatus", ctx, mock.MatchedBy(func(p job.UpdateStatusParams) bool {
			return p.Status == job.StatusFailure
		})).Return(nil).Once()

		u := &usecases.ProcessQuoteUpdateJobUsecaseImpl{
			jobRepository:        jobRepo,
			quoteRepository:      quoteRepo,
			exchangePriceService: priceSvc,
		}

		err := u.Execute(ctx, usecases.ProcessQuoteUpdateJobUsecaseParams{Id: jobID})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "external api error")
	})
}
