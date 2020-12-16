package model

import (
	"github.com/google/uuid"
	"time"
)

type Check struct {
	Id           string
	ServiceId    string
	LatencyInMs  int64
	PingTimeInMs int64
	IsFailure    bool
	CreatedAt    time.Time
}

type CheckVo struct {
	Id           string    `json:"id"`
	ServiceId    string    `json:"serviceId"`
	LatencyInMs  int64     `json:"latencyInMs"`
	PingTimeInMs int64     `json:"pingTimeInMs"`
	IsFailure    bool      `json:"isFailure"`
	CreatedAt    time.Time `json:"createdAt"`
}

func NewCheck(serviceId string, latency, pingTime int64, isFailure bool) *Check {
	return &Check{
		Id:           uuid.New().String(),
		ServiceId:    serviceId,
		LatencyInMs:  latency,
		PingTimeInMs: pingTime,
		IsFailure:    isFailure,
		CreatedAt:    time.Now(),
	}
}

func MapCheckEntityToVo(entity Check) CheckVo {
	return CheckVo{
		Id:           entity.Id,
		ServiceId:    entity.ServiceId,
		LatencyInMs:  entity.LatencyInMs,
		PingTimeInMs: entity.PingTimeInMs,
		IsFailure:    entity.IsFailure,
		CreatedAt:    entity.CreatedAt,
	}
}

func MapCheckEntitiesToVos(entities []Check) []CheckVo {
	result := make([]CheckVo, 0, len(entities))

	return result
}
