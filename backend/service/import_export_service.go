package service

import (
	"context"
	"encoding/csv"
	"errors"
	"github.com/google/uuid"
	"github.com/koloo91/monhttp/model"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidServiceType = errors.New("invalid service type. must be one of [HTTP, ICMP_PING]")
	ErrHttpMethod         = errors.New("invalid http method. must be one of [GET, POST, PUT, PATCH, DELETE]")
)

const (
	nameIndex = iota
	typeIndex
	intervalInSecondsIndex
	endpointIndex
	httpMethodIndex
	requestTimeoutInSecondsIndex
	httpHeadersIndex
	httpBodyIndex
	expectedResponseBodyIndex
	expectedStatusCodeIndex
	followRedirectsIndex
	verifySslIndex
	enableNotificationsIndex
	notifyAfterNumberOfFailuresIndex
	continuouslySendNotificationsIndex
	notifiersIndex
)

func ImportCsvData(ctx context.Context, file io.Reader) []model.ImportResult {
	csvReader := csv.NewReader(file)

	// read header line
	_, _ = csvReader.Read()

	result := make([]model.ImportResult, 0)
	rowNumber := 0

	for {
		rowNumber++

		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			result = append(result, model.ImportResult{RowNumber: rowNumber, Error: err})
			continue
		}

		service, err := csvRowToService(record)
		result = append(result, model.ImportResult{RowNumber: rowNumber, Service: service, Error: err})
	}

	for _, entry := range result {
		if entry.Error != nil {
			continue
		}

		createdService, err := CreateService(ctx, entry.Service)
		entry.Service = createdService
		entry.Error = err
	}

	return result
}

func csvRowToService(row []string) (model.Service, error) {
	name := strings.TrimSpace(row[nameIndex])

	var serviceType model.ServiceType
	switch strings.TrimSpace(row[typeIndex]) {
	case model.ServiceTypeIcmpPing:
		serviceType = model.ServiceTypeIcmpPing
	case model.ServiceTypeHttp:
		serviceType = model.ServiceTypeHttp
	default:
		return model.Service{}, ErrInvalidServiceType
	}

	intervalInSeconds := strings.TrimSpace(row[intervalInSecondsIndex])
	intervalInSecondsInt, err := strconv.Atoi(intervalInSeconds)
	if err != nil {
		return model.Service{}, err
	}

	endpoint := strings.TrimSpace(row[endpointIndex])
	httpMethod := strings.TrimSpace(row[httpMethodIndex])
	if httpMethod != "GET" && httpMethod != "POST" && httpMethod != "PUT" && httpMethod != "PATCH" && httpMethod != "DELETE" {
		return model.Service{}, ErrHttpMethod
	}

	requestTimeoutInSeconds := strings.TrimSpace(row[requestTimeoutInSecondsIndex])
	requestTimeoutInSecondsInt, err := strconv.Atoi(requestTimeoutInSeconds)
	if err != nil {
		return model.Service{}, err
	}

	httpHeaders := strings.TrimSpace(row[httpHeadersIndex])
	httpBody := strings.TrimSpace(row[httpBodyIndex])
	expectedResponseBody := strings.TrimSpace(row[expectedResponseBodyIndex])

	expectedStatusCode := strings.TrimSpace(row[expectedStatusCodeIndex])
	expectedStatusCodeInt, err := strconv.Atoi(expectedStatusCode)
	if err != nil {
		return model.Service{}, err
	}

	followRedirects := strings.TrimSpace(row[followRedirectsIndex])
	followRedirectsBool, err := strconv.ParseBool(followRedirects)
	if err != nil {
		return model.Service{}, err
	}

	verifySsl := strings.TrimSpace(row[verifySslIndex])
	verifySslBool, err := strconv.ParseBool(verifySsl)
	if err != nil {
		return model.Service{}, err
	}

	enableNotifications := strings.TrimSpace(row[enableNotificationsIndex])
	enableNotificationsBool, err := strconv.ParseBool(enableNotifications)
	if err != nil {
		return model.Service{}, err
	}

	notifyAfterNumberOfFailures := strings.TrimSpace(row[notifyAfterNumberOfFailuresIndex])
	notifyAfterNumberOfFailuresInt, err := strconv.Atoi(notifyAfterNumberOfFailures)
	if err != nil {
		return model.Service{}, err
	}

	continuouslySendNotifications := strings.TrimSpace(row[continuouslySendNotificationsIndex])
	continuouslySendNotificationsBool, err := strconv.ParseBool(continuouslySendNotifications)
	if err != nil {
		return model.Service{}, err
	}

	notifiers := strings.TrimSpace(row[notifiersIndex])
	notifiersSlice := strings.Split(notifiers, ",")

	return model.Service{
		Id:                            uuid.New().String(),
		Name:                          name,
		Type:                          serviceType,
		IntervalInSeconds:             intervalInSecondsInt,
		Endpoint:                      endpoint,
		HttpMethod:                    httpMethod,
		RequestTimeoutInSeconds:       requestTimeoutInSecondsInt,
		HttpHeaders:                   httpHeaders,
		HttpBody:                      httpBody,
		ExpectedHttpResponseBody:      expectedResponseBody,
		ExpectedHttpStatusCode:        expectedStatusCodeInt,
		FollowRedirects:               followRedirectsBool,
		VerifySsl:                     verifySslBool,
		EnableNotifications:           enableNotificationsBool,
		NotifyAfterNumberOfFailures:   notifyAfterNumberOfFailuresInt,
		ContinuouslySendNotifications: continuouslySendNotificationsBool,
		Notifiers:                     notifiersSlice,
		CreatedAt:                     time.Now(),
		UpdatedAt:                     time.Now(),
	}, nil
}
