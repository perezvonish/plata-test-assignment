package usecases

import "context"

type QuoteUpdateUsecase interface {
	Execute(ctx context.Context)
}
