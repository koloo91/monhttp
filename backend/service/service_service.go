package service

import (
	"context"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/repository"
)

func CreateService(ctx context.Context, service model.Service) (model.Service, error) {
	return service, repository.InsertService(ctx, service)
}

func GetServices(ctx context.Context) ([]model.Service, error) {
	return repository.SelectServices(ctx)
}
