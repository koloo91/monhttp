package service

import (
	"context"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/repository"
)

func CreateService(ctx context.Context, service model.Service) (model.Service, error) {
	tx, err := repository.BeginnTransaction()
	if err != nil {
		return model.Service{}, err
	}

	if err := repository.InsertServiceTx(ctx, tx, service); err != nil {
		if err := tx.Rollback(); err != nil {
			return model.Service{}, err
		}
		return model.Service{}, err
	}

	if err := repository.InsertJobTx(ctx, tx, model.NewJob(service.Id)); err != nil {
		if err := tx.Rollback(); err != nil {
			return model.Service{}, nil
		}
		return model.Service{}, nil
	}

	if err := tx.Commit(); err != nil {
		return model.Service{}, err
	}

	return service, repository.InsertService(ctx, service)
}

func GetServices(ctx context.Context) ([]model.Service, error) {
	return repository.SelectServices(ctx)
}

func GetServiceById(ctx context.Context, id string) (model.Service, error) {
	return repository.SelectServiceById(ctx, id)
}

func UpdateServiceById(ctx context.Context, id string, service model.Service) (model.Service, error) {
	if err := repository.UpdateServiceById(ctx, id, service); err != nil {
		return model.Service{}, nil
	}

	return repository.SelectServiceById(ctx, id)
}

func DeleteServiceById(ctx context.Context, id string) error {
	return repository.DeleteServiceById(ctx, id)
}
