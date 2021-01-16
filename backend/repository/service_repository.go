package repository

import (
	"context"
	"database/sql"
	"github.com/koloo91/monhttp/model"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	insertServiceQuery = `INSERT INTO service (id, name, type, interval_in_seconds, endpoint, http_method,
											 request_timeout_in_seconds, http_headers, http_body, expected_http_response_body,
											 expected_http_status_code, follow_redirects, verify_ssl, enable_notifications,
											 notify_after_number_of_failures, continuously_send_notifications, notifiers, created_at, updated_at)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19);`
)

var (
	insertServiceStatement     *sql.Stmt
	selectServicesStatement    *sql.Stmt
	selectServiceByIdStatement *sql.Stmt
	updateServiceByIdStatement *sql.Stmt
	deleteServiceByIdStatement *sql.Stmt
)

func prepareServiceStatements() {
	var err error

	insertServiceStatement, err = db.Prepare(insertServiceQuery)
	if err != nil {
		log.Fatal(err)
	}

	selectServicesStatement, err = db.Prepare(`SELECT id,
															   name,
															   type,
															   interval_in_seconds,
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
       														   continuously_send_notifications,
       														   notifiers,
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
       														   continuously_send_notifications,
       														   notifiers,
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
															endpoint=$5,
															http_method=$6,
															request_timeout_in_seconds=$7,
															http_headers=$8,
															http_body=$9,
															expected_http_response_body=$10,
															expected_http_status_code=$11,
															follow_redirects=$12,
															verify_ssl=$13,
															enable_notifications=$14,
															notify_after_number_of_failures=$15,
														    continuously_send_notifications=$16,
														    notifiers=$17,
															updated_at=$18
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
}

func InsertService(ctx context.Context, service model.Service) error {
	if _, err := insertServiceStatement.ExecContext(ctx,
		service.Id, service.Name, service.Type, service.IntervalInSeconds, service.Endpoint, service.HttpMethod,
		service.RequestTimeoutInSeconds, service.HttpHeaders, service.HttpBody, service.ExpectedHttpResponseBody,
		service.ExpectedHttpStatusCode, service.FollowRedirects, service.VerifySsl, service.EnableNotifications,
		service.NotifyAfterNumberOfFailures, service.ContinuouslySendNotifications, pq.Array(service.Notifiers),
		service.CreatedAt, service.UpdatedAt); err != nil {
		return err
	}

	return nil
}

func InsertServiceTx(ctx context.Context, tx *sql.Tx, service model.Service) error {
	if _, err := tx.ExecContext(ctx, insertServiceQuery,
		service.Id, service.Name, service.Type, service.IntervalInSeconds, service.Endpoint, service.HttpMethod,
		service.RequestTimeoutInSeconds, service.HttpHeaders, service.HttpBody, service.ExpectedHttpResponseBody,
		service.ExpectedHttpStatusCode, service.FollowRedirects, service.VerifySsl, service.EnableNotifications,
		service.NotifyAfterNumberOfFailures, service.ContinuouslySendNotifications, pq.Array(service.Notifiers),
		service.CreatedAt, service.UpdatedAt); err != nil {
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
	var followRedirects, verifySsl, enableNotifications, continuouslySendNotifications bool
	var notifiers []string
	var createdAt, updatedAt time.Time

	result := make([]model.Service, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &name, &serviceType, &intervalInSeconds, &endpoint, &httpMethod,
			&requestTimeoutInSeconds, &httpHeaders, &httpBody, &expectedHttpResponseBody,
			&expectedHttpStatusCode, &followRedirects, &verifySsl, &enableNotifications,
			&notifyAfterNumberOfFailures, &continuouslySendNotifications, pq.Array(&notifiers),
			&createdAt, &updatedAt); err != nil {
			return nil, err
		}

		result = append(result, model.Service{
			Id:                            id,
			Name:                          name,
			Type:                          serviceType,
			IntervalInSeconds:             intervalInSeconds,
			Endpoint:                      endpoint,
			HttpMethod:                    httpMethod,
			RequestTimeoutInSeconds:       requestTimeoutInSeconds,
			HttpHeaders:                   httpHeaders,
			HttpBody:                      httpBody,
			ExpectedHttpResponseBody:      expectedHttpResponseBody,
			ExpectedHttpStatusCode:        expectedHttpStatusCode,
			FollowRedirects:               followRedirects,
			VerifySsl:                     verifySsl,
			EnableNotifications:           enableNotifications,
			NotifyAfterNumberOfFailures:   notifyAfterNumberOfFailures,
			ContinuouslySendNotifications: continuouslySendNotifications,
			Notifiers:                     notifiers,
			CreatedAt:                     createdAt,
			UpdatedAt:                     updatedAt,
		})
	}

	return result, nil
}

func SelectServiceById(ctx context.Context, serviceId string) (model.Service, error) {
	row := selectServiceByIdStatement.QueryRowContext(ctx, serviceId)

	var id, name, endpoint, httpMethod, httpHeaders, httpBody, expectedHttpResponseBody string
	var serviceType model.ServiceType
	var intervalInSeconds, requestTimeoutInSeconds, expectedHttpStatusCode, notifyAfterNumberOfFailures int
	var followRedirects, verifySsl, enableNotifications, continuouslySendNotifications bool
	var notifiers []string
	var createdAt, updatedAt time.Time

	if err := row.Scan(&id, &name, &serviceType, &intervalInSeconds, &endpoint, &httpMethod,
		&requestTimeoutInSeconds, &httpHeaders, &httpBody, &expectedHttpResponseBody,
		&expectedHttpStatusCode, &followRedirects, &verifySsl, &enableNotifications,
		&notifyAfterNumberOfFailures, &continuouslySendNotifications, pq.Array(&notifiers),
		&createdAt, &updatedAt); err != nil {
		return model.Service{}, err
	}

	return model.Service{
		Id:                            id,
		Name:                          name,
		Type:                          serviceType,
		IntervalInSeconds:             intervalInSeconds,
		Endpoint:                      endpoint,
		HttpMethod:                    httpMethod,
		RequestTimeoutInSeconds:       requestTimeoutInSeconds,
		HttpHeaders:                   httpHeaders,
		HttpBody:                      httpBody,
		ExpectedHttpResponseBody:      expectedHttpResponseBody,
		ExpectedHttpStatusCode:        expectedHttpStatusCode,
		FollowRedirects:               followRedirects,
		VerifySsl:                     verifySsl,
		EnableNotifications:           enableNotifications,
		NotifyAfterNumberOfFailures:   notifyAfterNumberOfFailures,
		ContinuouslySendNotifications: continuouslySendNotifications,
		Notifiers:                     notifiers,
		CreatedAt:                     createdAt,
		UpdatedAt:                     updatedAt,
	}, nil
}

func UpdateServiceById(ctx context.Context, serviceId string, service model.Service) error {
	if _, err := updateServiceByIdStatement.ExecContext(ctx, serviceId, service.Name, service.Type, service.IntervalInSeconds,
		service.Endpoint, service.HttpMethod, service.RequestTimeoutInSeconds, service.HttpHeaders, service.HttpBody,
		service.ExpectedHttpResponseBody, service.ExpectedHttpStatusCode, service.FollowRedirects, service.VerifySsl,
		service.EnableNotifications, service.NotifyAfterNumberOfFailures, service.ContinuouslySendNotifications,
		pq.Array(service.Notifiers), time.Now()); err != nil {
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
