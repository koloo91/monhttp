package model

import (
	"github.com/google/uuid"
	"time"
)

type Job struct {
	Id        string
	ServiceId string
	ExecuteAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewJob(serviceId string) Job {
	return Job{
		Id:        uuid.New().String(),
		ServiceId: serviceId,
		ExecuteAt: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
