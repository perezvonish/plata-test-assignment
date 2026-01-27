package usecases

import (
	"context"
)

type HealthPingResult struct {
	Ok bool `json:"ok"`
}

type HealthPingUsecase interface {
	Execute(ctx context.Context) (HealthPingResult, error)
}

type HealthPingUsecaseImpl struct{}

func (h HealthPingUsecaseImpl) Execute(ctx context.Context) (HealthPingResult, error) {
	return HealthPingResult{Ok: true}, nil
}

func NewHealthPingUsecase() HealthPingUsecase {
	return HealthPingUsecaseImpl{}
}
