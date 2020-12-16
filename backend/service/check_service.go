package service

import (
	"context"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/repository"
	"time"
)

func GetChecks(ctx context.Context, serviceId string, from time.Time, to time.Time) ([]model.Check, error) {
	return repository.SelectChecks(ctx, serviceId, from, to)
}
