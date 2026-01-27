package usecases

import "context"

type QuoteGetLatestUsecase interface {
	Execute(ctx context.Context)
}
