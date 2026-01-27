package usecases

import "context"

type QuoteGetByUpdateIdUsecase interface {
	Execute(ctx context.Context)
}
