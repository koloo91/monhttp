package model

import (
	"github.com/docker/distribution/uuid"
	"time"
)

const (
	ServiceTypeHttp     = "HTTP"
	ServiceTypeIcmpPing = "ICMP_PING"
)

type ServiceType string

type Service struct {
	Id                            string
	Name                          string
	Type                          ServiceType
	IntervalInSeconds             int
	NextCheckTime                 time.Time
	Endpoint                      string
	HttpMethod                    string
	RequestTimeoutInSeconds       int
	HttpHeaders                   string
	HttpBody                      string
	ExpectedHttpResponseBody      string
	ExpectedHttpStatusCode        int
	FollowRedirects               bool
	VerifySsl                     bool
	EnableNotifications           bool
	NotifyAfterNumberOfFailures   int
	ContinuouslySendNotifications bool
	CreatedAt                     time.Time
	UpdatedAt                     time.Time
}

type ServiceVo struct {
	Id                            string      `json:"id"`
	Name                          string      `json:"name"`
	Type                          ServiceType `json:"type"`
	IntervalInSeconds             int         `json:"intervalInSeconds"`
	NextCheckTime                 time.Time   `json:"nextCheckTime"`
	Endpoint                      string      `json:"endpoint"`
	HttpMethod                    string      `json:"httpMethod"`
	RequestTimeoutInSeconds       int         `json:"requestTimeoutInSeconds"`
	HttpHeaders                   string      `json:"httpHeaders"`
	HttpBody                      string      `json:"httpBody"`
	ExpectedHttpResponseBody      string      `json:"expectedHttpResponseBody"`
	ExpectedHttpStatusCode        int         `json:"expectedHttpStatusCode"`
	FollowRedirects               bool        `json:"followRedirects"`
	VerifySsl                     bool        `json:"verifySsl"`
	EnableNotifications           bool        `json:"enableNotifications"`
	NotifyAfterNumberOfFailures   int         `json:"notifyAfterNumberOfFailures"`
	ContinuouslySendNotifications bool        `json:"continuouslySendNotifications"`
	CreatedAt                     time.Time   `json:"createdAt"`
	UpdatedAt                     time.Time   `json:"updatedAt"`
}

func MapServiceVoToEntity(vo ServiceVo) Service {
	return Service{
		Id:                            uuid.Generate().String(),
		Name:                          vo.Name,
		Type:                          vo.Type,
		IntervalInSeconds:             vo.IntervalInSeconds,
		NextCheckTime:                 time.Now(),
		Endpoint:                      vo.Endpoint,
		HttpMethod:                    vo.HttpMethod,
		RequestTimeoutInSeconds:       vo.RequestTimeoutInSeconds,
		HttpHeaders:                   vo.HttpHeaders,
		HttpBody:                      vo.HttpBody,
		ExpectedHttpResponseBody:      vo.ExpectedHttpResponseBody,
		ExpectedHttpStatusCode:        vo.ExpectedHttpStatusCode,
		FollowRedirects:               vo.FollowRedirects,
		VerifySsl:                     vo.VerifySsl,
		EnableNotifications:           vo.EnableNotifications,
		NotifyAfterNumberOfFailures:   vo.NotifyAfterNumberOfFailures,
		ContinuouslySendNotifications: vo.ContinuouslySendNotifications,
		CreatedAt:                     time.Now(),
		UpdatedAt:                     time.Now(),
	}
}

func MapServiceEntityToVo(entity Service) ServiceVo {
	return ServiceVo{
		Id:                            entity.Id,
		Name:                          entity.Name,
		Type:                          entity.Type,
		IntervalInSeconds:             entity.IntervalInSeconds,
		NextCheckTime:                 entity.NextCheckTime,
		Endpoint:                      entity.Endpoint,
		HttpMethod:                    entity.HttpMethod,
		RequestTimeoutInSeconds:       entity.RequestTimeoutInSeconds,
		HttpHeaders:                   entity.HttpHeaders,
		HttpBody:                      entity.HttpBody,
		ExpectedHttpResponseBody:      entity.ExpectedHttpResponseBody,
		ExpectedHttpStatusCode:        entity.ExpectedHttpStatusCode,
		FollowRedirects:               entity.FollowRedirects,
		VerifySsl:                     entity.VerifySsl,
		EnableNotifications:           entity.EnableNotifications,
		NotifyAfterNumberOfFailures:   entity.NotifyAfterNumberOfFailures,
		ContinuouslySendNotifications: entity.ContinuouslySendNotifications,
		CreatedAt:                     entity.CreatedAt,
		UpdatedAt:                     entity.UpdatedAt,
	}
}

func MapServiceEntitiesToVos(entities []Service) []ServiceVo {
	result := make([]ServiceVo, 0, len(entities))
	for _, entity := range entities {
		result = append(result, MapServiceEntityToVo(entity))
	}
	return result
}
