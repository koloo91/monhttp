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
	averageResponseTimeLastDay, err := repository.SelectAverageLatency(ctx, serviceId, oneDayBefore, time.Now())
	if err != nil {
		return model.Average{}, err
	}

	oneWeekBefore := time.Now().AddDate(0, 0, -7)
	averageResponseTimeLastWeek, err := repository.SelectAverageLatency(ctx, serviceId, oneWeekBefore, time.Now())
	if err != nil {
		return model.Average{}, err
	}

	successLastDay, failuresLastDay, err := repository.SelectUptime(ctx, serviceId, oneDayBefore, time.Now())
	if err != nil {
		return model.Average{}, err
	}

	successLastWeek, failuresLastWeek, err := repository.SelectUptime(ctx, serviceId, oneWeekBefore, time.Now())
	if err != nil {
		return model.Average{}, err
	}

	return model.Average{
		LastDayResponseTime:  averageResponseTimeLastDay,
		LastWeekResponseTime: averageResponseTimeLastWeek,
		LastDayUptime:        (successLastDay / (successLastDay + failuresLastDay)) * 100,
		LastWeekUptime:       (successLastWeek / (successLastWeek + failuresLastWeek)) * 100,
	}, nil
}

func GetIsOnline(ctx context.Context, serviceId string) (bool, error) {
	isOnline, err := repository.SelectIsOnline(ctx, serviceId)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return isOnline, err
}