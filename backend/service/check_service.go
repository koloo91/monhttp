package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/koloo91/monhttp/model"
	"github.com/koloo91/monhttp/repository"
	"time"
)

func GetChecks(ctx context.Context, serviceId string, from, to time.Time, reduceByFactor int) ([]model.Check, error) {
	if reduceByFactor <= 0 {
		reduceByFactor = 1
	}
	return repository.SelectChecks(ctx, serviceId, from, to, reduceByFactor)
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

	var lastDayUptime float64 = 0
	if successLastDay+failuresLastDay > 0 {
		lastDayUptime = (successLastDay / (successLastDay + failuresLastDay)) * 100
	}

	var lastWeekUptime float64 = 0
	if successLastWeek+failuresLastWeek > 0 {
		lastWeekUptime = (successLastWeek / (successLastWeek + failuresLastWeek)) * 100
	}

	return model.Average{
		LastDayResponseTime:  averageResponseTimeLastDay,
		LastWeekResponseTime: averageResponseTimeLastWeek,
		LastDayUptime:        lastDayUptime,
		LastWeekUptime:       lastWeekUptime,
	}, nil
}

func GetIsOnline(ctx context.Context, serviceId string) (bool, error) {
	isOnline, err := repository.SelectIsOnline(ctx, serviceId)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return isOnline, err
}
