package repository

import (
	"context"
	"database/sql"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	insertServiceStatement  *sql.Stmt
	selectServicesStatement *sql.Stmt
)

func prepareServiceStatements() {
	var err error

	insertServiceStatement, err = db.Prepare(`INSERT INTO service (id, name, type, interval_in_seconds, next_check_time, endpoint, http_method,
											 request_timeout_in_seconds, http_headers, http_body, expected_http_response_body,
											 expected_http_status_code, follow_redirects, verify_ssl, enable_notifications,
											 notify_after_number_of_failures, created_at, updated_at)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);`)
	if err != nil {
		log.Fatal(err)
	}

	selectServicesStatement, err = db.Prepare(`SELECT id,
						   name,
						   type,
						   interval_in_seconds,
						   next_check_time,
						   endpoint,
						   http_method,
						   request_timeout_in_seconds,
						   http_headers,
						   http_body,
						   expected_http_response_body,
						   expected_http_status_code,
						   follow_redirects,
						   verify_ssl,
						   enable_notifications,
						   notify_after_number_of_failures,
						   created_at,
						   updated_at
					FROM service;`)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertService(ctx context.Context, service model.Service) error {
	if _, err := insertServiceStatement.ExecContext(ctx,
		service.Id, service.Name, service.Type, service.IntervalInSeconds, service.NextCheckTime, service.Endpoint, service.HttpMethod,
		service.RequestTimeoutInSeconds, service.HttpHeaders, service.HttpBody, service.ExpectedHttpResponseBody,
		service.ExpectedHttpStatusCode, service.FollowRedirects, service.VerifySsl, service.EnableNotifications,
		service.NotifyAfterNumberOfFailures, service.CreatedAt, service.UpdatedAt); err != nil {
		return err
	}

	return nil
}

// TODO: add paging
func SelectServices(ctx context.Context) ([]model.Service, error) {
	rows, err := selectServicesStatement.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var id, name, endpoint, httpMethod, httpHeaders, httpBody, expectedHttpResponseBody string
	var serviceType model.ServiceType
	var intervalInSeconds, requestTimeoutInSeconds, expectedHttpStatusCode, notifyAfterNumberOfFailures int
	var followRedirects, verifySsl, enableNotifications bool
	var nextCheckTime, createdAt, updatedAt time.Time

	result := make([]model.Service, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &name, &serviceType, &intervalInSeconds, &nextCheckTime, &endpoint, &httpMethod,
			&requestTimeoutInSeconds, &httpHeaders, &httpBody, &expectedHttpResponseBody,
			&expectedHttpStatusCode, &followRedirects, &verifySsl, &enableNotifications,
			&notifyAfterNumberOfFailures, &createdAt, &updatedAt); err != nil {
			return nil, err
		}

		result = append(result, model.Service{
			Id:                          id,
			Name:                        name,
			Type:                        serviceType,
			IntervalInSeconds:           intervalInSeconds,
			NextCheckTime:               nextCheckTime,
			Endpoint:                    endpoint,
			HttpMethod:                  httpMethod,
			RequestTimeoutInSeconds:     requestTimeoutInSeconds,
			HttpHeaders:                 httpHeaders,
			HttpBody:                    httpBody,
			ExpectedHttpResponseBody:    expectedHttpResponseBody,
			ExpectedHttpStatusCode:      expectedHttpStatusCode,
			FollowRedirects:             followRedirects,
			VerifySsl:                   verifySsl,
			EnableNotifications:         enableNotifications,
			NotifyAfterNumberOfFailures: notifyAfterNumberOfFailures,
			CreatedAt:                   createdAt,
			UpdatedAt:                   updatedAt,
		})
	}

	return result, nil
}
