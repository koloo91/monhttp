package service

import (
	"context"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/repository"
	"time"
)

func GetFailures(ctx context.Context, serviceId string, from, to time.Time) ([]model.Failure, error) {
	return repository.SelectFailures(ctx, serviceId, from, to)
}

func GetFailuresCount(ctx context.Context, serviceId string, from, to time.Time) (model.FailureCount, error) {
	count, err := repository.SelectFailuresCount(ctx, serviceId, from, to)
	return model.FailureCount{Count: count}, err
}
