package usecases_test

import (
	"context"
	"fmt"
	"perezvonish/plata-test-assignment/internal/application/quote/usecases"
	mocksServices "perezvonish/plata-test-assignment/mocks/application/quote/services"
	mocksRepositories "perezvonish/plata-test-assignment/mocks/domain/job"
	mocks "perezvonish/plata-test-assignment/mocks/domain/quote"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"perezvonish/plata-test-assignment/internal/domain/job"
	"perezvonish/plata-test-assignment/internal/domain/quote"
)

func TestProcessQuoteUpdateJobUsecase_Execute_TableDriven(t *testing.T) {
	type fields struct {
		jobRepo   *mocksRepositories.Repository
		quoteRepo *mocks.Repository
		priceSvc  *mocksServices.ExchangePrice
	}

	jobID := uuid.New()
	quoteID := uuid.New()
	ctx := context.Background()

	mJob := &job.Job{Id: jobID, QuoteId: quoteID, Status: job.StatusPending}
	mQuote := &quote.Quote{Id: quoteID, FromCurrency: "USD", ToCurrency: "RUB"}

	tests := []struct {
		name    string
		params  usecases.ProcessQuoteUpdateJobUsecaseParams
		setup   func(f fields)
		wantErr bool
	}{
		{
			name:   "Success path",
			params: usecases.ProcessQuoteUpdateJobUsecaseParams{Id: jobID},
			setup: func(f fields) {
				f.jobRepo.On("GetById", ctx, jobID).Return(mJob, nil).Once()
				f.jobRepo.On("UpdateStatus", ctx, mock.MatchedBy(func(p job.UpdateStatusParams) bool {
					return p.Status == job.StatusProcessing
				})).Return(nil).Once()
				f.quoteRepo.On("GetById", ctx, quoteID).Return(mQuote, nil).Once()
				f.priceSvc.On("GetRate", ctx, mock.Anything).Return(100.0, nil).Once()
				f.jobRepo.On("UpdatePrice", ctx, mock.Anything).Return(nil).Once()
				f.quoteRepo.On("UpdatePrice", ctx, mock.Anything).Return(nil).Once()
				f.jobRepo.On("UpdateStatus", ctx, mock.MatchedBy(func(p job.UpdateStatusParams) bool {
					return p.Status == job.StatusSuccess
				})).Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "Job not found",
			params: usecases.ProcessQuoteUpdateJobUsecaseParams{Id: jobID},
			setup: func(f fields) {
				f.jobRepo.On("GetById", ctx, jobID).Return(nil, fmt.Errorf("db error")).Once()
			},
			wantErr: true,
		},
		{
			name:   "External API error should mark job as failure",
			params: usecases.ProcessQuoteUpdateJobUsecaseParams{Id: jobID},
			setup: func(f fields) {
				f.jobRepo.On("GetById", ctx, jobID).Return(mJob, nil).Once()
				f.jobRepo.On("UpdateStatus", ctx, mock.MatchedBy(func(p job.UpdateStatusParams) bool {
					return p.Status == job.StatusProcessing
				})).Return(nil).Once()
				f.quoteRepo.On("GetById", ctx, quoteID).Return(mQuote, nil).Once()
				f.priceSvc.On("GetRate", ctx, mock.Anything).Return(0.0, fmt.Errorf("api timeout")).Once()
				f.jobRepo.On("UpdateStatus", ctx, mock.MatchedBy(func(p job.UpdateStatusParams) bool {
					return p.Status == job.StatusFailure
				})).Return(nil).Once()
			},
			wantErr: true,
		},
		{
			name:   "Quote repository error",
			params: usecases.ProcessQuoteUpdateJobUsecaseParams{Id: jobID},
			setup: func(f fields) {
				f.jobRepo.On("GetById", ctx, jobID).Return(mJob, nil).Once()
				f.jobRepo.On("UpdateStatus", ctx, mock.Anything).Return(nil).Once()
				f.quoteRepo.On("GetById", ctx, quoteID).Return(nil, fmt.Errorf("quote db error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := fields{
				jobRepo:   new(mocksRepositories.Repository),
				quoteRepo: new(mocks.Repository),
				priceSvc:  new(mocksServices.ExchangePrice),
			}

			if tt.setup != nil {
				tt.setup(f)
			}

			u := usecases.NewProcessQuoteUpdateUsecase(usecases.ProcessQuoteUpdateUsecaseInitParams{
				JobRepository:   f.jobRepo,
				QuoteRepository: f.quoteRepo,
				ExchangeService: f.priceSvc,
			})

			err := u.Execute(ctx, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			f.jobRepo.AssertExpectations(t)
			f.quoteRepo.AssertExpectations(t)
			f.priceSvc.AssertExpectations(t)
		})
	}
}
