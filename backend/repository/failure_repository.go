package repository

import (
	"context"
	"database/sql"
	"github.com/koloo91/monhttp/model"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	selectFailuresByServiceIdAndCreateAtStatement              *sql.Stmt
	selectFailuresCountByServiceIdAnCreatedAtStatement         *sql.Stmt
	selectFailuresCountByServiceIdAndGroupByCreatedAtStatement *sql.Stmt
)

func prepareFailureStatements() {
	var err error
	selectFailuresByServiceIdAndCreateAtStatement, err = db.Prepare(`SELECT id, reason, created_at
																			FROM failure
																			WHERE service_id = $1
																			  AND created_at >= $2
																			  AND created_at <= $3
																			ORDER BY created_at DESC
																			LIMIT $4 
																			OFFSET $5;`)

	if err != nil {
		log.Fatal(err)
	}

	selectFailuresCountByServiceIdAnCreatedAtStatement, err = db.Prepare(`SELECT COUNT(id)
																					FROM failure
																					WHERE service_id = $1
																					AND created_at >= $2
																					AND created_at <= $3;`)
	if err != nil {
		log.Fatal(err)
	}

	selectFailuresCountByServiceIdAndGroupByCreatedAtStatement, err = db.Prepare(`SELECT generated.day, COALESCE(failure.count, 0)
																							FROM (
																									 SELECT date_trunc('day', dd)::date as day, 0 as count
																									 FROM generate_series($2::timestamptz, $3::timestamptz,
																														  '1 day'::interval) dd) as generated
																									 LEFT JOIN (SELECT date_trunc('day', created_at)::date as day, COUNT(*) as count
																												FROM failure
																												WHERE service_id = $1
																												  AND created_at >= $2
																												  AND created_at <= $3
																												GROUP BY 1
																												ORDER BY 1) as failure ON generated.day = failure.day;`)
}

func InsertFailure(ctx context.Context, tx *sql.Tx, failure model.Failure) error {
	if _, err := tx.ExecContext(ctx, `INSERT INTO failure (id, service_id, reason, created_at) 
											VALUES ($1, $2, $3, $4)`,
		failure.Id, failure.ServiceId, failure.Reason, failure.CreatedAt); err != nil {
		return err
	}
	return nil
}

func SelectFailures(ctx context.Context, serviceId string, from, to time.Time, limit, offset int) ([]model.Failure, error) {
	rows, err := selectFailuresByServiceIdAndCreateAtStatement.QueryContext(ctx, serviceId, from, to, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id, reason string
	var createdAt time.Time

	result := make([]model.Failure, 0)

	for rows.Next() {
		if err := rows.Scan(&id, &reason, &createdAt); err != nil {
			return nil, err
		}

		result = append(result, model.Failure{
			Id:        id,
			ServiceId: serviceId,
			Reason:    reason,
			CreatedAt: createdAt,
		})
	}

	return result, nil
}

func SelectFailuresCount(ctx context.Context, serviceId string, from, to time.Time) (int, error) {
	row := selectFailuresCountByServiceIdAnCreatedAtStatement.QueryRowContext(ctx, serviceId, from, to)

	var count int

	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func SelectFailuresGroupedByDay(ctx context.Context, serviceId string, from, to time.Time) ([]model.FailureCountByDay, error) {
	rows, err := selectFailuresCountByServiceIdAndGroupByCreatedAtStatement.QueryContext(ctx, serviceId, from, to)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var day time.Time
	var count int

	result := make([]model.FailureCountByDay, 0)

	for rows.Next() {
		if err := rows.Scan(&day, &count); err != nil {
			return nil, err
		}

		result = append(result, model.FailureCountByDay{
			Day:   day,
			Count: count,
		})
	}

	return result, nil
}
