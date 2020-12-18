package repository

import (
	"context"
	"database/sql"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	insertServiceStatement            *sql.Stmt
	selectServicesStatement           *sql.Stmt
	selectServiceByIdStatement        *sql.Stmt
	updateServiceByIdStatement        *sql.Stmt
	deleteServiceByIdStatement        *sql.Stmt
	getNextScheduledServicesStatement *sql.Stmt
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
														FROM service ORDER BY name;`)
	if err != nil {
		log.Fatal(err)
	}

	selectServiceByIdStatement, err = db.Prepare(`SELECT id,
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
														FROM service WHERE id = $1;`)
	if err != nil {
		log.Fatal(err)
	}

	updateServiceByIdStatement, err = db.Prepare(`UPDATE service
														SET name=$2,
															type=$3,
															interval_in_seconds=$4,
															next_check_time=$5,
															endpoint=$6,
															http_method=$7,
															request_timeout_in_seconds=$8,
															http_headers=$9,
															http_body=$10,
															expected_http_response_body=$11,
															expected_http_status_code=$12,
															follow_redirects=$13,
															verify_ssl=$14,
															enable_notifications=$15,
															notify_after_number_of_failures=$16,
															updated_at=$17
														WHERE id = $1;`)
	if err != nil {
		log.Fatal(err)
	}

	deleteServiceByIdStatement, err = db.Prepare(`DELETE
														FROM service
														WHERE id = $1;`)
	if err != nil {
		log.Fatal(err)
	}

	getNextScheduledServicesStatement, err = db.Prepare(`SELECT id FROM service WHERE next_check_time <= now();`)
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

func SelectServiceById(ctx context.Context, serviceId string) (model.Service, error) {
	row := selectServiceByIdStatement.QueryRowContext(ctx, serviceId)

	var id, name, endpoint, httpMethod, httpHeaders, httpBody, expectedHttpResponseBody string
	var serviceType model.ServiceType
	var intervalInSeconds, requestTimeoutInSeconds, expectedHttpStatusCode, notifyAfterNumberOfFailures int
	var followRedirects, verifySsl, enableNotifications bool
	var nextCheckTime, createdAt, updatedAt time.Time

	if err := row.Scan(&id, &name, &serviceType, &intervalInSeconds, &nextCheckTime, &endpoint, &httpMethod,
		&requestTimeoutInSeconds, &httpHeaders, &httpBody, &expectedHttpResponseBody,
		&expectedHttpStatusCode, &followRedirects, &verifySsl, &enableNotifications,
		&notifyAfterNumberOfFailures, &createdAt, &updatedAt); err != nil {
		return model.Service{}, err
	}

	return model.Service{
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
	}, nil
}

func SelectServiceByIdLocked(ctx context.Context, tx *sql.Tx, serviceId string) (model.Service, error) {
	row := tx.QueryRowContext(ctx, `SELECT id,
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
											FROM service WHERE id = $1 
											FOR UPDATE;`, serviceId)

	var id, name, endpoint, httpMethod, httpHeaders, httpBody, expectedHttpResponseBody string
	var serviceType model.ServiceType
	var intervalInSeconds, requestTimeoutInSeconds, expectedHttpStatusCode, notifyAfterNumberOfFailures int
	var followRedirects, verifySsl, enableNotifications bool
	var nextCheckTime, createdAt, updatedAt time.Time

	if err := row.Scan(&id, &name, &serviceType, &intervalInSeconds, &nextCheckTime, &endpoint, &httpMethod,
		&requestTimeoutInSeconds, &httpHeaders, &httpBody, &expectedHttpResponseBody,
		&expectedHttpStatusCode, &followRedirects, &verifySsl, &enableNotifications,
		&notifyAfterNumberOfFailures, &createdAt, &updatedAt); err != nil {
		return model.Service{}, err
	}

	return model.Service{
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
	}, nil
}

func UpdateServiceById(ctx context.Context, serviceId string, service model.Service) error {
	if _, err := updateServiceByIdStatement.ExecContext(ctx, serviceId, service.Name, service.Type, service.IntervalInSeconds, service.NextCheckTime, service.Endpoint, service.HttpMethod,
		service.RequestTimeoutInSeconds, service.HttpHeaders, service.HttpBody, service.ExpectedHttpResponseBody,
		service.ExpectedHttpStatusCode, service.FollowRedirects, service.VerifySsl, service.EnableNotifications,
		service.NotifyAfterNumberOfFailures, time.Now()); err != nil {
		return err
	}
	return nil
}

func UpdateServiceByIdNextCheckTime(ctx context.Context, tx *sql.Tx, serviceId string, nextCheckTime time.Time) error {
	if _, err := tx.ExecContext(ctx, `UPDATE service SET next_check_time = $2 WHERE id = $1;`,
		serviceId, nextCheckTime); err != nil {
		return err
	}
	return nil
}

func DeleteServiceById(ctx context.Context, serviceId string) error {
	if _, err := deleteServiceByIdStatement.ExecContext(ctx, serviceId); err != nil {
		return err
	}
	return nil
}

func GetNextScheduledServiceIds(ctx context.Context) ([]string, error) {
	rows, err := getNextScheduledServicesStatement.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id string

	result := make([]string, 0)

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		result = append(result, id)
	}

	return result, nil
}
