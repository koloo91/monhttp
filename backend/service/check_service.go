package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/repository"
	"time"
)

func GetChecks(ctx context.Context, serviceId string, from, to time.Time) ([]model.Check, error) {
	return repository.SelectChecks(ctx, serviceId, from, to)
}

func GetAverageValues(ctx context.Context, serviceId string) (model.Average, error) {
	oneDayBefore := time.Now().AddDate(0, 0, -1)
	averageLastDay, err := repository.SelectAverageLatency(ctx, serviceId, oneDayBefore, time.Now())
	if err != nil {
		return model.Average{}, err
	}

	oneWeekBefore := time.Now().AddDate(0, 0, -7)
	averageLastWeek, err := repository.SelectAverageLatency(ctx, serviceId, oneWeekBefore, time.Now())
	if err != nil {
		return model.Average{}, err
	}

	return model.Average{
		LastDay:  averageLastDay,
		LastWeek: averageLastWeek,
	}, nil
}

func GetIsOnline(ctx context.Context, serviceId string) (bool, error) {
	isOnline, err := repository.SelectIsOnline(ctx, serviceId)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return isOnline, err
}
