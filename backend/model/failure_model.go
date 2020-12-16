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

func NewFailure(serviceId string, reason string) *Failure {
	return &Failure{
		Id:        uuid.New().String(),
		ServiceId: serviceId,
		Reason:    reason,
		CreatedAt: time.Now(),
	}
}
