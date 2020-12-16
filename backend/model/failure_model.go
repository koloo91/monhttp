package model

import (
	"github.com/google/uuid"
	"time"
)

type Failure struct {
	Id        string
	ServiceId string
	Reason    string
	CreatedAt time.Time
}

type FailureVo struct {
	Id        string    `json:"id"`
	ServiceId string    `json:"serviceId"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewFailure(serviceId string, reason string) *Failure {
	return &Failure{
		Id:        uuid.New().String(),
		ServiceId: serviceId,
		Reason:    reason,
		CreatedAt: time.Now(),
	}
}

func MapFailureEntityToVo(entity Failure) FailureVo {
	return FailureVo{
		Id:        entity.Id,
		ServiceId: entity.ServiceId,
		Reason:    entity.Reason,
		CreatedAt: entity.CreatedAt,
	}
}

func MapFailureEntitiesToVos(entities []Failure) []FailureVo {
	result := make([]FailureVo, 0, len(entities))
	for _, entity := range entities {
		result = append(result, MapFailureEntityToVo(entity))
	}
	return result
}
