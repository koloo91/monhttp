package repository

import (
	"context"
	"database/sql"
	"github.com/koloo91/monhttp/model"
	"time"
)

const (
	insertJobQuery           = `INSERT INTO job(id, service_id, execute_at, created_at, updated_at) VALUES($1, $2, $3, $4, $5);`
	getNextJobIdsQuery       = `SELECT id FROM job WHERE execute_at <= now();`
	selectJobByIdLockedQuery = `SELECT service_id,
									   execute_at,
									   created_at,
									   updated_at
								FROM job 
								WHERE id = $1 
								FOR UPDATE;`
	updateJobByIdExecuteAtQuery = `UPDATE job SET execute_at = $2 WHERE id = $1;`
)

func InsertJobTx(ctx context.Context, tx *sql.Tx, job model.Job) error {
	if _, err := tx.ExecContext(ctx, insertJobQuery, job.Id, job.ServiceId, job.ExecuteAt, job.CreatedAt, job.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func GetNextJobIds(ctx context.Context) ([]string, error) {
	rows, err := db.QueryContext(ctx, getNextJobIdsQuery)
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

func SelectJobByIdLockedTx(ctx context.Context, tx *sql.Tx, id string) (model.Job, error) {
	row := tx.QueryRowContext(ctx, selectJobByIdLockedQuery, id)

	var serviceId string
	var executeAt, createdAt, updatedAt time.Time

	if err := row.Scan(&serviceId, &executeAt, &createdAt, &updatedAt); err != nil {
		return model.Job{}, err
	}

	return model.Job{
		Id:        id,
		ServiceId: serviceId,
		ExecuteAt: executeAt,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func UpdateJobByIdExecuteAtTx(ctx context.Context, tx *sql.Tx, id string, executeAt time.Time) error {
	if _, err := tx.ExecContext(ctx, updateJobByIdExecuteAtQuery, id, executeAt); err != nil {
		return err
	}
	return nil
}
